package util

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/eris-ltd/eris-pm/definitions"
)

const logFileName = "epm.log"

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
