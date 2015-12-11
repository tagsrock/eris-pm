package ebi

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
)

//This file is for storage and retrieval functions of abi's in abi subdirectory

//Write ABI []byte data into hash-named file
func WriteAbi(abiData []byte) (string, error) {
	//Construct file path based on data hash
	hash := sha256.Sum256(abiData)
	abiHash := hex.EncodeToString(hash[:])
	abiPath := path.Join(Raw, abiHash)

	//Write data
	err := ioutil.WriteFile(abiPath, abiData, 0644)
	if err != nil {
		return "", err
	}

	return abiHash, nil
}

func ReadAbiFile(abiPath string) ([]byte, string, error) {
	//Check it exists first
	if _, err := os.Stat(abiPath); err != nil {
		return nil, "", fmt.Errorf("Could not read ABI file %s: Does not Exist", abiPath)
	}

	abiData, err := ioutil.ReadFile(abiPath)
	if err != nil {
		log.Println("Failed to read abi file:", err)
		return nil, "", err
	}

	dataHash := sha256.Sum256(abiData)
	hashStr := hex.EncodeToString(dataHash[:])

	return abiData, hashStr, nil
}

func ReadAbi(abiHash string) ([]byte, string, error) {
	abiPath := path.Join(Raw, abiHash)

	abiData, dataHash, err := ReadAbiFile(abiPath)
	if err != nil {
		return nil, "", err
	}

	if dataHash != abiHash {
		return nil, "", fmt.Errorf("The retrieved Abi file's hash did not match requested")
	}

	return abiData, dataHash, nil
}

func ImportAbi(abiPath string) (string, error) {
	abiData, _, err := ReadAbiFile(abiPath)
	if err != nil {
		return "", err
	}

	abiHash, err := WriteAbi(abiData)
	if err != nil {
		return "", err
	}

	return abiHash, nil
}

func VerifyAbiHash(abiPath, abiHash string) error {
	_, dataHash, err := ReadAbiFile(abiPath)
	if err != nil {
		return err
	}

	if dataHash != abiHash {
		return fmt.Errorf("The abi data does not match its hash")
	}

	return nil
}

func VerifyAbiFile(abiPath string) error {
	//Get the filename for hash
	finfo, err := os.Stat(abiPath)
	if err != nil {
		return fmt.Errorf("Could not read ABI file %s: Does not Exist", abiPath)
	}

	filename := finfo.Name()

	//Check filename is a hexstring
	_, err = hex.DecodeString(filename)
	if err != nil {
		return fmt.Errorf("File name %s is not a hex string. Can't Compare", filename)
	}

	return VerifyAbiHash(abiPath, filename)
}
