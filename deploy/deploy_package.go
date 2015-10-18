package deploy

import (
	"github.com/eris-ltd/eris-pm/definitions"
	"github.com/eris-ltd/eris-pm/loaders"
)

func Deploy(do *definitions.Do) error {
	// check if Do struct has Package added. If not then load.
	if do.Package == nil {
		var err error

		do.Package, err = loaders.LoadPackage(do.YAMLPath)
		if err != nil {
			return err
		}
	}

	// run jobs

	return nil
}