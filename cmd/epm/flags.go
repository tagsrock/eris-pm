package main

import (
	"github.com/codegangsta/cli"

	"github.com/eris-ltd/eris-pm/commands"
)

var (
	chainFlag = cli.StringFlag{
		Name:   "chain",
		Value:  "",
		Usage:  "set the chain by <ref name> or by <type>/<id>",
		EnvVar: "",
	}

	typeFlag = cli.StringFlag{
		Name:   "type",
		Value:  "thelonious",
		Usage:  "set the chain type (thelonious, genesis, bitcoin, ethereum)",
		EnvVar: "",
	}

	interactiveFlag = cli.BoolFlag{
		Name:   "i",
		Usage:  "run epm in interactive mode",
		EnvVar: "",
	}

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
		Value:  2,
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
		Value: 5,
		Usage: "set the rpc port",
	}

	rpcLocalFlag = cli.BoolFlag{
		Name:  "local",
		Usage: "let the rpc server handle keys (sign txs)",
	}

	compilerFlag = cli.StringFlag{
		Name:  "compiler",
		Usage: "specify <host>:<port> to use for compile server",
	}
)
