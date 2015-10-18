package commands

import (
	"github.com/eris-ltd/eris-pm/deploy"

	. "github.com/eris-ltd/eris-pm/Godeps/_workspace/src/github.com/eris-ltd/common/go/common"
	"github.com/eris-ltd/eris-pm/Godeps/_workspace/src/github.com/spf13/cobra"
)

//----------------------------------------------------

// Primary Deploy Sub-Command
var Deploy = &cobra.Command{
	Use:   "deploy",
	Short: "Deploy your smart contract package.",
	Long:  `Deploy your smart contract package.`,
	Run:   DeployPackage,
}

// build the deploy subcommand
func buildDeployCommand() {
	addDeployFlags()
}

//----------------------------------------------------

func addDeployFlags() {}

//----------------------------------------------------
func DeployPackage(cmd *cobra.Command, args []string) {
	IfExit(deploy.Deploy(do))
}
