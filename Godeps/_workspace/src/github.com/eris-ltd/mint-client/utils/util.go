package utils

import (
	cfg "github.com/eris-ltd/eris-pm/Godeps/_workspace/src/github.com/tendermint/tendermint/config"
	tmcfg "github.com/eris-ltd/eris-pm/Godeps/_workspace/src/github.com/tendermint/tendermint/config/tendermint"
)

func init() {
	cfg.ApplyConfig(tmcfg.GetConfig(""))
}
