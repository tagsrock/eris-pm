package util

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/eris-ltd/eris-pm/definitions"

	"github.com/eris-ltd/eris-pm/Godeps/_workspace/src/github.com/eris-ltd/common/go/common"
	ebi "github.com/eris-ltd/eris-pm/Godeps/_workspace/src/github.com/eris-ltd/eris-abi/core"
	"github.com/eris-ltd/eris-pm/Godeps/_workspace/src/github.com/eris-ltd/mint-client/mintx/core"
)

// This is a closer function which is called by most of the tx_run functions
func ReadTxSignAndBroadcast(result *core.TxResult, err error) error {
	// if there's an error just return.
	if err != nil {
		return err
	}

	// if there is nothing to unpack then just return.
	if result == nil {
		return nil
	}

	// Unpack and display for the user.
	logger.Printf("Transaction Hash =>\t\t%X\n", result.Hash)
	if result.Address != nil {
		logger.Infof("Contract Address =>\t\t%X\n", result.Address)
	}
	if result.Return != nil {
		logger.Debugf("Block Hash =>\t\t\t%X\n", result.BlockHash)
		logger.Debugf("Return Value =>\t\t\t%X\n", result.Return)
		logger.Debugf("Exception =>\t\t\t%s\n", result.Exception)
	}

	return nil
}

func ReadAbiFormulateCall(contract, dataRaw string, do *definitions.Do) (string, error) {
	abiSpecBytes, err := readAbi(do.ABIPath, contract)
	if err != nil {
		return "", err
	}
	logger.Debugf("ABI Spec =>\t\t\t%s\n", abiSpecBytes)

	// Process and Pack the Call
	funcName, args := abiPreProcess(dataRaw, do)

	var totalArgs []string
	totalArgs = append(totalArgs, funcName)
	totalArgs = append(totalArgs, args...)
	logger.Debugf("Packing Call via ABI =>\t\t%v:%v\n", funcName, args)

	return ebi.Packer(abiSpecBytes, totalArgs...)
}

func ReadAndDecodeContractReturn(contract, dataRaw, resultRaw string, do *definitions.Do) (string, error) {
	abiSpecBytes, err := readAbi(do.ABIPath, contract)
	if err != nil {
		return "", err
	}
	logger.Debugf("ABI Spec =>\t\t\t%s\n", abiSpecBytes)

	// Process and Pack the Call
	funcName, _ := abiPreProcess(dataRaw, do)

	// Unpack the result
	res, err := ebi.UnPacker(abiSpecBytes, funcName, resultRaw, false)
	if err != nil {
		return "", err
	}

	// Wrangle these returns
	type ContractReturn struct {
		Name  string `mapstructure:"," json:","`
		Type  string `mapstructure:"," json:","`
		Value string `mapstructure:"," json:","`
	}
	var resTotal []ContractReturn
	err = json.Unmarshal([]byte(res), &resTotal)
	if err != nil {
		return "", err
	}

	// Get the value
	result := resTotal[0].Value
	return result, nil
}

func abiPreProcess(dataRaw string, do *definitions.Do) (string, []string) {
	var dataNew []string

	data := strings.Split(dataRaw, " ")
	for _, d := range data {
		d, _ = PreProcess(d, do)
		dataNew = append(dataNew, d)
	}

	funcName := dataNew[0]
	args := dataNew[1:]

	return funcName, args
}

func readAbi(root, contract string) ([]byte, error) {
	p := path.Join(root, common.StripHex(contract))
	if _, err := os.Stat(p); err != nil {
		return []byte{}, fmt.Errorf("Abi doesn't exist for =>\t%s", p)
	}

	b, err := ioutil.ReadFile(p)
	if err != nil {
		return []byte{}, err
	}

	return b, nil
}
