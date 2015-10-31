package packages

import (
	"github.com/eris-ltd/eris-pm/definitions"
	"github.com/eris-ltd/eris-pm/perform"
	"github.com/eris-ltd/eris-pm/util"
)

func Deploy(do *definitions.Do) error {
	// check if Do struct has Package added (via possible marshalling from Tests). If not then load.
	if do.Package == nil {
		var err error
		do.Package, err = LoadPackage(do.YAMLPath)
		if err != nil {
			return err
		}
		util.BundleHttpPathCorrect(do)
		util.PrintPathPackage(do)
	}

	perform.RunDeployJobs(do)

	return nil
}
