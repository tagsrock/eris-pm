package packages

import (
	"github.com/eris-ltd/eris-pm/definitions"
	"github.com/eris-ltd/eris-pm/perform"
)

func Test(do *definitions.Do) error {
	var err error

	// Load the package
	do.Package, err = LoadPackage(do.YAMLPath)
	if err != nil {
		return err
	}

	if err := Deploy(do); err != nil {
		return err
	}

	return perform.RunTestJobs(do)
}
