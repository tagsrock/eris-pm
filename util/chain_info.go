package util

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/eris-ltd/eris-pm/definitions"

	"github.com/eris-ltd/eris-pm/Godeps/_workspace/src/github.com/eris-ltd/mint-client/mintx/core"
	cclient "github.com/eris-ltd/eris-pm/Godeps/_workspace/src/github.com/tendermint/tendermint/rpc/core_client"
)

func ChainStatus(field string, do *definitions.Do) (string, error) {
	client := cclient.NewClient(do.Chain, "HTTP")

	r, err := client.Status()
	if err != nil {
		return "", err
	}

	s, err := FormatOutput([]string{field}, 0, r)
	if err != nil {
		return "", err
	}

	return s, nil
}

func AccountsInfo(account, field string, do *definitions.Do) (string, error) {
	client := cclient.NewClient(do.Chain, "HTTP")

	addrBytes, err := hex.DecodeString(account)
	if err != nil {
		return "", fmt.Errorf("Account Addr %s is improper hex: %v", account, err)
	}

	r, err := client.GetAccount(addrBytes)
	if err != nil {
		return "", err
	}
	if r == nil {
		return "", fmt.Errorf("Account %s does not exist", account)
	}

	r2 := r.Account
	if r2 == nil {
		return "", fmt.Errorf("Account %s does not exist", account)
	}

	var s string
	if strings.Contains(field, "permissions") {

		type BasePermission struct {
			PermissionValue int `mapstructure:"perms" json:"perms"`
			SetBitValue     int `mapstructure:"set" json:"set"`
		}

		type AccountPermission struct {
			Base  *BasePermission `mapstructure:"base" json:"base"`
			Roles []string        `mapstructure:"roles" json:"roles"`
		}

		fields := strings.Split(field, ".")

		s, err = FormatOutput([]string{"permissions"}, 0, r2)

		var deconstructed AccountPermission
		err := json.Unmarshal([]byte(s), &deconstructed)
		if err != nil {
			return "", err
		}

		if len(fields) > 1 {
			switch fields[1] {
			case "roles":
				s = strings.Join(deconstructed.Roles, ",")
			case "base", "perms":
				s = strconv.Itoa(deconstructed.Base.PermissionValue)
			case "set":
				s = strconv.Itoa(deconstructed.Base.SetBitValue)
			}
		}
	} else {
		s, err = FormatOutput([]string{field}, 0, r2)
	}

	if err != nil {
		return "", err
	}

	return s, nil
}

func NamesInfo(account, field string, do *definitions.Do) (string, error) {
	client := cclient.NewClient(do.Chain, "HTTP")

	r, err := client.GetName(account)
	if err != nil {
		return "", err
	}
	if r == nil {
		return "", fmt.Errorf("Account %s does not exist", account)
	}

	r2 := r.Entry
	s, err := FormatOutput([]string{field}, 0, r2)
	if err != nil {
		return "", err
	}

	s, err = strconv.Unquote(s)
	if err != nil {
		return "", err
	}

	return s, nil
}

func ValidatorsInfo(field string, do *definitions.Do) (string, error) {
	client := cclient.NewClient(do.Chain, "HTTP")

	r, err := client.ListValidators()
	if err != nil {
		return "", err
	}

	s, err := FormatOutput([]string{field}, 0, r)
	if err != nil {
		return "", err
	}

	type Account struct {
		Address string `mapstructure:"address" json:"address"`
	}

	var deconstructed []Account
	err = json.Unmarshal([]byte(s), &deconstructed)
	if err != nil {
		return "", err
	}

	vals := []string{}
	for _, v := range deconstructed {
		vals = append(vals, v.Address)
	}
	s = strings.Join(vals, ",")

	return s, nil
}

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
