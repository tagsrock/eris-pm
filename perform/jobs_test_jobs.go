package perform

import (
	"encoding/hex"
	"fmt"
	"strconv"

	"github.com/eris-ltd/eris-pm/definitions"
	"github.com/eris-ltd/eris-pm/util"

	log "github.com/eris-ltd/eris-pm/Godeps/_workspace/src/github.com/Sirupsen/logrus"
	cclient "github.com/eris-ltd/eris-pm/Godeps/_workspace/src/github.com/eris-ltd/tendermint/rpc/core_client"
)

func QueryContractJob(query *definitions.QueryContract, do *definitions.Do) (string, error) {
	// Preprocess variables. We don't preprocess data as it is processed by ReadAbiFormulateCall
	query.Source, _ = util.PreProcess(query.Source, do)
	query.Destination, _ = util.PreProcess(query.Destination, do)
	query.ABI, _ = util.PreProcess(query.ABI, do)

	// Set the from and the to
	fromAddrBytes, err := hex.DecodeString(query.Source)
	if err != nil {
		return "", err
	}
	toAddrBytes, err := hex.DecodeString(query.Destination)
	if err != nil {
		return "", err
	}

	// Get the packed data from the ABI functions
	var data string
	if query.ABI == "" {
		data, err = util.ReadAbiFormulateCall(query.Destination, query.Data, do)
	} else {
		data, err = util.ReadAbiFormulateCall(query.ABI, query.Data, do)
	}
	if err != nil {
		return util.ABIErrorHandler(do, err, nil, query)
	}
	dataBytes, err := hex.DecodeString(data)
	if err != nil {
		return "", err
	}

	// Call the client
	client := cclient.NewClient(do.Chain, "HTTP")
	retrn, err := client.Call(fromAddrBytes, toAddrBytes, dataBytes)
	if err != nil {
		return "", err
	}

	// Preprocess return
	result, err := util.FormatOutput([]string{"return"}, 0, retrn)
	if err != nil {
		return "", err
	}
	result, err = strconv.Unquote(result)
	if err != nil {
		return "", err
	}

	// Formally process the return
	log.WithField("res", result).Debug("Decoding Raw Result")
	if query.ABI == "" {
		log.WithField("abi", query.Destination).Debug()
		result, err = util.ReadAndDecodeContractReturn(query.Destination, query.Data, result, do)
	} else {
		log.WithField("abi", query.ABI).Debug()
		result, err = util.ReadAndDecodeContractReturn(query.ABI, query.Data, result, do)
	}
	if err != nil {
		return "", err
	}

	// Finalize
	if result != "" {
		log.WithField("=>", result).Warn("Return Value")
	} else {
		log.Debug("No return.")
	}
	return result, nil
}

func QueryAccountJob(query *definitions.QueryAccount, do *definitions.Do) (string, error) {
	// Preprocess variables
	query.Account, _ = util.PreProcess(query.Account, do)
	query.Field, _ = util.PreProcess(query.Field, do)

	// Perform Query
	arg := fmt.Sprintf("%s:%s", query.Account, query.Field)
	log.WithField("=>", arg).Info("Querying Account")

	result, err := util.AccountsInfo(query.Account, query.Field, do)
	if err != nil {
		return "", err
	}

	// Result
	if result != "" {
		log.WithField("=>", result).Warn("Return Value")
	} else {
		log.Debug("No return.")
	}
	return result, nil
}

func QueryNameJob(query *definitions.QueryName, do *definitions.Do) (string, error) {
	// Preprocess variables
	query.Name, _ = util.PreProcess(query.Name, do)
	query.Field, _ = util.PreProcess(query.Field, do)

	// Peform query
	log.WithFields(log.Fields{
		"name":  query.Name,
		"field": query.Field,
	}).Info("Querying")
	result, err := util.NamesInfo(query.Name, query.Field, do)
	if err != nil {
		return "", err
	}

	if result != "" {
		log.WithField("=>", result).Warn("Return Value")
	} else {
		log.Debug("No return.")
	}
	return result, nil
}

func QueryValsJob(query *definitions.QueryVals, do *definitions.Do) (string, error) {
	var result string

	// Preprocess variables
	query.Field, _ = util.PreProcess(query.Field, do)

	// Peform query
	log.WithField("=>", query.Field).Info("Querying Vals")
	result, err := util.ValidatorsInfo(query.Field, do)
	if err != nil {
		return "", err
	}

	if result != "" {
		log.WithField("=>", result).Warn("Return Value")
	} else {
		log.Debug("No return.")
	}
	return result, nil
}

func AssertJob(assertion *definitions.Assert, do *definitions.Do) (string, error) {
	var result string
	// Preprocess variables
	assertion.Key, _ = util.PreProcess(assertion.Key, do)
	assertion.Relation, _ = util.PreProcess(assertion.Relation, do)
	assertion.Value, _ = util.PreProcess(assertion.Value, do)

	// Switch on relation
	log.WithFields(log.Fields{
		"key":      assertion.Key,
		"relation": assertion.Relation,
		"value":    assertion.Value,
	}).Info("Assertion =>")

	switch assertion.Relation {
	case "==", "eq":
		if assertion.Key == assertion.Value {
			return assertPass("==", assertion.Key, assertion.Value)
		} else {
			return assertFail("==", assertion.Key, assertion.Value)
		}
	case "!=", "ne":
		if assertion.Key != assertion.Value {
			return assertPass("!=", assertion.Key, assertion.Value)
		} else {
			return assertFail("!=", assertion.Key, assertion.Value)
		}
	case ">", "gt":
		k, v, err := bulkConvert(assertion.Key, assertion.Value)
		if err != nil {
			return convFail()
		}
		if k > v {
			return assertPass(">", assertion.Key, assertion.Value)
		} else {
			return assertFail(">", assertion.Key, assertion.Value)
		}
	case ">=", "ge":
		k, v, err := bulkConvert(assertion.Key, assertion.Value)
		if err != nil {
			return convFail()
		}
		if k >= v {
			return assertPass(">=", assertion.Key, assertion.Value)
		} else {
			return assertFail(">=", assertion.Key, assertion.Value)
		}
	case "<", "lt":
		k, v, err := bulkConvert(assertion.Key, assertion.Value)
		if err != nil {
			return convFail()
		}
		if k < v {
			return assertPass("<", assertion.Key, assertion.Value)
		} else {
			return assertFail("<", assertion.Key, assertion.Value)
		}
	case "<=", "le":
		k, v, err := bulkConvert(assertion.Key, assertion.Value)
		if err != nil {
			return convFail()
		}
		if k <= v {
			return assertPass("<=", assertion.Key, assertion.Value)
		} else {
			return assertFail("<=", assertion.Key, assertion.Value)
		}
	}

	return result, nil
}

func bulkConvert(key, value string) (int, int, error) {
	k, err := strconv.Atoi(key)
	if err != nil {
		return 0, 0, err
	}
	v, err := strconv.Atoi(value)
	if err != nil {
		return 0, 0, err
	}
	return k, v, nil
}

func assertPass(typ, key, val string) (string, error) {
	log.WithField("=>", fmt.Sprintf("%s %s %s", key, typ, val)).Warn("Assertion Succeeded")
	return "passed", nil
}

func assertFail(typ, key, val string) (string, error) {
	log.WithField("=>", fmt.Sprintf("%s %s %s", key, typ, val)).Warn("Assertion Failed")
	return "failed", fmt.Errorf("assertion failed")
}

func convFail() (string, error) {
	return "", fmt.Errorf("The Key of your assertion cannot be converted into an integer.\nFor string conversions please use the equal or not equal relations.")
}
