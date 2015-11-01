package perform

import (
	"github.com/eris-ltd/eris-pm/definitions"
	"github.com/eris-ltd/eris-pm/util"
)

func SetAccountJob(account *definitions.Account, do *definitions.Do) (string, error) {
	var result string
	account.Address, _ = util.PreProcess(account.Address, do)

	do.Package.Account = account.Address
	result = account.Address
	return result, nil
}

func SetValJob(set *definitions.Set, do *definitions.Do) (string, error) {
	var result string
	set.Value, _ = util.PreProcess(set.Value, do)

	result = set.Value
	return result, nil
}
