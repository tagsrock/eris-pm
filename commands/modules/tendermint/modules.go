package commands

import (
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/eris-ltd/eris-pm/epm"

	"github.com/eris-ltd/eris-pm/Godeps/_workspace/src/github.com/eris-ltd/mint-client/mintx/core"
	"github.com/eris-ltd/eris-pm/Godeps/_workspace/src/github.com/gorilla/websocket"
	rpcserver "github.com/eris-ltd/eris-pm/Godeps/_workspace/src/github.com/tendermint/tendermint/rpc/server"
	rpctypes "github.com/eris-ltd/eris-pm/Godeps/_workspace/src/github.com/tendermint/tendermint/rpc/types"
	"github.com/eris-ltd/eris-pm/Godeps/_workspace/src/github.com/tendermint/tendermint/types"
	"github.com/eris-ltd/eris-pm/Godeps/_workspace/src/github.com/tendermint/tendermint/wire"
)

var (
	SIGN      = true
	BROADCAST = true
	WAIT      = // we implement waiting here so we only subscribe once on big pdxs
	false
)

func NewChain(chainID, nodeAddr, signAddr, pubkey string) (epm.ChainClient, error) {
	return NewTendermint(chainID, nodeAddr, signAddr, pubkey)
}

// implements epm.ChainClient
type Tendermint struct {
	chainID        string
	nodeAddr       string
	wsAddr         string
	signAddr       string
	pubkey         string
	inputAddr      string
	inputAddrBytes []byte

	// txs are added to the txPool after they are broadcast
	// they should be removed once they are confirmed in a block
	mtx    sync.Mutex
	txPool map[string]struct{}

	quit           chan struct{}
	cleared        chan struct{} // fires when the txPool is emptied if requestCleared == true
	requestCleared bool

	eid  string          // event id
	conn *websocket.Conn // ws conn
}

func NewTendermint(chainID, nodeAddr, signAddr, pubkey string) (*Tendermint, error) {

	//TODO: get addr from pubkey!
	var addr string
	var inputAddrBytes []byte

	wsAddr := strings.TrimPrefix(nodeAddr, "http://")
	wsAddr = "ws://" + wsAddr + "websocket"
	t := &Tendermint{
		chainID:        chainID,
		nodeAddr:       nodeAddr,
		wsAddr:         wsAddr,
		signAddr:       signAddr,
		pubkey:         pubkey,
		inputAddr:      addr,
		inputAddrBytes: inputAddrBytes,
		txPool:         make(map[string]struct{}),
		quit:           make(chan struct{}),
		cleared:        make(chan struct{}),
	}

	// setup a websocket connection
	// and subscribe to inputs from the sending address
	if err := t.setupWSConn(); err != nil {
		return nil, err
	}

	// waits for transactions sent by this client to be confirmed
	// and clears them from the txPool
	go t.clearTransactions()

	return t, nil
}

func (t *Tendermint) Tx(addr, amt string) (string, error) {
	tx, err := core.Send(t.nodeAddr, t.pubkey, t.inputAddr, addr, amt, "")
	if err != nil {
		return "", err
	}
	result, err := core.SignAndBroadcast(t.chainID, t.nodeAddr, t.signAddr, tx, SIGN, BROADCAST, WAIT)
	if err != nil {
		return "", err
	}
	hashString := fmt.Sprintf("%X", result.Hash)
	t.mtx.Lock()
	t.txPool[hashString] = struct{}{}
	t.mtx.Unlock()
	return hashString, nil
}

func (t *Tendermint) Msg(addr string, data []string) (string, error) {
	return "", nil
}

func (t *Tendermint) Call(addr string, data []string) (string, error) {
	return "", nil

}

func (t *Tendermint) Script(code string) (string, string, error) {
	return "", "", nil

}

func (t *Tendermint) NameReg(name, value string) (string, error) {
	return "", nil

}

func (t *Tendermint) StorageAt(target, storage string) (string, error) {
	return "", nil
}

// Commit waits until the pool of transactions broadcast by the client are confirmed in blocks
func (t *Tendermint) Commit() {
	t.mtx.Lock()
	t.requestCleared = true
	l := len(t.txPool)
	t.mtx.Unlock()

	if l > 0 {
		<-t.cleared
	}

	t.mtx.Lock()
	t.requestCleared = false
	t.mtx.Unlock()
}

//--------------------------------------------------------
// functions for managing web socket conn

func (t *Tendermint) setupWSConn() error {
	dialer := websocket.DefaultDialer
	rHeader := http.Header{}
	conn, _, err := dialer.Dial(t.wsAddr, rHeader)
	if err != nil {
		return fmt.Errorf("Error establishing websocket connection to wait for tx to get committed: %v", err)
	}
	eid := types.EventStringAccInput(t.inputAddrBytes)

	if err := conn.WriteJSON(rpctypes.RPCRequest{
		JSONRPC: "2.0",
		Id:      "",
		Method:  "subscribe",
		Params:  []interface{}{eid},
	}); err != nil {
		return fmt.Errorf("Error subscribing to AccInput event: %v", err)
	}

	// TODO: we should also subscribe to and track new blocks ...

	// run the ping loop
	go func() {
		pingTicker := time.NewTicker((time.Second * rpcserver.WSReadTimeoutSeconds) / 2)
		for {
			select {
			case <-pingTicker.C:
				if err := conn.WriteControl(websocket.PingMessage, []byte("whatevs"), time.Now().Add(time.Second)); err != nil {
					logger.Debugln("error writing ping:", err)
				}
			case <-t.quit:
				return
			}
		}
	}()
	t.conn = conn
	t.eid = eid
	return nil
}

func (t *Tendermint) clearTransactions() {
	// receive events on the websocket
	// and clear transactions from the txPool
	// if t.requestCleared == true, fire on t.cleared when all txs in pool are confirmed

	// Read message
	// errors are logged and ignored
	for {
		_, p, err := t.conn.ReadMessage()
		if err != nil {
			logger.Debugln(fmt.Errorf("Error reading ws: %v\n", err))
			// TODO: the connection was probably lost or something.
			// we should try and renew it ...
			continue
		}

		var response struct {
			Result struct {
				Event string      `json:"event"`
				Data  interface{} `json:"data"`
			} `json:"result"`
			Error string `json:"error"`
		}

		wire.ReadJSON(&response, p, &err)
		if err != nil {
			logger.Debugln(fmt.Errorf("error unmarshaling event data: %v", err))
			continue
		}

		if response.Error != "" {
			logger.Debugln(fmt.Errorf("response error: %v", response.Error))
			continue
		}

		if response.Result.Event != t.eid {
			logger.Debugf("received unsolicited event! Got %s, expected %s\n", response.Result.Event, t.eid)
			continue
		}

		var txid []byte
		// TODO: get txid from response.Result.Data

		// TODO: keep return values and exceptions from call txs

		txidString := fmt.Sprintf("%X", txid)

		// if the tx is known, clear it
		t.mtx.Lock()
		if _, ok := t.txPool[txidString]; ok {
			delete(t.txPool, txidString)
		}
		if t.requestCleared && len(t.txPool) == 0 {
			t.cleared <- struct{}{}
		}
		t.mtx.Unlock()
	}
}
