package packages

import (
	"github.com/eris-ltd/eris-pm/definitions"
	"github.com/eris-ltd/eris-pm/perform"
)

func Deploy(do *definitions.Do) error {
	var err error

	// Load the package if it doesn't exist (as will happen with testing)
	if do.Package == nil {
		do.Package, err = LoadPackage(do.YAMLPath)
		if err != nil {
			return err
		}
	}

	return perform.RunDeployJobs(do)
}
