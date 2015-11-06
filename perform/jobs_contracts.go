package perform

import (
	"encoding/hex"
	"fmt"
	"github.com/eris-ltd/eris-pm/definitions"
	"github.com/eris-ltd/eris-pm/util"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/eris-ltd/eris-pm/Godeps/_workspace/src/github.com/eris-ltd/common/go/common"
	"github.com/eris-ltd/eris-pm/Godeps/_workspace/src/github.com/eris-ltd/lllc-server"
	"github.com/eris-ltd/eris-pm/Godeps/_workspace/src/github.com/eris-ltd/mint-client/mintx/core"
	"github.com/eris-ltd/eris-pm/Godeps/_workspace/src/github.com/tendermint/tendermint/types"
)

func PackageDeployJob(pkgDeploy *definitions.PackageDeploy, do *definitions.Do) (string, error) {
	// todo
	var result string

	return result, nil
}

func DeployJob(deploy *definitions.Deploy, do *definitions.Do) (string, error) {
	// Preprocess variables
	deploy.Source, _ = util.PreProcess(deploy.Source, do)
	deploy.Contract, _ = util.PreProcess(deploy.Contract, do)
	deploy.Amount, _ = util.PreProcess(deploy.Amount, do)
	deploy.Nonce, _ = util.PreProcess(deploy.Nonce, do)
	deploy.Fee, _ = util.PreProcess(deploy.Fee, do)
	deploy.Gas, _ = util.PreProcess(deploy.Gas, do)

	// Use default
	deploy.Source = useDefault(deploy.Source, do.Package.Account)
	deploy.Amount = useDefault(deploy.Amount, do.DefaultAmount)
	deploy.Fee = useDefault(deploy.Fee, do.DefaultFee)
	deploy.Gas = useDefault(deploy.Gas, do.DefaultGas)

	// assemble contract
	var p string
	if _, err := os.Stat(deploy.Contract); err == nil {
		p = deploy.Contract
	} else {
		p = filepath.Join(do.ContractsPath, deploy.Contract)
	}
	logger.Debugf("Contract path =>\t\t%s\n", p)

	// compile
	bytecode, abiSpec, err := lllcserver.Compile(p)
	if err != nil {
		return "", err
	}
	logger.Debugf("Abi spec =>\t\t\t%s\n", string(abiSpec))
	contractCode := hex.EncodeToString(bytecode)

	// Don't use pubKey if account override
	var oldKey string
	if deploy.Source != do.Package.Account {
		oldKey = do.PublicKey
		do.PublicKey = ""
	}

	// additional data may be sent along with the contract
	// these are naively added to the end of the contract code using standard
	// mint packing
	if deploy.Data != "" {
		splitout := strings.Split(deploy.Data, " ")
		for _, s := range splitout {
			s, _ = util.PreProcess(s, do)
			addOns := common.LeftPadString(common.StripHex(common.Coerce2Hex(s)), 64)
			logger.Debugf("Contract Code =>\t\t%s\n", contractCode)
			logger.Debugf("\tAdditional Data =>\t%s\n", addOns)
			contractCode = contractCode + addOns
		}
	}

	// Deploy contract
	logger.Infof("Deploying Contract =>\t\t%s:%v\n", deploy.Source, contractCode)
	tx, err := core.Call(do.Chain, do.Signer, do.PublicKey, deploy.Source, "", deploy.Amount, deploy.Nonce, deploy.Gas, deploy.Fee, contractCode)
	if err != nil {
		return "", fmt.Errorf("Error deploying contract %s: %v", p, err)
	}

	// Sign, broadcast, display
	var result string
	result, err = deployFinalize(do, tx, deploy.Wait)

	// Save ABI
	if _, err := os.Stat(do.ABIPath); os.IsNotExist(err) {
		if err := os.Mkdir(do.ABIPath, 0775); err != nil {
			return "", err
		}
	}
	abiLocation := filepath.Join(do.ABIPath, result)
	logger.Debugf("Saving ABI =>\t\t\t%s\n", abiLocation)
	if err := ioutil.WriteFile(abiLocation, []byte(abiSpec), 0664); err != nil {
		return "", err
	}

	// Don't use pubKey if account override
	if deploy.Source != do.Package.Account {
		do.PublicKey = oldKey
	}

	return result, nil
}

func CallJob(call *definitions.Call, do *definitions.Do) (string, error) {
	// Preprocess variables
	call.Source, _ = util.PreProcess(call.Source, do)
	call.Destination, _ = util.PreProcess(call.Destination, do)
	call.Amount, _ = util.PreProcess(call.Amount, do)
	call.Nonce, _ = util.PreProcess(call.Nonce, do)
	call.Fee, _ = util.PreProcess(call.Fee, do)
	call.Gas, _ = util.PreProcess(call.Gas, do)

	// Use default
	call.Source = useDefault(call.Source, do.Package.Account)
	call.Amount = useDefault(call.Amount, do.DefaultAmount)
	call.Fee = useDefault(call.Fee, do.DefaultFee)
	call.Gas = useDefault(call.Gas, do.DefaultGas)

	var err error
	call.Data, err = util.ReadAbiFormulateCall(call.Destination, call.Data, do)
	if err != nil {
		return "", err
	}

	// Don't use pubKey if account override
	var oldKey string
	if call.Source != do.Package.Account {
		oldKey = do.PublicKey
		do.PublicKey = ""
	}

	logger.Infof("Calling =>\t\t\t%s:%v\n", call.Destination, call.Data)
	tx, err := core.Call(do.Chain, do.Signer, do.PublicKey, call.Source, call.Destination, call.Amount, call.Nonce, call.Gas, call.Fee, call.Data)
	if err != nil {
		return "", err
	}

	// Don't use pubKey if account override
	if call.Source != do.Package.Account {
		do.PublicKey = oldKey
	}

	// Sign, broadcast, display
	return txFinalize(do, tx, call.Wait)
}

func deployFinalize(do *definitions.Do, tx interface{}, wait bool) (string, error) {
	var result string

	res, err := core.SignAndBroadcast(do.ChainID, do.Chain, do.Signer, tx.(types.Tx), true, true, wait)
	if err != nil {
		logger.Errorf("ERROR =>\n")
		return "", err
	}

	if err := util.ReadTxSignAndBroadcast(res, err); err != nil {
		logger.Errorf("ERROR =>\n")
		return "", err
	}

	result = fmt.Sprintf("%X", res.Address)
	return result, nil
}
