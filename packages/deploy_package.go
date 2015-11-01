package packages

import (
	"github.com/eris-ltd/eris-pm/definitions"
	"github.com/eris-ltd/eris-pm/perform"
	"github.com/eris-ltd/eris-pm/util"
)

func Deploy(do *definitions.Do) error {
	var err error

	// Load the package if it doesn't exist (as will happen with testing)
	if do.Package == nil {
		do.Package, err = LoadPackage(do.YAMLPath)
		if err != nil {
			return err
		}
		util.BundleHttpPathCorrect(do)
		util.PrintPathPackage(do)
	}

	return perform.RunDeployJobs(do)
}
