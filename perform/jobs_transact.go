package perform

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"

	"github.com/eris-ltd/eris-pm/definitions"
	"github.com/eris-ltd/eris-pm/util"

	"github.com/eris-ltd/eris-pm/Godeps/_workspace/src/github.com/eris-ltd/mint-client/mintx/core"
	"github.com/eris-ltd/eris-pm/Godeps/_workspace/src/github.com/tendermint/tendermint/types"
)

func SendJob(send *definitions.Send, do *definitions.Do) (string, error) {
	// Process Variables
	send.Source, _ = util.PreProcess(send.Source, do)
	send.Destination, _ = util.PreProcess(send.Destination, do)
	send.Amount, _ = util.PreProcess(send.Amount, do)

	// Use Default
	send.Source = useDefault(send.Source, do.Package.Account)

	// Formulate tx
	logger.Infof("Sending Transaction =>\t\t%s:%s:%s\n", send.Source, send.Destination, send.Amount)
	tx, err := core.Send(do.Chain, do.Signer, do.PublicKey, send.Source, send.Destination, send.Amount, send.Nonce)
	if err != nil {
		logger.Errorf("ERROR =>\n")
		return "", err
	}

	// Sign, broadcast, display
	return txFinalize(do, tx, send.Wait)
}

func RegisterNameJob(name *definitions.RegisterName, do *definitions.Do) (string, error) {
	// Process Variables
	name.DataFile, _ = util.PreProcess(name.DataFile, do)

	// If a data file is given it should be in csv format and
	// it will be read first. Once the file is parsed and sent
	// to the chain then a single nameRegTx will be sent if that
	// has been populated.
	if name.DataFile != "" {
		// open the file and use a reader
		fileReader, err := os.Open(name.DataFile)
		if err != nil {
			logger.Errorf("ERROR =>\n")
			return "", err
		}

		defer fileReader.Close()
		r := csv.NewReader(fileReader)

		// loop through the records
		for {
			// Read the record
			record, err := r.Read()

			// Catch the errors
			if err == io.EOF {
				break
			}
			if err != nil {
				logger.Errorf("ERROR =>\n")
				return "", err
			}

			// Sink the Amount into the third slot in the record if
			// it doesn't exist
			if len(record) <= 2 {
				record = append(record, name.Amount)
			}

			// Send an individual Tx for the record
			// [TODO]: move these to async using goroutines?
			r, err := registerNameTx(&definitions.RegisterName{
				Source: name.Source,
				Name:   record[0],
				Data:   record[1],
				Amount: record[2],
				Fee:    name.Fee,
				Nonce:  name.Nonce,
				Wait:   name.Wait,
			}, do)

			if err != nil {
				logger.Errorf("ERROR =>\n")
				return "", err
			}

			n := fmt.Sprintf("%s:%s", record[0], record[1])
			if err = util.WriteJobResult(n, r); err != nil {
				logger.Errorf("ERROR =>\n")
				return "", err
			}
		}
	}

	// If the data field is populated then there is a single
	// nameRegTx to send. So do that *now*.
	if name.Data != "" {
		return registerNameTx(name, do)
	} else {
		return "data_file_parsed", nil
	}
}

// Runs an individual nametx.
func registerNameTx(name *definitions.RegisterName, do *definitions.Do) (string, error) {
	// Process Variables
	name.Source, _ = util.PreProcess(name.Source, do)
	name.Name, _ = util.PreProcess(name.Name, do)
	name.Data, _ = util.PreProcess(name.Data, do)
	name.Amount, _ = util.PreProcess(name.Amount, do)
	name.Fee, _ = util.PreProcess(name.Fee, do)

	// Set Defaults
	name.Source = useDefault(name.Source, do.Package.Account)
	name.Fee = useDefault(name.Fee, "1234")       // TODO: less hackify this.
	name.Amount = useDefault(name.Amount, "9999") // TODO: less hackify this.

	// Formulate tx
	logger.Infof("NameReg Transaction =>\t\t%s:%s:%s\n", name.Name, name.Data, name.Amount)
	tx, err := core.Name(do.Chain, do.Signer, do.PublicKey, name.Source, name.Amount, name.Nonce, name.Fee, name.Name, name.Data)
	if err != nil {
		logger.Errorf("ERROR =>\n")
		return "", err
	}

	// Sign, broadcast, display
	return txFinalize(do, tx, name.Wait)
}

func PermissionJob(perm *definitions.Permission, do *definitions.Do) (string, error) {
	// Process Variables
	perm.Source, _ = util.PreProcess(perm.Source, do)
	perm.Action, _ = util.PreProcess(perm.Action, do)
	perm.PermissionFlag, _ = util.PreProcess(perm.PermissionFlag, do)
	perm.Value, _ = util.PreProcess(perm.Value, do)
	perm.Target, _ = util.PreProcess(perm.Target, do)
	perm.Role, _ = util.PreProcess(perm.Role, do)

	// Set defaults
	perm.Source = useDefault(perm.Source, do.Package.Account)

	// Populate the transaction appropriately
	var args []string
	switch perm.Action {
	case "set_global":
		args = []string{perm.PermissionFlag, perm.Value}
	case "set_base":
		args = []string{perm.Target, perm.PermissionFlag, perm.Value}
	case "unset_base":
		args = []string{perm.Target, perm.PermissionFlag}
	case "add_role", "rm_role":
		args = []string{perm.Target, perm.Role}
	}

	// Formulate tx
	logger.Infof("Setting Permissions =>\t\t%s:%v\n", perm.Action, args)
	tx, err := core.Permissions(do.Chain, do.Signer, do.PublicKey, perm.Source, perm.Nonce, perm.Action, args)
	if err != nil {
		logger.Errorf("ERROR =>\n")
		return "", err
	}

	// Sign, broadcast, display
	return txFinalize(do, tx, perm.Wait)
}

func BondJob(bond *definitions.Bond, do *definitions.Do) (string, error) {
	// Process Variables
	bond.Account, _ = util.PreProcess(bond.Account, do)
	bond.Amount, _ = util.PreProcess(bond.Amount, do)
	bond.PublicKey, _ = util.PreProcess(bond.PublicKey, do)

	// Use Defaults
	bond.Account = useDefault(bond.Account, do.Package.Account)
	do.PublicKey = useDefault(do.PublicKey, bond.PublicKey)

	// Formulate tx
	logger.Infof("Bond Transaction =>\t\t%s:%s\n", do.PublicKey, bond.Amount)
	tx, err := core.Bond(do.Chain, do.Signer, do.PublicKey, bond.Account, bond.Amount, bond.Nonce)
	if err != nil {
		logger.Errorf("ERROR =>\n")
		return "", err
	}

	// Sign, broadcast, display
	return txFinalize(do, tx, bond.Wait)
}

func UnbondJob(unbond *definitions.Unbond, do *definitions.Do) (string, error) {
	// Process Variables
	unbond.Account, _ = util.PreProcess(unbond.Account, do)
	unbond.Height, _ = util.PreProcess(unbond.Height, do)

	// Use defaults
	unbond.Account = useDefault(unbond.Account, do.Package.Account)

	// Formulate tx
	logger.Infof("Unbond Transaction =>\t\t%s:%s\n", unbond.Account, unbond.Height)
	tx, err := core.Unbond(unbond.Account, unbond.Height)
	if err != nil {
		logger.Errorf("ERROR =>\n")
		return "", err
	}

	// Sign, broadcast, display
	return txFinalize(do, tx, unbond.Wait)
}

func RebondJob(rebond *definitions.Rebond, do *definitions.Do) (string, error) {
	// Process Variables
	rebond.Account, _ = util.PreProcess(rebond.Account, do)
	rebond.Height, _ = util.PreProcess(rebond.Height, do)

	// Use defaults
	rebond.Account = useDefault(rebond.Account, do.Package.Account)

	// Formulate tx
	logger.Infof("Rebond Transaction =>\t\t%s:%s\n", rebond.Account, rebond.Height)
	tx, err := core.Rebond(rebond.Account, rebond.Height)
	if err != nil {
		logger.Errorf("ERROR =>\n")
		return "", err
	}

	// Sign, broadcast, display
	return txFinalize(do, tx, rebond.Wait)
}

func txFinalize(do *definitions.Do, tx interface{}, wait bool) (string, error) {
	var result string

	res, err := core.SignAndBroadcast(do.ChainID, do.Chain, do.Signer, tx.(types.Tx), true, true, wait)
	if err != nil {
		logger.Errorf("ERROR =>\n")
		return "", err
	}

	if err := util.UnpackSignAndBroadcast(res, err); err != nil {
		logger.Errorf("ERROR =>\n")
		return "", err
	}

	result = fmt.Sprintf("%X", res.Hash)
	return result, nil
}

func useDefault(thisOne, defaultOne string) string {
	if thisOne == "" {
		return defaultOne
	}
	return thisOne
}
