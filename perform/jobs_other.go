package perform

import (
	"github.com/eris-ltd/eris-pm/definitions"
)

func SetAccountJob(account *definitions.Account, do *definitions.Do) (string, error) {
	var result string
	do.Package.Account = account.Address
	result = account.Address
	return result, nil
}

func SetValJob(set *definitions.Set, do *definitions.Do) (string, error) {
	var result string
	result = set.Value
	return result, nil
}

func DumpStateJob(dump *definitions.DumpState, do *definitions.Do) (string, error) {
	var result string

	return result, nil
}

func RestoreStateJob(restore *definitions.RestoreState, do *definitions.Do) (string, error) {
	var result string

	return result, nil
}
