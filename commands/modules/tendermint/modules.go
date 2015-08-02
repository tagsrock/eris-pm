package commands

import (
	"log"

	"github.com/eris-ltd/eris-pm/epm"

	"github.com/eris-ltd/eris-pm/Godeps/_workspace/src/github.com/eris-ltd/mint-client/mintx/core"
)

var (
	SIGN      = true
	BROADCAST = true
	WAIT      = false
)

func NewChain(chainType string, rpc bool) epm.ChainClient {
	switch chainType {
	case "tendermint", "mint":
		if rpc {
			log.Fatal("Tendermint rpc not implemented yet")
		} else {
			return &Tendermint{} // TODO!
		}
	}
	return nil

}

func NewTendermint(chainID, nodeAddr, signAddr, pubkey, addr string) *Tendermint {
	return &Tendermint{
		chainID:  chainID,
		nodeAddr: nodeAddr,
		signAddr: signAddr,
		pubkey:   pubkey,
		addr:     addr,
	}
}

type Tendermint struct {
	chainID  string
	nodeAddr string
	signAddr string
	pubkey   string
	addr     string
}

func (t *Tendermint) Tx(addr, amt string) (string, error) {
	tx, err := core.Send(t.nodeAddr, t.pubkey, t.addr, addr, amt, "")
	if err != nil {
		return "", err
	}
	result, err = core.SignAndBroadcast(t.chainID, t.nodeAddr, t.signAddr, tx, SIGN, BROADCAST, WAIT)
	return fmt.Sprintf("%X", result.Hash), err
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
