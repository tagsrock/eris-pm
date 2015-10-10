package main

import (
	"os"

	"github.com/eris-ltd/eris-pm/version"

	"github.com/eris-ltd/eris-pm/Godeps/_workspace/src/github.com/codegangsta/cli"
	"github.com/eris-ltd/eris-pm/Godeps/_workspace/src/github.com/eris-ltd/common/go/log"
)

func main() {

	app := cli.NewApp()
	app.Name = "epm"
	app.Usage = "The Eris Package Manager Tests and Deploys Smart Contract Engines"
	app.Version = version.VERSION
	app.Author = "Ethan Buchman"
	app.Email = "ethan@erisindustries.com"

	app.Before = before
	app.After = after

	app.Flags = []cli.Flag{
		// // which chain
		// chainFlag,
		// chainIDFlag,
		// typeFlag,

		// log
		logLevelFlag,

		// rpc
		rpcHostFlag,
		rpcPortFlag,
		rpcAddrFlag,

		// key server
		signPortFlag,
		signHostFlag,

		// languages
		compilerFlag,

		// pubkey
		pubkeyFlag,

		// files
		contractPathFlag,
		pdxPathFlag,
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
