package util

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"

	"github.com/eris-ltd/eris-pm/definitions"

	"github.com/eris-ltd/eris-pm/Godeps/_workspace/src/github.com/tendermint/tendermint/wire"
)

const logFileName = "epm.log"

// ------------------------------------------------------------------------
// Logging
// ------------------------------------------------------------------------

// WriteJobResult takes two strings and writes those to the delineated log
// file, which is currently epm.log in the same directory as the epm.yaml
func WriteJobResult(name, result string) error {
	// TODO: add logging path besides pwd
	logFile := setPath()

	var file *os.File
	var err error

	if _, err := os.Stat(logFile); os.IsNotExist(err) {
		file, err = os.Create(logFile)
	} else {
		file, err = os.OpenFile(logFile, os.O_APPEND|os.O_WRONLY, 0600)
	}

	if err != nil {
		return err
	}

	defer file.Close()

	text := fmt.Sprintf("%s,%s\n", name, result)
	if _, err = file.WriteString(text); err != nil {
		return err
	}

	return nil
}

func ClearJobResults() error {
	// TODO: add logging path besides pwd
	return os.Remove(setPath())
}

func PrintPathPackage(do *definitions.Do) {
	logger.Infof("Using Compiler at =>\t\t%s\n", do.Compiler)
	logger.Infof("Using Chain at =>\t\t%s\n", do.Chain)
	logger.Debugf("\twith ChainID =>\t\t%s\n", do.ChainID)
	logger.Infof("Using Signer at =>\t\t%s\n", do.Signer)
}

func setPath() string {
	pwd, _ := os.Getwd()
	return filepath.Join(pwd, logFileName)
}

// ------------------------------------------------------------------------
// Writers of Arbitrary stuff
// ------------------------------------------------------------------------

// FormatOutput formats arbitrary json in a viewable manner using reflection
func FormatOutput(args []string, i int, o interface{}) (string, error) {
	if len(args) < i+1 {
		return prettyPrint(o)
	}
	arg0 := args[i]
	v := reflect.ValueOf(o).Elem()
	name, err := fieldFromTag(v, arg0)
	if err != nil {
		return "", err
	}
	f := v.FieldByName(name)
	return prettyPrint(f.Interface())
}

func fieldFromTag(v reflect.Value, field string) (string, error) {
	iv := v.Interface()
	st := reflect.TypeOf(iv)
	for i := 0; i < v.NumField(); i++ {
		tag := st.Field(i).Tag.Get("json")
		if tag == field {
			return st.Field(i).Name, nil
		}
	}
	return "", fmt.Errorf("Invalid field name")
}

func prettyPrint(o interface{}) (string, error) {
	var prettyJSON bytes.Buffer
	err := json.Indent(&prettyJSON, wire.JSONBytes(o), "", "\t")
	if err != nil {
		return "", err
	}
	return string(prettyJSON.Bytes()), nil
}
