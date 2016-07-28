package util

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/eris-ltd/eris-pm/definitions"

	log "github.com/eris-ltd/eris-logger"
	cclient "github.com/eris-ltd/eris-pm/Godeps/_workspace/src/github.com/eris-ltd/tendermint/rpc/core_client"
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

func GetChainID(do *definitions.Do) error {
	if do.ChainID == "" {
		status, err := ChainStatus("node_info", do)
		if err != nil {
			return err
		}

		// Wrangle these returns
		type NodeInfo struct {
			ChainID string `mapstructure:"chain_id" json:"chain_id"`
		}
		var ret NodeInfo
		err = json.Unmarshal([]byte(status), &ret)
		if err != nil {
			return err
		}

		do.ChainID = ret.ChainID
		log.WithField("=>", do.ChainID).Info("Using ChainID from Node")
	}

	return nil
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
