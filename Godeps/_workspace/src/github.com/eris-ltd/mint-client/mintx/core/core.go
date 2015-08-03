package core

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/eris-ltd/eris-pm/Godeps/_workspace/src/github.com/gorilla/websocket"

	"github.com/eris-ltd/eris-pm/Godeps/_workspace/src/github.com/tendermint/tendermint/account"
	ptypes "github.com/eris-ltd/eris-pm/Godeps/_workspace/src/github.com/tendermint/tendermint/permission/types"
	rtypes "github.com/eris-ltd/eris-pm/Godeps/_workspace/src/github.com/tendermint/tendermint/rpc/core/types"
	cclient "github.com/eris-ltd/eris-pm/Godeps/_workspace/src/github.com/tendermint/tendermint/rpc/core_client"
	rpcserver "github.com/eris-ltd/eris-pm/Godeps/_workspace/src/github.com/tendermint/tendermint/rpc/server"
	rpctypes "github.com/eris-ltd/eris-pm/Godeps/_workspace/src/github.com/tendermint/tendermint/rpc/types"
	"github.com/eris-ltd/eris-pm/Godeps/_workspace/src/github.com/tendermint/tendermint/types"
	"github.com/eris-ltd/eris-pm/Godeps/_workspace/src/github.com/tendermint/tendermint/wire"
)

//------------------------------------------------------------------------------------
// core functions with string args.
// validates strings and forms transaction

func Output(addr, amtS string) ([]byte, error) {
	if amtS == "" {
		return nil, fmt.Errorf("output must specify an amount with the --amt flag")
	}

	if addr == "" {
		return nil, fmt.Errorf("output must specify an addr with the --addr flag")
	}

	addrBytes, err := hex.DecodeString(addr)
	if err != nil {
		return nil, fmt.Errorf("addr is bad hex: %v", err)
	}

	amt, err := strconv.ParseInt(amtS, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("amt is misformatted: %v", err)
	}
	// TODO: validate amt!

	txOutput := types.TxOutput{
		Address: addrBytes,
		Amount:  amt,
	}

	n, errPtr := new(int64), new(error)
	buf := new(bytes.Buffer)
	txOutput.WriteSignBytes(buf, n, errPtr)
	if *errPtr != nil {
		return nil, *errPtr
	}
	return buf.Bytes(), nil
}

func Input(nodeAddr, pubkey, amtS, nonceS, addr string) ([]byte, error) {
	pub, addrBytes, amt, nonce, err := checkCommon(nodeAddr, pubkey, addr, amtS, nonceS)
	if err != nil {
		return nil, err
	}

	txInput := types.TxInput{
		Address:  addrBytes,
		Amount:   amt,
		Sequence: int(nonce),
		PubKey:   pub,
	}

	n, errPtr := new(int64), new(error)
	buf := new(bytes.Buffer)
	txInput.WriteSignBytes(buf, n, errPtr)
	if *errPtr != nil {
		return nil, *errPtr
	}
	return buf.Bytes(), nil
}

func Send(nodeAddr, pubkey, addr, toAddr, amtS, nonceS string) (*types.SendTx, error) {
	pub, addrBytes, amt, nonce, err := checkCommon(nodeAddr, pubkey, addr, amtS, nonceS)
	if err != nil {
		return nil, err
	}

	if toAddr == "" {
		return nil, fmt.Errorf("destination address must be given with --to flag")
	}

	toAddrBytes, err := hex.DecodeString(toAddr)
	if err != nil {
		return nil, fmt.Errorf("toAddr is bad hex: %v", err)
	}

	tx := types.NewSendTx()
	_ = addrBytes // TODO!
	tx.AddInputWithNonce(pub, amt, int(nonce))
	tx.AddOutput(toAddrBytes, amt)

	return tx, nil
}

func Call(nodeAddr, pubkey, addr, toAddr, amtS, nonceS, gasS, feeS, data string) (*types.CallTx, error) {
	pub, _, amt, nonce, err := checkCommon(nodeAddr, pubkey, addr, amtS, nonceS)
	if err != nil {
		return nil, err
	}

	toAddrBytes, err := hex.DecodeString(toAddr)
	if err != nil {
		return nil, fmt.Errorf("toAddr is bad hex: %v", err)
	}

	fee, err := strconv.ParseInt(feeS, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("fee is misformatted: %v", err)
	}

	gas, err := strconv.ParseInt(gasS, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("gas is misformatted: %v", err)
	}

	dataBytes, err := hex.DecodeString(data)
	if err != nil {
		return nil, fmt.Errorf("data is bad hex: %v", err)
	}

	tx := types.NewCallTxWithNonce(pub, toAddrBytes, dataBytes, amt, gas, fee, int(nonce))
	return tx, nil
}

func Name(nodeAddr, pubkey, addr, amtS, nonceS, feeS, name, data string) (*types.NameTx, error) {
	pub, _, amt, nonce, err := checkCommon(nodeAddr, pubkey, addr, amtS, nonceS)
	if err != nil {
		return nil, err
	}

	fee, err := strconv.ParseInt(feeS, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("fee is misformatted: %v", err)
	}

	tx := types.NewNameTxWithNonce(pub, name, data, amt, fee, int(nonce))
	return tx, nil
}

func Permissions(nodeAddr, pubkey, addrS, nonceS, permFunc string, argsS []string) (*types.PermissionsTx, error) {
	pub, _, _, nonce, err := checkCommon(nodeAddr, pubkey, addrS, "0", "0")
	if err != nil {
		return nil, err
	}
	var args ptypes.PermArgs
	switch permFunc {
	case "set_base":
		addr, pF, err := decodeAddressPermFlag(argsS[0], argsS[1])
		if err != nil {
			return nil, err
		}
		if len(argsS) != 3 {
			return nil, fmt.Errorf("set_base also takes a value (true or false)")
		}
		var value bool
		if argsS[2] == "true" {
			value = true
		} else if argsS[2] == "false" {
			value = false
		} else {
			return nil, fmt.Errorf("Unknown value %s", argsS[2])
		}
		args = &ptypes.SetBaseArgs{addr, pF, value}
	case "unset_base":
		addr, pF, err := decodeAddressPermFlag(argsS[0], argsS[1])
		if err != nil {
			return nil, err
		}
		args = &ptypes.UnsetBaseArgs{addr, pF}
	case "set_global":
		pF, err := ptypes.PermStringToFlag(argsS[0])
		if err != nil {
			return nil, err
		}
		var value bool
		if argsS[1] == "true" {
			value = true
		} else if argsS[1] == "false" {
			value = false
		} else {
			return nil, fmt.Errorf("Unknown value %s", argsS[1])
		}
		args = &ptypes.SetGlobalArgs{pF, value}
	case "add_role":
		addr, err := hex.DecodeString(argsS[0])
		if err != nil {
			return nil, err
		}
		args = &ptypes.AddRoleArgs{addr, argsS[1]}
	case "rm_role":
		addr, err := hex.DecodeString(argsS[0])
		if err != nil {
			return nil, err
		}
		args = &ptypes.RmRoleArgs{addr, argsS[1]}
	default:
		return nil, fmt.Errorf("Invalid permission function for use in PermissionsTx: %s", permFunc)
	}
	// args := snativeArgs(
	tx := types.NewPermissionsTxWithNonce(pub, args, int(nonce))
	return tx, nil
}

func decodeAddressPermFlag(addrS, permFlagS string) (addr []byte, pFlag ptypes.PermFlag, err error) {
	if addr, err = hex.DecodeString(addrS); err != nil {
		return
	}
	if pFlag, err = ptypes.PermStringToFlag(permFlagS); err != nil {
		return
	}
	return
}

type NameGetter struct {
	client cclient.Client
}

func (n NameGetter) GetNameRegEntry(name string) *types.NameRegEntry {
	entry, err := n.client.GetName(name)
	if err != nil {
		panic(err)
	}
	return entry
}

/*
func coreNewAccount(nodeAddr, pubkey, chainID string) (*types.NewAccountTx, error) {
	pub, _, _, _, err := checkCommon(nodeAddr, pubkey, "", "0", "0")
	if err != nil {
		return nil, err
	}

	client := cclient.NewClient(nodeAddr, "HTTP")
	return types.NewNewAccountTx(NameGetter{client}, pub, chainID)
}
*/

func Bond(nodeAddr, pubkey, unbondAddr, amtS, nonceS string) (*types.BondTx, error) {
	pub, addrBytes, amt, nonce, err := checkCommon(nodeAddr, pubkey, "", amtS, nonceS)
	if err != nil {
		return nil, err
	}

	if unbondAddr == "" {
		return nil, fmt.Errorf("Unbond address must be given with --unbond-to flag")
	}

	unbondAddrBytes, err := hex.DecodeString(unbondAddr)
	if err != nil {
		return nil, fmt.Errorf("unbondAddr is bad hex: %v", err)
	}

	tx, err := types.NewBondTx(pub)
	if err != nil {
		return nil, err
	}
	_ = addrBytes
	tx.AddInputWithNonce(pub, amt, int(nonce))
	tx.AddOutput(unbondAddrBytes, amt)

	return tx, nil
}

func Unbond(addrS, heightS string) (*types.UnbondTx, error) {
	if addrS == "" {
		return nil, fmt.Errorf("Validator address must be given with --addr flag")
	}

	addrBytes, err := hex.DecodeString(addrS)
	if err != nil {
		return nil, fmt.Errorf("addr is bad hex: %v", err)
	}

	height, err := strconv.ParseInt(heightS, 10, 32)
	if err != nil {
		return nil, fmt.Errorf("height is misformatted: %v", err)
	}

	return &types.UnbondTx{
		Address: addrBytes,
		Height:  int(height),
	}, nil
}

func Rebond(addrS, heightS string) (*types.RebondTx, error) {
	if addrS == "" {
		return nil, fmt.Errorf("Validator address must be given with --addr flag")
	}

	addrBytes, err := hex.DecodeString(addrS)
	if err != nil {
		return nil, fmt.Errorf("addr is bad hex: %v", err)
	}

	height, err := strconv.ParseInt(heightS, 10, 32)
	if err != nil {
		return nil, fmt.Errorf("height is misformatted: %v", err)
	}

	return &types.RebondTx{
		Address: addrBytes,
		Height:  int(height),
	}, nil
}

//------------------------------------------------------------------------------------
// sign and broadcast

func Sign(signBytes, signAddr, signRPC string) (sig [64]byte, err error) {
	args := map[string]string{
		"hash": signBytes,
		"addr": signAddr,
	}
	b, err := json.Marshal(args)
	if err != nil {
		return
	}
	logger.Debugln("Sending request body:", string(b))
	req, err := http.NewRequest("POST", signRPC+"/sign", bytes.NewBuffer(b))
	if err != nil {
		return
	}
	req.Header.Add("Content-Type", "application/json")
	sigS, errS, err := requestResponse(req)
	if err != nil {
		return sig, fmt.Errorf("Error calling signing daemon: %s", err.Error())
	}
	if errS != "" {
		return sig, fmt.Errorf("Error (string) calling signing daemon: %s", errS)
	}
	sigBytes, err := hex.DecodeString(sigS)
	if err != nil {
		return
	}
	copy(sig[:], sigBytes)
	return
}

func Broadcast(tx types.Tx, broadcastRPC string) (*rtypes.Receipt, error) {
	client := cclient.NewClient(broadcastRPC, "JSONRPC")
	rec, err := client.BroadcastTx(tx)
	if err != nil {
		return nil, err
	}
	return rec, nil
}

//------------------------------------------------------------------------------------
// utils for talking to the key server

type HTTPResponse struct {
	Response string
	Error    string
}

func requestResponse(req *http.Request) (string, string, error) {
	client := new(http.Client)
	resp, err := client.Do(req)
	if err != nil {
		return "", "", err
	}
	if resp.StatusCode >= 400 {
		return "", "", fmt.Errorf(resp.Status)
	}
	return unpackResponse(resp)
}

func unpackResponse(resp *http.Response) (string, string, error) {
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", "", err
	}
	r := new(HTTPResponse)
	if err := json.Unmarshal(b, r); err != nil {
		return "", "", err
	}
	return r.Response, r.Error, nil
}

//------------------------------------------------------------------------------------
// sign and broadcast convenience

// tx has either one input or we default to the first one (ie for send/bond)
// TODO: better support for multisig..
func signTx(signAddr, chainID string, tx_ types.Tx) ([]byte, types.Tx, error) {
	signBytes := fmt.Sprintf("%X", account.SignBytes(chainID, tx_))
	var inputAddr []byte
	var sigED account.SignatureEd25519
	switch tx := tx_.(type) {
	case *types.SendTx:
		inputAddr = tx.Inputs[0].Address
		defer func(s *account.SignatureEd25519) { tx.Inputs[0].Signature = *s }(&sigED)
	case *types.NameTx:
		inputAddr = tx.Input.Address
		defer func(s *account.SignatureEd25519) { tx.Input.Signature = *s }(&sigED)
	case *types.CallTx:
		inputAddr = tx.Input.Address
		defer func(s *account.SignatureEd25519) { tx.Input.Signature = *s }(&sigED)
	case *types.PermissionsTx:
		inputAddr = tx.Input.Address
		defer func(s *account.SignatureEd25519) { tx.Input.Signature = *s }(&sigED)
	case *types.BondTx:
		inputAddr = tx.Inputs[0].Address
		defer func(s *account.SignatureEd25519) { tx.Inputs[0].Signature = *s }(&sigED)
	case *types.UnbondTx:
		inputAddr = tx.Address
		defer func(s *account.SignatureEd25519) { tx.Signature = *s }(&sigED)
	case *types.RebondTx:
		inputAddr = tx.Address
		defer func(s *account.SignatureEd25519) { tx.Signature = *s }(&sigED)
	}
	addrHex := fmt.Sprintf("%X", inputAddr)
	sig, err := Sign(signBytes, addrHex, signAddr)
	if err != nil {
		return nil, nil, err
	}
	sigED = account.SignatureEd25519(sig)
	logger.Debugf("SIG: %X\n", sig)
	return inputAddr, tx_, nil
}

type TxResult struct {
	Hash []byte // all txs get a hash

	// only CallTx
	Address   []byte // only for new contracts
	Return    []byte
	Exception string

	//TODO: make Broadcast() errors more responsive so we
	// can differentiate mempool errors from other
}

func SignAndBroadcast(chainID, nodeAddr, signAddr string, tx types.Tx, sign, broadcast, wait bool) (txResult *TxResult, err error) {
	var inputAddr []byte
	if sign {
		inputAddr, tx, err = signTx(signAddr, chainID, tx)
		if err != nil {
			return nil, err
		}
	}

	if broadcast {
		if wait {
			var ch chan Msg
			ch, err = subscribeAndWait(tx, chainID, nodeAddr, inputAddr)
			if err != nil {
				return nil, err
			} else {
				defer func() {
					if err != nil {
						// if broadcast threw an error, just return
						return
					}
					logger.Debugln("Waiting for tx to be committed ...")
					msg := <-ch
					if msg.Error != nil {
						logger.Infoln("Encountered error waiting for event: %v\n", msg.Error)
						err = msg.Error
					} else {
						txResult.Return = msg.Value
						txResult.Exception = msg.Exception
					}
				}()
			}
		}
		var receipt *rtypes.Receipt
		receipt, err = Broadcast(tx, nodeAddr)
		if err != nil {
			return nil, err
		}
		txResult = &TxResult{
			Hash: receipt.TxHash,
		}
		if tx_, ok := tx.(*types.CallTx); ok {
			if len(tx_.Address) == 0 {
				txResult.Address = types.NewContractAddress(tx_.Input.Address, tx_.Input.Sequence)
			}
		}
	}
	return
}

//------------------------------------------------------------------------------------
// wait for events

type Msg struct {
	Value     []byte
	Exception string
	Error     error
}

func subscribeAndWait(tx types.Tx, chainID, nodeAddr string, inputAddr []byte) (chan Msg, error) {
	// subscribe to event and wait for tx to be committed
	wsAddr := strings.TrimPrefix(nodeAddr, "http://")
	wsAddr = "ws://" + wsAddr + "websocket"
	logger.Debugln(wsAddr)
	dialer := websocket.DefaultDialer
	rHeader := http.Header{}
	conn, _, err := dialer.Dial(wsAddr, rHeader)
	if err != nil {
		return nil, fmt.Errorf("Error establishing websocket connection to wait for tx to get committed: %v", err)
	}
	eid := types.EventStringAccInput(inputAddr)

	if err := conn.WriteJSON(rpctypes.RPCRequest{
		JSONRPC: "2.0",
		Id:      "",
		Method:  "subscribe",
		Params:  []interface{}{eid},
	}); err != nil {
		return nil, fmt.Errorf("Error subscribing to AccInput event: %v", err)
	}

	go func() {
		pingTicker := time.NewTicker((time.Second * rpcserver.WSReadTimeoutSeconds) / 2)
		for {
			select {
			case <-pingTicker.C:
				if err := conn.WriteControl(websocket.PingMessage, []byte("whatevs"), time.Now().Add(time.Second)); err != nil {
					logger.Debugln("error writing ping:", err)
				}
			}
		}
	}()

	resultChan := make(chan Msg, 1)

	// Read message
	go func() {
		for {
			_, p, err := conn.ReadMessage()
			if err != nil {
				resultChan <- Msg{Error: err}
				break
			} else {
				var response struct {
					Result struct {
						Event string               `json:"event"`
						Data  types.EventMsgCallTx `json:"data"`
					} `json:"result"`
					Error string `json:"error"`
				}
				err := new(error)
				wire.ReadJSON(&response, p, err)
				if *err != nil {
					resultChan <- Msg{Error: fmt.Errorf("error unmarshaling event data: %v", *err)}
					return
				}
				if response.Error != "" {
					resultChan <- Msg{Error: fmt.Errorf("response error: %v", response.Error)}
					return
				}
				if response.Result.Event != eid {
					logger.Debugf("received unsolicited event! Got %s, expected %s\n", response.Result.Event, eid)
					continue
				}
				if !bytes.Equal(types.TxID(chainID, response.Result.Data.Tx), types.TxID(chainID, tx)) {
					logger.Debugf("Received event for same input from another transaction: %X\n", types.TxID(chainID, response.Result.Data.Tx))
					continue
				}

				if response.Result.Data.Exception != "" {
					resultChan <- Msg{Value: response.Result.Data.Return, Exception: response.Result.Data.Exception}
					return
				}

				// GOOD!
				resultChan <- Msg{Value: response.Result.Data.Return}
				return
			}
		}
	}()

	// txs should take no more than 10 seconds
	timeoutTicker := time.Tick(10 * time.Second)

	go func() {
		<-timeoutTicker
		resultChan <- Msg{Error: fmt.Errorf("timed out waiting for event")}
		return
	}()
	return resultChan, nil
}

//------------------------------------------------------------------------------------
// convenience function

func checkCommon(nodeAddr, pubkey, addr, amtS, nonceS string) (pub account.PubKey, addrBytes []byte, amt int64, nonce int64, err error) {
	if amtS == "" {
		err = fmt.Errorf("input must specify an amount with the --amt flag")
		return
	}

	if pubkey == "" && addr == "" {
		err = fmt.Errorf("at least one of --pubkey or --addr must be given")
		return
	}

	pubKeyBytes, err := hex.DecodeString(pubkey)
	if err != nil {
		err = fmt.Errorf("pubkey is bad hex: %v", err)
		return
	}

	addrBytes, err = hex.DecodeString(addr)
	if err != nil {
		err = fmt.Errorf("addr is bad hex: %v", err)
		return
	}

	amt, err = strconv.ParseInt(amtS, 10, 64)
	if err != nil {
		err = fmt.Errorf("amt is misformatted: %v", err)
	}

	if len(pubKeyBytes) > 0 {
		var pubArray [32]byte
		copy(pubArray[:], pubKeyBytes)
		pub = account.PubKeyEd25519(pubArray)
		addrBytes = pub.Address()
	}

	if nonceS == "" {
		if nodeAddr == "" {
			err = fmt.Errorf("input must specify a nonce with the --nonce flag or use --node-addr (or MINTX_NODE_ADDR) to fetch the nonce from a node")
			return
		}

		// fetch nonce from node
		client := cclient.NewClient(nodeAddr, "HTTP")
		var ac *account.Account
		ac, err = client.GetAccount(addrBytes)
		if err != nil {
			err = fmt.Errorf("Error connecting to node (%s) to fetch nonce: %s", nodeAddr, err.Error())
			return
		}
		nonce = int64(ac.Sequence) + 1
	} else {
		nonce, err = strconv.ParseInt(nonceS, 10, 64)
		if err != nil {
			err = fmt.Errorf("nonce is misformatted: %v", err)
			return
		}
	}

	return
}
