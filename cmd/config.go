package commands

import (
	cfg "github.com/eris-ltd/eris-pm/Godeps/_workspace/src/github.com/tendermint/tendermint/config"

	_ "github.com/eris-ltd/eris-pm/Godeps/_workspace/src/github.com/eris-ltd/mint-client/utils" // calls ApplyConfig
)

var config cfg.Config = nil

func init() {
	cfg.OnConfig(func(newConfig cfg.Config) {
		config = newConfig
	})
}
