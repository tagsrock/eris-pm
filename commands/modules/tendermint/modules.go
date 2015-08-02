package commands

import (
	"fmt"
	"log"

	"github.com/eris-ltd/eris-pm/epm"
	//	mintconfig "github.com/tendermint/tendermint/config"
)

func NewChain(chainType string, rpc bool) epm.ChainClient {
	switch chainType {
	case "tendermint", "mint":
		if rpc {
			log.Fatal("Tendermint rpc not implemented yet")
		} else {
			// return mint.NewMint()
		}
	}
	return nil

}

func ChainSpecificDeploy(chain epm.ChainClient, deployGen, root string, novi bool) error {
	return nil
}

func Fetch(chainType, peerserver string) ([]byte, error) {
	return nil, fmt.Errorf("Fetch not supported for mint")
}
