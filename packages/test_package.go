package packages

import (
	"github.com/eris-ltd/eris-pm/definitions"
)

func Test(do *definitions.Do) error {
	var err error

	do.Package, err = LoadPackage(do.YAMLPath)
	if err != nil {
		return err
	}

	if err := Deploy(do); err != nil {
		return err
	}

	// do some testing.

	return nil
}
