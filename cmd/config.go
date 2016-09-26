package commands

import (
	cfg "github.com/eris-ltd/tendermint/config"

	_ "github.com/eris-ltd/mint-client/utils" // calls ApplyConfig
)

var config cfg.Config = nil

func init() {
	cfg.OnConfig(func(newConfig cfg.Config) {
		config = newConfig
	})
}
