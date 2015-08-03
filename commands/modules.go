package commands

import (
	"fmt"

	//epm-binary-generator:IMPORT
	mod "github.com/eris-ltd/eris-pm/commands/modules/tendermint"

	"github.com/eris-ltd/eris-pm/epm"
)

// this needs to match the type of the chain we're trying to run
// it should be blank for the base epm (even though it includes thel ...)
//epm-binary-generator:CHAIN
const CHAIN = ""

// chainroot is a full path to the dir
func LoadChain(c *Context, chainType string) (epm.ChainClient, error) {
	logger.Debugln("Loading chain ", c.String("type"))

	// TODO: these things from flags

	// need a toml config

	rpcAddr := c.String("node-addr")
	if rpcAddr == "" {
		rpcHost, rpcPort := c.String("host"), c.Int("port")
		rpcAddr = fmt.Sprintf("http://%s:%d/", rpcHost, rpcPort)
	}
	signHost, signPort := c.String("sign-addr-host"), c.Int("sign-addr-port")

	signAddr := fmt.Sprintf("http://%s:%d", signHost, signPort)
	chainID, pubkey := c.String("chainid"), c.String("pubkey")

	chain, err := mod.NewChain(chainID, rpcAddr, signAddr, pubkey)
	if err != nil {
		return nil, err
	}

	// TODO: set contract path?

	return chain, nil
}
