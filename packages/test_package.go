package packages

import (
	"github.com/eris-ltd/eris-pm/definitions"
	"github.com/eris-ltd/eris-pm/perform"
	"github.com/eris-ltd/eris-pm/util"
)

func Test(do *definitions.Do) error {
	var err error

	// Load the package
	do.Package, err = LoadPackage(do.YAMLPath)
	if err != nil {
		return err
	}
	util.BundleHttpPathCorrect(do)
	util.PrintPathPackage(do)

	if err := Deploy(do); err != nil {
		return err
	}

	return perform.RunTestJobs(do)
}
