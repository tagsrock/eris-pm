package main

import (
	"os"

	"github.com/eris-ltd/eris-pm/Godeps/_workspace/src/github.com/codegangsta/cli"
)

func main() {

	app := cli.NewApp()
	app.Name = "epm"
	app.Usage = "The Eris Package Manager Tests and Operates Blockchains and Smart Contract Systems"
	app.Version = "0.10.0"
	app.Author = "Ethan Buchman"
	app.Email = "ethan@erisindustries.com"

	app.Flags = []cli.Flag{
		// which chain
		chainFlag,

		// log
		logLevelFlag,

		// rpc
		rpcHostFlag,
		rpcPortFlag,
		rpcLocalFlag,

		// languages
		compilerFlag,
	}

	app.Commands = []cli.Command{
		deployCmd,
	}

	app.Run(os.Args)

}
