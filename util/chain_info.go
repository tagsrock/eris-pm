package util

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/eris-ltd/eris-pm/definitions"

	log "github.com/eris-ltd/eris-logger"

	"github.com/eris-ltd/eris-db/client"
)

func ChainStatus(do *definitions.Do) (chainId string, latestBlockHeight int, err error) {
	nodeClient := client.NewErisNodeClient(do.Chain)
	// NOTE: ChainId should be bytes, convert it here to string
	// NodeInfo is no longer exposed through Status();
	// other values are currentlu not used by e-pm
	chainIdBytes, _, _, latestBlockHeight, _, err := nodeClient.Status()
	if err != nil {
		return "", 0, err
	}
	// set return values
	chainId = string(chainIdBytes)
	return
}



func GetChainID(do *definitions.Do) error {
	if do.ChainID == "" {
		chainId, _, err := ChainStatus(do)
		if err != nil {
			return err
		}

		do.ChainID = chainId
		log.WithField("=>", do.ChainID).Info("Using ChainID from Node")
	}

	return nil
}

func AccountsInfo(account, field string, do *definitions.Do) (string, error) {

	addrBytes, err := hex.DecodeString(account)
	if err != nil {
		return "", fmt.Errorf("Account Addr %s is improper hex: %v", account, err)
	}
	nodeClient := client.NewErisNodeClient(do.Chain)
	r, err := nodeClient.GetAccount(addrBytes)
	if err != nil {
		return "", err
	}
	if r == nil {
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

		s, err = FormatOutput([]string{"permissions"}, 0, r)

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
		s, err = FormatOutput([]string{field}, 0, r)
	}

	if err != nil {
		return "", err
	}

	return s, nil
}

func NamesInfo(name, field string, do *definitions.Do) (string, error) {
	nodeClient := client.NewErisNodeClient(do.Chain)
	owner, data, expirationBlock, err := nodeClient.GetName(name)
	if err != nil {
		return "", err
	}
	if r == nil {
		return "", fmt.Errorf("Name %s does not exist", name)
	}

	switch field {
	case "name":
		return name, nil
	case "owner":
		return string(owner), nil
	case "data":
		return data, nil
	case "expires":
		return strvonv.Itoa(expirationBlock), nil
	default:
		return "", fmt.Errorf("Field %s not recognized", field)
	}
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
