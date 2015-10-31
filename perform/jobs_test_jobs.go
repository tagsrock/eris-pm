package perform

import (
	// "encoding/hex"

	"github.com/eris-ltd/eris-pm/definitions"

	// cclient "github.com/eris-ltd/eris-pm/Godeps/_workspace/src/github.com/tendermint/tendermint/rpc/core_client"
)

func QueryJob(query *definitions.Query) (string, error) {
	var result string

	// 	logger.Infof("Querying =>\t\t\t%s:%v", query.Destination, query.Data)

	// 	fromAddrBytes, err := hex.DecodeString(query.Source)
	// 	if err != nil {
	// 		return err
	// 	}
	// 	toAddrBytes, err := hex.DecodeString(query.Destination)
	// 	if err != nil {
	// 		return err
	// 	}
	// 	dataBytes, err := hex.DecodeString(query.Data)
	// 	if err != nil {
	// 		return err
	// 	}

	// 	r, err := cclient.Client.Call(fromAddrBytes, toAddrBytes, dataBytes)
	// 	if err != nil {
	// 		return err
	// 	}
	// 	s, err := util.FormatOutput(args, 3, r)
	// 	if err != nil {
	// 		return err
	// 	}
	// 	logger.Println(s)
	return result, nil
}

func GetNameEntryJob(name *definitions.GetNameEntry) (string, error) {
	var result string

	return result, nil
}

func AssertJob(assertion *definitions.Assert) (string, error) {
	var result string

	return result, nil
}
