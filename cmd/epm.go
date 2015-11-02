package commands

import (
	"os"

	"github.com/eris-ltd/eris-pm/definitions"
	"github.com/eris-ltd/eris-pm/util"
	"github.com/eris-ltd/eris-pm/version"

	"github.com/eris-ltd/eris-pm/Godeps/_workspace/src/github.com/eris-ltd/common/go/log"
	"github.com/eris-ltd/eris-pm/Godeps/_workspace/src/github.com/spf13/cobra"
	cfg "github.com/eris-ltd/eris-pm/Godeps/_workspace/src/github.com/tendermint/tendermint/config"
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

		util.ClearJobResults()

		// TODO: get chainid from node endpoint
		config.Set("chain_id", do.ChainID)
		config.Set("log_level", "error")
		cfg.ApplyConfig(config)
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
	EPMCmd.PersistentFlags().StringVarP(&do.PublicKey, "key", "k", "", "full public key to use by default")
	EPMCmd.PersistentFlags().StringVarP(&do.ChainID, "chain-id", "i", "", "identifier of the chain to work against")
	EPMCmd.PersistentFlags().UintVarP(&do.DefaultGas, "gas", "g", 1111111111, "default gas to use; can be overridden for a single job")
	EPMCmd.PersistentFlags().BoolVarP(&do.Verbose, "verbose", "v", false, "verbose output")
	EPMCmd.PersistentFlags().BoolVarP(&do.Debug, "debug", "d", false, "debug level output")
}
