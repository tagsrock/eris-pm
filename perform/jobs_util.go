package perform

import (
	"github.com/eris-ltd/eris-pm/definitions"
	"github.com/eris-ltd/eris-pm/util"
)

func SetAccountJob(account *definitions.Account, do *definitions.Do) (string, error) {
	var result string
	account.Address, _ = util.PreProcess(account.Address, do)
	do.Package.Account = account.Address
	logger.Infof("Setting Account =>\t\t%s\n", do.Package.Account)
	result = account.Address
	return result, nil
}

func SetValJob(set *definitions.Set, do *definitions.Do) (string, error) {
	var result string
	set.Value, _ = util.PreProcess(set.Value, do)
	logger.Infof("Setting Variable =>\t\t%s\n", set.Value)
	result = set.Value
	return result, nil
}
