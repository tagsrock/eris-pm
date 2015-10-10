package main

import (
	"github.com/eris-ltd/eris-pm/commands"

	"github.com/eris-ltd/eris-pm/Godeps/_workspace/src/github.com/codegangsta/cli"
	"github.com/eris-ltd/eris-pm/Godeps/_workspace/src/github.com/eris-ltd/common/go/common"
)

// wraps a epm-go/commands function in a closure that accepts cli.Context
func cliCall(f func(*commands.Context) error) func(*cli.Context) {
	return func(c *cli.Context) {
		c2 := commands.TransformContext(c)
		common.IfExit(f(c2))
	}
}

var (
	setCmd = cli.Command{
		Name:   "set",
		Usage:  "set an epm var",
		Action: cliCall(commands.Set),
	}

	deployCmd = cli.Command{
		Name:   "deploy",
		Usage:  "deploy a .pdx file onto a blockchain",
		Action: cliCall(commands.Deploy),
		Flags: []cli.Flag{
			// chainFlag,
			diffFlag,
			dontClearFlag,
			contractPathFlag,
		},
	}

	plopCmd = cli.Command{
		Name:   "plop",
		Usage:  "machine readable variable display: epm plop <addr | chainid | config | genesis | key | pid | vars>",
		Action: cliCall(commands.Plop),
		Flags:  []cli.Flag{
		// chainFlag,
		},
	}

	/*

		headCmd = cli.Command{
			Name:   "head",
			Usage:  "display the current working blockchain",
			Action: cliCall(commands.Head),
		}

		commandCmd = cli.Command{
			Name:   "cmd",
			Usage:  "run a command (useful when combined with RPC): epm cmd <deploy contract.lll>",
			Action: cliCall(commands.Command),
			Flags: []cli.Flag{
				chainFlag,
				multiFlag,
				contractPathFlag,
			},
		consoleCmd = cli.Command{
			Name:   "console",
			Usage:  "run epm in interactive mode",
			Action: cliCall(commands.Console),
			Flags: []cli.Flag{
				chainFlag,
				multiFlag,
				diffFlag,
				dontClearFlag,
				contractPathFlag,
			},
		}
		headCmd = cli.Command{
			Name:   "head",
			Usage:  "display the current working blockchain",
			Action: cliCall(commands.Head),
		}

		commandCmd = cli.Command{
			Name:   "cmd",
			Usage:  "run a command (useful when combined with RPC): epm cmd <deploy contract.lll>",
			Action: cliCall(commands.Command),
			Flags: []cli.Flag{
				chainFlag,
				multiFlag,
				contractPathFlag,
			},

		testCmd = cli.Command{
			Name:   "test",
			Usage:  "run all pdx/pdt in the directory",
			Action: cliCall(commands.Test),
			Flags: []cli.Flag{
				chainFlag,
				contractPathFlag,
			},
		}
	*/
)
