package perform

import (
	"fmt"
	"strconv"

	"github.com/eris-ltd/eris-pm/definitions"
	"github.com/eris-ltd/eris-pm/util"
)

func QueryAccountJob(query *definitions.QueryAccount, do *definitions.Do) (string, error) {
	var result string

	// Preprocess variables
	query.Account, _ = util.PreProcess(query.Account, do)
	query.Field, _ = util.PreProcess(query.Field, do)

	// Perform Query
	logger.Infof("Querying Account =>\t\t%s:%s\n", query.Account, query.Field)
	result, err := util.AccountsInfo(query.Account, query.Field, do)
	if err != nil {
		return "", err
	}

	// Result
	return result, nil
}

func QueryNameJob(query *definitions.QueryName, do *definitions.Do) (string, error) {
	var result string

	// Preprocess variables
	query.Name, _ = util.PreProcess(query.Name, do)
	query.Field, _ = util.PreProcess(query.Field, do)

	// Peform query
	logger.Infof("Querying Name =>\t\t%s:%s\n", query.Name, query.Field)
	result, err := util.NamesInfo(query.Name, query.Field, do)
	if err != nil {
		return "", err
	}

	return result, nil
}

func QueryContractJob(query *definitions.QueryContract, do *definitions.Do) (string, error) {
	var result string

	return result, nil
}

func AssertJob(assertion *definitions.Assert, do *definitions.Do) (string, error) {
	var result string
	// Preprocess variables
	assertion.Key, _ = util.PreProcess(assertion.Key, do)
	assertion.Relation, _ = util.PreProcess(assertion.Relation, do)
	assertion.Value, _ = util.PreProcess(assertion.Value, do)

	// Switch on relation
	logger.Infof("Assertion =>\t\t\t%s:%s:%s\n", assertion.Key, assertion.Relation, assertion.Value)
	switch assertion.Relation {
	case "==", "eq":
		if assertion.Key == assertion.Value {
			return assertPass()
		} else {
			return assertFail(assertion.Key, assertion.Value)
		}
	case "!=", "ne":
		if assertion.Key != assertion.Value {
			return assertPass()
		} else {
			return assertFail(assertion.Key, assertion.Value)
		}
	case ">", "gt":
		k, v, err := bulkConvert(assertion.Key, assertion.Value)
		if err != nil {
			return convFail()
		}
		if k > v {
			return assertPass()
		} else {
			return assertFail(assertion.Key, assertion.Value)
		}
	case ">=", "ge":
		k, v, err := bulkConvert(assertion.Key, assertion.Value)
		if err != nil {
			return convFail()
		}
		if k >= v {
			return assertPass()
		} else {
			return assertFail(assertion.Key, assertion.Value)
		}
	case "<", "lt":
		k, v, err := bulkConvert(assertion.Key, assertion.Value)
		if err != nil {
			return convFail()
		}
		if k < v {
			return assertPass()
		} else {
			return assertFail(assertion.Key, assertion.Value)
		}
	case "<=", "le":
		k, v, err := bulkConvert(assertion.Key, assertion.Value)
		if err != nil {
			return convFail()
		}
		if k <= v {
			return assertPass()
		} else {
			return assertFail(assertion.Key, assertion.Value)
		}
	}

	return result, nil
}

func QueryValsJob(query *definitions.QueryVals, do *definitions.Do) (string, error) {
	var result string

	// Preprocess variables
	query.Field, _ = util.PreProcess(query.Field, do)

	// Peform query
	logger.Infof("Querying Vals =>\t\t%s\n", query.Field)
	result, err := util.ValidatorsInfo(query.Field, do)
	if err != nil {
		return "", err
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

func assertPass() (string, error) {
	return "passed", nil
}

func assertFail(expect, receive string) (string, error) {
	return "failed", fmt.Errorf("Assertion Failed =>\t\t%s:%s", expect, receive)
}

func convFail() (string, error) {
	return "", fmt.Errorf("The Key of your assertion cannot be converted into an integer.\nFor string conversions please use the equal or not equal relations.")
}
