package commands

import (

	//epm-binary-generator:IMPORT
	mod "github.com/eris-ltd/eris-pm/commands/modules/tendermint"

	"github.com/eris-ltd/eris-pm/epm"
)

// this needs to match the type of the chain we're trying to run
// it should be blank for the base epm (even though it includes thel ...)
//epm-binary-generator:CHAIN
const CHAIN = ""

// chainroot is a full path to the dir
func LoadChain(c *Context, chainType string) epm.ChainClient {
	rpc := c.Bool("rpc")
	logger.Debugln("Loading chain ", c.String("type"))

	chain := mod.NewChain(chainType, rpc)
	// XXX: setupModule ...
	applyFlags(c, chain)
	return chain
}

func applyFlags(c *Context, m epm.ChainClient) {
	// then apply flags
	//TODO:
	/*
		setLogLevel(c, m)
		setKeysFile(c, m)
		setGenesisPath(c, m)
		setContractPath(c, m)
		setMining(c, m)
		setRpc(c, m)
	*/
}
