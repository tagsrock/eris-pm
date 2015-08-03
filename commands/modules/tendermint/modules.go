package commands

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"net/http"
	"path"
	"strings"
	"sync"
	"time"

	"github.com/eris-ltd/eris-pm/epm"

	"github.com/eris-ltd/eris-pm/Godeps/_workspace/src/github.com/eris-ltd/mint-client/mintx/core"
	"github.com/eris-ltd/eris-pm/Godeps/_workspace/src/github.com/gorilla/websocket"
	"github.com/eris-ltd/eris-pm/Godeps/_workspace/src/github.com/tendermint/tendermint/account"
	cclient "github.com/eris-ltd/eris-pm/Godeps/_workspace/src/github.com/tendermint/tendermint/rpc/core_client"
	rpcserver "github.com/eris-ltd/eris-pm/Godeps/_workspace/src/github.com/tendermint/tendermint/rpc/server"
	rpctypes "github.com/eris-ltd/eris-pm/Godeps/_workspace/src/github.com/tendermint/tendermint/rpc/types"
	"github.com/eris-ltd/eris-pm/Godeps/_workspace/src/github.com/tendermint/tendermint/types"
	"github.com/eris-ltd/eris-pm/Godeps/_workspace/src/github.com/tendermint/tendermint/wire"
)

// switches for SignAndBroadcast
var (
	SIGN      = true
	BROADCAST = true
	WAIT      = false // we implement waiting here so we only subscribe once on big pdxs
)

func NewChain(rootDir, chainID, nodeAddr, signAddr, pubkey string) (epm.ChainClient, error) {
	return NewTendermint(rootDir, chainID, nodeAddr, signAddr, pubkey)
}

// implements epm.ChainClient
type Tendermint struct {
	chainID string
	rootDir string

	// hosts
	nodeAddr string
	wsAddr   string
	signAddr string

	// pubkey and address
	pubkey         string
	inputAddr      string
	inputAddrBytes []byte

	quit chan struct{}

	// txs are added to the txPool after they are broadcast
	// they should be removed once they are confirmed in a block
	mtx    sync.Mutex
	txPool map[string]struct{}

	// Commit() blocks waiting for the txPool to empty.
	// When len(txPool) == 0, fire on cleared
	cleared        chan struct{}
	requestCleared bool // = true when we call Commit()

	// we subscribe to our address
	eid  string          // event id
	conn *websocket.Conn // ws conn
}

func NewTendermint(rootDir, chainID, nodeAddr, signAddr, pubkey string) (*Tendermint, error) {

	//TODO: get addr from pubkey!
	pubkeyBytes, err := hex.DecodeString(pubkey)
	if err != nil {
		return nil, fmt.Errorf("Pubkey is invalid hex: %v", err)
	}

	if len(pubkeyBytes) != 32 {
		return nil, fmt.Errorf("Pubkey is wrong length. Got %d, expected 32", len(pubkeyBytes))
	}
	var pubKey account.PubKeyEd25519
	copy(pubKey[:], pubkeyBytes)
	inputAddrBytes := pubKey.Address()
	inputAddr := hex.EncodeToString(inputAddrBytes)

	wsAddr := strings.TrimPrefix(nodeAddr, "http://")
	wsAddr = fmt.Sprintf("ws://%swebsocket", wsAddr)
	t := &Tendermint{
		chainID:        chainID,
		rootDir:        rootDir,
		nodeAddr:       nodeAddr,
		wsAddr:         wsAddr,
		signAddr:       signAddr,
		pubkey:         pubkey,
		inputAddr:      inputAddr,
		inputAddrBytes: inputAddrBytes,
		txPool:         make(map[string]struct{}),
		quit:           make(chan struct{}),
		cleared:        make(chan struct{}),
	}

	logger.Debugf("Starting tendermint client with chainID:%s, rootDir:%s, nodeAddr:%s, wsAddr:%s, signAddr:%s, pubkey:%s\n", t.chainID, t.rootDir, t.nodeAddr, t.wsAddr, t.signAddr, t.pubkey)

	// setup a websocket connection
	// and subscribe to inputs from the sending address
	if err := t.setupWSConn(); err != nil {
		return nil, err
	}

	// waits for transactions sent by this client to be confirmed
	// and clears them from the txPool
	// Commit() blocks waiting for clearTransactions to clear the pool
	go t.clearTransactions()

	return t, nil
}

func (t *Tendermint) RootDir() string {
	return t.rootDir
}

func (t *Tendermint) AddTxToPool(txHash string) {
	t.mtx.Lock()
	t.txPool[txHash] = struct{}{}
	t.mtx.Unlock()
}

// A simple SendTx
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
	t.AddTxToPool(hashString)
	return hashString, nil
}

// CallTx to a contract with some data
func (t *Tendermint) Msg(addr string, data string) (string, error) {
	//nodeAddr, pubkey, fromAddr, toAddr, amt, nonce, gas, fee, data
	tx, err := core.Call(t.nodeAddr, t.pubkey, t.inputAddr, addr, "1", "", "1000", "0", data)
	if err != nil {
		return "", err
	}
	result, err := core.SignAndBroadcast(t.chainID, t.nodeAddr, t.signAddr, tx, SIGN, BROADCAST, WAIT)
	if err != nil {
		return "", err
	}
	hashString := fmt.Sprintf("%X", result.Hash)
	t.AddTxToPool(hashString)
	return hashString, nil
}

// Simulated call to contract with some data
func (t *Tendermint) Call(addr string, data string) (string, error) {
	client := cclient.NewClient(t.nodeAddr, "HTTP")
	addrBytes, err := hex.DecodeString(addr)
	if err != nil {
		return "", err
	}
	dataBytes, err := hex.DecodeString(data)
	if err != nil {
		return "", err
	}
	r, err := client.Call(t.inputAddrBytes, addrBytes, dataBytes)
	if err != nil {
		return "", nil
	}
	v := hex.EncodeToString(r.Return)
	return v, nil
}

// Create new contract
func (t *Tendermint) Script(code string) (string, string, error) {
	tx, err := core.Call(t.nodeAddr, t.pubkey, t.inputAddr, "", "1", "", "10000", "0", code)
	if err != nil {
		return "", "", err
	}
	result, err := core.SignAndBroadcast(t.chainID, t.nodeAddr, t.signAddr, tx, SIGN, BROADCAST, WAIT)
	if err != nil {
		return "", "", err
	}
	hashString := fmt.Sprintf("%X", result.Hash)
	addrString := fmt.Sprintf("%X", result.Address)
	t.AddTxToPool(hashString)
	return hashString, addrString, nil
}

func (t *Tendermint) NameReg(name, value string) (string, error) {
	return "", nil

}

func (t *Tendermint) StorageAt(target, storage string) (string, error) {
	client := cclient.NewClient(t.nodeAddr, "HTTP")
	targetBytes, err := hex.DecodeString(target)
	if err != nil {
		return "", err
	}
	storageBytes, err := hex.DecodeString(storage)
	if err != nil {
		return "", err
	}
	r, err := client.GetStorage(targetBytes, storageBytes)
	if err != nil {
		return "", nil
	}
	v := hex.EncodeToString(r.Value)
	return v, nil
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

// dial the ws host, subscribe to inputs from our address
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

// wait for events on the ws, clear transactions from the txPool accordingly
// if t.requestCleared == true, fire on t.cleared when all txs in pool are confirmed
func (t *Tendermint) clearTransactions() {
	// Read messages off the ws. Errors are logged and ignored
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
				Event string   `json:"event"`
				Data  types.Tx `json:"data"`
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

		txid := types.TxID(t.chainID, response.Result.Data)

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

		logger.Debugln("cleared tx", txidString)

		// TODO: right the tx hash and block number to file
	}
}

func (t *Tendermint) WriteConfig() error {
	buf := new(bytes.Buffer)
	buf.WriteString("\n\n")
	buf.WriteString(fmt.Sprintf("chain_id = \"%s\" \n", t.chainID))
	buf.WriteString(fmt.Sprintf("root_dir = \"%s\" \n", t.rootDir))
	buf.WriteString(fmt.Sprintf("node_addr = \"%s\" \n", t.nodeAddr))
	buf.WriteString(fmt.Sprintf("sign_addr = \"%s\" \n", t.signAddr))
	buf.WriteString(fmt.Sprintf("pubkey = \"%s\" \n", t.pubkey))
	return ioutil.WriteFile(path.Join(t.rootDir, "config.toml"), buf.Bytes(), 0644)
}
