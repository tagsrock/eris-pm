package main

import (
	"os"

	"github.com/eris-ltd/eris-pm/Godeps/_workspace/src/github.com/codegangsta/cli"
	"github.com/eris-ltd/eris-pm/Godeps/_workspace/src/github.com/eris-ltd/common/go/log"
)

func main() {

	app := cli.NewApp()
	app.Name = "epm"
	app.Usage = "The Eris Package Manager Tests and Operates Blockchains and Smart Contract Systems"
	app.Version = "0.10.0"
	app.Author = "Ethan Buchman"
	app.Email = "ethan@erisindustries.com"

	app.Before = before
	app.After = after

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
		setCmd,
		plopCmd,
	}

	app.Run(os.Args)
}

func before(c *cli.Context) error {
	log.SetLogLevelGlobal(log.LogLevel(c.Int("log")))
	return nil
}

func after(c *cli.Context) error {
	log.Flush()
	return nil
}
