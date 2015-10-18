package commands

import (
	"fmt"
	"os"

	"github.com/eris-ltd/eris-pm/version"
	"github.com/eris-ltd/eris-pm/definitions"

	"github.com/eris-ltd/eris-pm/Godeps/_workspace/src/github.com/eris-ltd/common/go/log"
	"github.com/eris-ltd/eris-pm/Godeps/_workspace/src/github.com/spf13/cobra"
)

const VERSION = version.VERSION

// Global Do struct
var do *definitions.Do

// Defining the root command
var EPMCmd = &cobra.Command{
	Use:   "epm [command] [flags]",
	Short: "The Eris Package Manager Tests and Deploys Smart Contract Systems",
	Long: `The Eris Package Manager Tests and Deploys Smart Contract Systems

Made with <3 by Eris Industries.

Complete documentation is available at https://docs.erisindustries.com
` + "\nVersion:\n  " + VERSION,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		var logLevel log.LogLevel
		if do.Verbose {
			logLevel = 2
		} else if do.Debug {
			logLevel = 3
		}
		log.SetLoggers(logLevel, os.Stdout, os.Stderr) // TODO: make this better....
	},
	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		log.Flush()
	},
}

func Execute() {
	do = definitions.NowDo()
	AddGlobalFlags()
	AddCommands()
	EPMCmd.Execute()
}

// Define the commands
func AddCommands() {
	buildTestCommand()
	EPMCmd.AddCommand(Test)
	buildDeployCommand()
	EPMCmd.AddCommand(Deploy)
}

// Flags that are to be used by commands are handled by the Do struct
// Define the persistent commands (globals)
func AddGlobalFlags() {
	EPMCmd.PersistentFlags().StringVarP(&do.YAMLPath, "file", "f", "./epm.yaml", "path to package file which EPM should use")
	EPMCmd.PersistentFlags().StringVarP(&do.ContractsPath, "contracts-path", "p", ".", "path to the contracts EPM should use")
	EPMCmd.PersistentFlags().StringVarP(&do.Chain, "chain", "c", "chain:46657", "<ip:port> of chain which EPM should use")
	EPMCmd.PersistentFlags().StringVarP(&do.Signer, "sign", "s", "keys:4767", "<ip:port> of signer daemon which EPM should use")
	EPMCmd.PersistentFlags().StringVarP(&do.Compiler, "compiler", "m", "compilers.eris.industries:8091", "<ip:port> of compiler which EPM should use")
	EPMCmd.PersistentFlags().UintVarP(&do.DefaultGas, "gas", "g", 1111111111, "default gas to use; can be overridden for a single job")
	EPMCmd.PersistentFlags().BoolVarP(&do.Verbose, "verbose", "v", false, "verbose output")
	EPMCmd.PersistentFlags().BoolVarP(&do.Debug, "debug", "d", false, "debug level output")
}

func ArgCheck(num int, comp string, cmd *cobra.Command, args []string) error {
	switch comp {
	case "eq":
		if len(args) != num {
			cmd.Help()
			return fmt.Errorf("\n**Note** you sent our marmots the wrong number of arguments.\nPlease send the marmots %d arguments only.", num)
		}
	case "ge":
		if len(args) < num {
			cmd.Help()
			return fmt.Errorf("\n**Note** you sent our marmots the wrong number of arguments.\nPlease send the marmots at least %d argument(s).", num)
		}
	}
	return nil
}
