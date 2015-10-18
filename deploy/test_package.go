package deploy

import (
	"github.com/eris-ltd/eris-pm/definitions"
	"github.com/eris-ltd/eris-pm/loaders"
)

func Test(do *definitions.Do) error {
	var err error

	do.Package, err = loaders.LoadPackage(do.YAMLPath)
	if err != nil {
		return err
	}

	Deploy(do)

	// do some testing.
	return nil
}