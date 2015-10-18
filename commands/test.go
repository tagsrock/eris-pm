package commands

import (
	"github.com/eris-ltd/eris-pm/test"

	. "github.com/eris-ltd/eris-pm/Godeps/_workspace/src/github.com/eris-ltd/common/go/common"
	"github.com/eris-ltd/eris-pm/Godeps/_workspace/src/github.com/spf13/cobra"
)

//----------------------------------------------------

// Primary Test Sub-Command
var Test = &cobra.Command{
	Use:   "test",
	Short: "Test your smart contract package.",
	Long:  `Test your smart contract package.`,
	Run:   TestPackage,
}

// build the test subcommand
func buildTestCommand() {
	addTestFlags()
}

//----------------------------------------------------

func addTestFlags() {}

//----------------------------------------------------
func TestPackage(cmd *cobra.Command, args []string) {
	IfExit(test.Test(do))
}
