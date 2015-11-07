package perform

import (
	"github.com/eris-ltd/eris-pm/definitions"
	"github.com/eris-ltd/eris-pm/util"

	keys "github.com/eris-ltd/eris-pm/Godeps/_workspace/src/github.com/eris-ltd/eris-keys/eris-keys"
)

func SetAccountJob(account *definitions.Account, do *definitions.Do) (string, error) {
	var result string
	var err error

	// Preprocess
	account.Address, _ = util.PreProcess(account.Address, do)

	// Set the Account in the Package & Announce
	do.Package.Account = account.Address
	logger.Infof("Setting Account =>\t\t%s\n", do.Package.Account)

	// Set the public key from eris-keys
	keys.DaemonAddr = do.Signer
	logger.Infof("Getting Public Key =>\t\t%s\n", keys.DaemonAddr)
	do.PublicKey, err = keys.Call("pub", map[string]string{"addr": do.Package.Account, "name": ""})
	if _, ok := err.(keys.ErrConnectionRefused); ok {
		keys.ExitConnectErr(err)
	}

	if err != nil {
		return "", err
	}

	// Set result and return
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
