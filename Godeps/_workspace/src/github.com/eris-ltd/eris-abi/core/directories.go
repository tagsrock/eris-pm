package ebi

import (
	"fmt"
	"os"
	"path"

	"github.com/eris-ltd/eris-pm/Godeps/_workspace/src/github.com/eris-ltd/common/go/common"
)

func GetAbiRoot() string {
	var abiroot string
	if os.Getenv("ERIS_ABI_ROOT") != "" {
		abiroot = os.Getenv("ERIS_ABI_ROOT")
	} else {
		abiroot = path.Join(common.ErisRoot, "abi")
	}
	return abiroot
}

//Directory Structure
//.eris
//   +-abi
//      +-index
//      |    + abi indexing jsons
//      |
//      +-raw
//           + abi files (hash named)

var (
	Root  = GetAbiRoot()
	Index = path.Join(Root, "index")
	Raw   = path.Join(Root, "raw")
)

func InitPaths() {
	Root = GetAbiRoot()
	Index = path.Join(Root, "index")
	Raw = path.Join(Root, "raw")
}

func BuildDirTree() error {
	//Check if abi root exists.
	if _, err := os.Stat(Root); err != nil {
		//create it
		err = os.MkdirAll(Root, 0700)
		if err != nil {
			return fmt.Errorf("Failed to create the abi root directory")
		}
	}

	//Create Indexing folder
	if _, err := os.Stat(Index); err != nil {
		//create it
		err = os.MkdirAll(Index, 0700)
		if err != nil {
			return fmt.Errorf("Failed to create the abi index directory")
		}
	}

	//Create Raw Storage Folder
	if _, err := os.Stat(Raw); err != nil {
		//create it
		err = os.MkdirAll(Raw, 0700)
		if err != nil {
			return fmt.Errorf("Failed to create the abi raw directory")
		}
	}

	return nil
}

func CheckDirTree() error {
	if _, err := os.Stat(Root); err != nil {
		return fmt.Errorf("Abi root directory does not exist.")
	}

	if _, err := os.Stat(Index); err != nil {
		return fmt.Errorf("Abi index directory does not exist.")
	}

	if _, err := os.Stat(Raw); err != nil {
		return fmt.Errorf("Abi raw directory does not exist.")
	}

	return nil
}
