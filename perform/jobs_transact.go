package perform

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"regexp"

	"github.com/eris-ltd/eris-pm/definitions"
	"github.com/eris-ltd/eris-pm/util"

	"github.com/eris-ltd/eris-pm/Godeps/_workspace/src/github.com/eris-ltd/mint-client/mintx/core"
)

func SendJob(send *definitions.Send, do *definitions.Do) (string, error) {
	var result string

	send.Source = txPreProcess(send.Source, do.Package)
	send.Destination = txPreProcess(send.Destination, do.Package)
	send.Amount = txPreProcess(send.Amount, do.Package)

	send.Source = useDefault(send.Source, do.Package.Account)

	logger.Infof("Sending Transaction =>\t\t%s:%s:%s\n", send.Source, send.Destination, send.Amount)

	tx, err := core.Send(do.Chain, do.Signer, do.PublicKey, send.Source, send.Destination, send.Amount, send.Nonce)
	if err != nil {
		logger.Errorf("ERROR =>\t\t\t%v\n", err)
		return "", err
	}

	res, err := core.SignAndBroadcast(do.ChainID, do.Chain, do.Signer, tx, true, true, send.Wait)
	if err != nil {
		logger.Errorf("ERROR =>\t\t\t%v\n", err)
		return "", err
	}

	if err := util.UnpackSignAndBroadcast(res, err); err != nil {
		logger.Errorf("ERROR =>\t\t\t%v\n", err)
		return "", err
	}

	result = fmt.Sprintf("%X", res.Hash)
	return result, nil
}

func RegisterNameJob(name *definitions.RegisterName, do *definitions.Do) (string, error) {
	var result string

	// If a data file is given it should be in csv format and
	// it will be read first. Once the file is parsed and sent
	// to the chain then a single nameRegTx will be sent if that
	// has been populated.
	if name.DataFile != "" {
		// open the file and use a reader
		fileReader, err := os.Open(name.DataFile)
		if err != nil {
			logger.Errorf("ERROR =>\t\t\t%v\n", err)
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
				logger.Errorf("ERROR =>\t\t\t%v\n", err)
				return "", err
			}

			// Sink the Fee into the third slot in the record if
			//   it doesn't exist
			if record[2] == "" {
				record[2] = name.Amount
			}

			// Send an individual Tx for the record
			// [TODO]: move these to async using goroutines?
			registerNameTx(&definitions.RegisterName{
				Source: name.Source,
				Fee:    name.Fee,
				Wait:   name.Wait,
				Name:   record[0],
				Data:   record[1],
				Amount: record[2],
			}, do)
		}
	}

	// If the data field is populated then there is a single
	// nameRegTx to send. So do that *now*.
	if name.Data != "" {
		return registerNameTx(name, do)
	} else {
		return result, nil
	}
}

// Runs an individual nametx.
func registerNameTx(name *definitions.RegisterName, do *definitions.Do) (string, error) {
	var result string

	name.Source = txPreProcess(name.Source, do.Package)
	name.Name = txPreProcess(name.Name, do.Package)
	name.Data = txPreProcess(name.Data, do.Package)
	name.Amount = txPreProcess(name.Amount, do.Package)
	name.Fee = txPreProcess(name.Fee, do.Package)

	name.Source = useDefault(name.Source, do.Package.Account)

	logger.Infof("NameReg Transaction =>\t\t%s:%s\n", name.Name, name.Data)

	tx, err := core.Name(do.Chain, do.Signer, do.PublicKey, name.Source, name.Amount, name.Nonce, name.Fee, name.Name, name.Data)
	if err != nil {
		logger.Errorf("ERROR =>\t\t\t%v\n", err)
		return "", err
	}

	res, err := core.SignAndBroadcast(do.ChainID, do.Chain, do.Signer, tx, true, true, name.Wait)
	if err != nil {
		logger.Errorf("ERROR =>\t\t\t%v\n", err)
		return "", err
	}

	if err := util.UnpackSignAndBroadcast(res, err); err != nil {
		logger.Errorf("ERROR =>\t\t\t%v\n", err)
		return "", err
	}

	result = fmt.Sprintf("%X", res.Hash)
	return result, nil
}

func CallJob(call *definitions.Call, do *definitions.Do) (string, error) {
	var result string

	logger.Infof("Calling =>\t\t\t%s:%v", call.Destination, call.Data)
	tx, err := core.Call(do.Chain, do.Signer, do.PublicKey, call.Source, call.Destination, call.Amount, call.Nonce, call.Gas, call.Gas, call.Data)
	if err != nil {
		return "", err
	}

	if err := util.UnpackSignAndBroadcast(core.SignAndBroadcast(do.ChainID, do.Chain, do.Signer, tx, true, true, call.Wait)); err != nil {
		return "", err
	}
	return result, nil
}

func PermissionJob(perm *definitions.Permission, do *definitions.Do) (string, error) {
	var result string

	perm.Source = txPreProcess(perm.Source, do.Package)
	perm.Action = txPreProcess(perm.Action, do.Package)
	perm.PermissionFlag = txPreProcess(perm.PermissionFlag, do.Package)
	perm.Value = txPreProcess(perm.Value, do.Package)
	perm.Target = txPreProcess(perm.Target, do.Package)
	perm.Role = txPreProcess(perm.Role, do.Package)

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

	logger.Infof("Setting Permissions =>\t\t%s:%v\n", perm.Action, args)
	tx, err := core.Permissions(do.Chain, do.Signer, do.PublicKey, perm.Source, perm.Nonce, perm.Action, args)
	if err != nil {
		logger.Errorf("ERROR =>\t\t\t%v\n", err)
		return "", err
	}

	res, err := core.SignAndBroadcast(do.ChainID, do.Chain, do.Signer, tx, true, true, perm.Wait)
	if err != nil {
		logger.Errorf("ERROR =>\t\t\t%v\n", err)
		return "", err
	}

	if err := util.UnpackSignAndBroadcast(res, err); err != nil {
		logger.Errorf("ERROR =>\t\t\t%v\n", err)
		return "", err
	}

	result = fmt.Sprintf("%X", res.Hash)
	return result, nil
}

func BondJob(bond *definitions.Bond, do *definitions.Do) (string, error) {
	var result string

	logger.Infof("Bond Transaction =>\t\t%s:%s", bond.UnbondAccount, bond.Amount)
	tx, err := core.Bond(do.Chain, do.Signer, do.PublicKey, bond.UnbondAccount, bond.Amount, bond.Nonce)
	if err != nil {
		return "", err
	}

	logger.Debugf("Bond Transaction Sent =>\t\t%v\n", tx)
	if err := util.UnpackSignAndBroadcast(core.SignAndBroadcast(do.ChainID, do.Chain, do.Signer, tx, true, true, bond.Wait)); err != nil {
		return "", err
	}
	return result, nil
}

func UnbondJob(unbond *definitions.Unbond, do *definitions.Do) (string, error) {
	var result string

	logger.Infof("Unbond Transaction =>\t\t%s:%s", unbond.Account, unbond.Height)
	tx, err := core.Unbond(unbond.Account, unbond.Height)
	if err != nil {
		return "", err
	}

	logger.Debugf("Unbond Transaction Sent =>\t\t%v\n", tx)
	if err := util.UnpackSignAndBroadcast(core.SignAndBroadcast(do.ChainID, do.Chain, do.Signer, tx, true, true, unbond.Wait)); err != nil {
		return "", err
	}
	return result, nil
}

func RebondJob(rebond *definitions.Rebond, do *definitions.Do) (string, error) {
	var result string

	logger.Infof("Rebond Transaction =>\t\t%s:%s", rebond.Account, rebond.Height)
	tx, err := core.Rebond(rebond.Account, rebond.Height)
	if err != nil {
		return "", err
	}

	logger.Debugf("Rebond Transaction Sent =>\t\t%v\n", tx)
	if err := util.UnpackSignAndBroadcast(core.SignAndBroadcast(do.ChainID, do.Chain, do.Signer, tx, true, true, rebond.Wait)); err != nil {
		return "", err
	}
	return result, nil
}

func txPreProcess(toProcess string, pkg *definitions.Package) string {
	catchEr := regexp.MustCompile("^\\$(.*)$")
	if catchEr.MatchString(toProcess) {
		jobName := catchEr.FindStringSubmatch(toProcess)[1]
		for _, job := range pkg.Jobs {
			if string(jobName) == job.JobName {
				logger.Debugf("Fixing Variables =>\t\t$%s:%s\n", string(jobName), job.JobResult)
				return job.JobResult
			}
		}
	}
	return toProcess
}

func useDefault(thisOne, defaultOne string) string {
	if thisOne == "" {
		return defaultOne
	}
	return thisOne
}
