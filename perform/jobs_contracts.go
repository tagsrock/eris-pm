package perform

import (
	"github.com/eris-ltd/eris-pm/definitions"
	"github.com/eris-ltd/eris-pm/util"

	"github.com/eris-ltd/eris-pm/Godeps/_workspace/src/github.com/eris-ltd/mint-client/mintx/core"
)

func DeployJob(deploy *definitions.Deploy, do *definitions.Do) (string, error) {
	// todo
	var result string

	return result, nil
}

func PackageDeployJob(pkgDeploy *definitions.PackageDeploy, do *definitions.Do) (string, error) {
	// todo
	var result string

	return result, nil
}

func CallJob(call *definitions.Call, do *definitions.Do) (string, error) {
	var result string

	logger.Infof("Calling =>\t\t\t%s:%v", call.Destination, call.Data)
	tx, err := core.Call(do.Chain, do.Signer, do.PublicKey, call.Source, call.Destination, call.Amount, call.Nonce, call.Gas, call.Gas, call.Data)
	if err != nil {
		return "", err
	}

	if err := util.UnpackSignAndBroadcast(core.SignAndBroadcast(do.ChainID, do.Chain, do.Signer, tx, true, true, call.Wait)); err != nil {
		return "", err
	}
	return result, nil
}
