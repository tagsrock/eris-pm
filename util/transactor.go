package util

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/eris-ltd/eris-pm/Godeps/_workspace/src/github.com/eris-ltd/mint-client/mintx/core"
	"github.com/eris-ltd/eris-pm/Godeps/_workspace/src/github.com/tendermint/tendermint/wire"
)

// This is a closer function which is called by most of the tx_run functions
func UnpackSignAndBroadcast(result *core.TxResult, err error) error {
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
		logger.Debugf("Block Hash =>\t\t%X\n", result.BlockHash)
		logger.Debugf("Return Value =>\t\t%X\n", result.Return)
		logger.Debugf("Exception =>\t\t\t%s\n", result.Exception)
	}

	return nil
}

func FormatOutput(args []string, i int, o interface{}) (string, error) {
	if len(args) < i+1 {
		return prettyPrint(o)
	}
	arg0 := args[i]
	v := reflect.ValueOf(o).Elem()
	name, err := fieldFromTag(v, arg0)
	if err != nil {
		return "", err
	}
	f := v.FieldByName(name)
	return prettyPrint(f.Interface())
}

func fieldFromTag(v reflect.Value, field string) (string, error) {
	iv := v.Interface()
	st := reflect.TypeOf(iv)
	for i := 0; i < v.NumField(); i++ {
		tag := st.Field(i).Tag.Get("json")
		if tag == field {
			return st.Field(i).Name, nil
		}
	}
	return "", fmt.Errorf("Invalid field name")
}

func prettyPrint(o interface{}) (string, error) {
	var prettyJSON bytes.Buffer
	err := json.Indent(&prettyJSON, wire.JSONBytes(o), "", "\t")
	if err != nil {
		return "", err
	}
	return string(prettyJSON.Bytes()), nil
}
