package main

import (
	"github.com/eris-ltd/eris-pm/Godeps/_workspace/src/github.com/codegangsta/cli"
	"github.com/eris-ltd/eris-pm/commands"
)

var (
	// chainFlag = cli.StringFlag{
	// 	Name:   "chain",
	// 	Value:  "",
	// 	Usage:  "set the chain by <ref name> or by <type>/<id>",
	// 	EnvVar: "",
	// }

	// typeFlag = cli.StringFlag{
	// 	Name:   "type",
	// 	Value:  "tendermint",
	// 	Usage:  "set the chain type (thelonious, genesis, bitcoin, ethereum)",
	// 	EnvVar: "",
	// }

	// interactiveFlag = cli.BoolFlag{
	// 	Name:   "i",
	// 	Usage:  "run epm in interactive mode",
	// 	EnvVar: "",
	// }

	// chainIDFlag = cli.StringFlag{
	// 	Name:  "chain_id",
	// 	Usage: "specify the chain id",
	// 	Value: "etcb_testnet",
	// }

	diffFlag = cli.BoolFlag{
		Name:   "diff",
		Usage:  "show a diff of all contract storage",
		EnvVar: "",
	}

	dontClearFlag = cli.BoolFlag{
		Name:   "dont-clear",
		Usage:  "stop epm from clearing the epm cache on startup",
		EnvVar: "",
	}

	contractPathFlag = cli.StringFlag{
		Name:  "contracts, c",
		Value: commands.DefaultContractPath,
		Usage: "set the contract path",
	}

	pdxPathFlag = cli.StringFlag{
		Name:   "p",
		Value:  ".",
		Usage:  "specify a .pdx file to deploy",
		EnvVar: "DEPLOY_PDX",
	}

	logLevelFlag = cli.IntFlag{
		Name:   "log",
		Value:  0,
		Usage:  "set the log level",
		EnvVar: "EPM_LOG",
	}

	rpcHostFlag = cli.StringFlag{
		Name:  "host",
		Value: "localhost",
		Usage: "set the rpc host",
	}

	rpcPortFlag = cli.IntFlag{
		Name:  "port",
		Value: 46657,
		Usage: "set the rpc port",
	}

	rpcAddrFlag = cli.StringFlag{
		Name:  "node-addr",
		Value: "",
		Usage: "set the full http address of the rpc node",
	}

	signPortFlag = cli.IntFlag{
		Name:  "sign_port",
		Usage: "set the port for the eris-keys server",
		Value: 4767,
	}

	signHostFlag = cli.StringFlag{
		Name:  "sign_host",
		Usage: "set the host for the eris-keys server",
		Value: "localhost",
	}

	pubkeyFlag = cli.StringFlag{
		Name:  "pubkey",
		Usage: "specify pubkey to use",
	}

	compilerFlag = cli.StringFlag{
		Name:  "compiler",
		Usage: "specify <host>:<port> to use for compile server",
	}
)
