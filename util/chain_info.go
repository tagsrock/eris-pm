package util

import (
	"github.com/eris-ltd/eris-pm/Godeps/_workspace/src/github.com/eris-ltd/mint-client/mintx/core"
	cclient "github.com/eris-ltd/eris-pm/Godeps/_workspace/src/github.com/tendermint/tendermint/rpc/core_client"
)

func ChainStatus(nodeAddr, field string) (string, error) {
	var client cclient.Client
	client = cclient.NewClient(nodeAddr, "HTTP")

	r, err := client.Status()
	if err != nil {
		return "", err
	}

	s, err := FormatOutput([]string{field}, 0, r)
	if err != nil {
		return "", err
	}

	return s, nil
}

// This is a closer function which is called by most of the tx_run functions
func UnpackSignAndBroadcast(result *core.TxResult, err error) error {
	// if there's an error just return.
	if err != nil {
		return err
	}

	// if there is nothing to unpack then just return.
	if result == nil {
		return nil
	}

	// Unpack and display for the user.
	logger.Printf("Transaction Hash =>\t\t%X\n", result.Hash)
	if result.Address != nil {
		logger.Infof("Contract Address =>\t\t%X\n", result.Address)
	}
	if result.Return != nil {
		logger.Debugf("Block Hash =>\t\t%X\n", result.BlockHash)
		logger.Debugf("Return Value =>\t\t%X\n", result.Return)
		logger.Debugf("Exception =>\t\t\t%s\n", result.Exception)
	}

	return nil
}
