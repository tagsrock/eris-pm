package commands

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"strings"

	//epm-binary-generator:IMPORT
	// mod "github.com/eris-ltd/eris-pm/commands/modules/tendermint"

	"github.com/eris-ltd/eris-pm/Godeps/_workspace/src/github.com/eris-ltd/common/go/common"
	"github.com/eris-ltd/eris-pm/epm"
)

func cleanupEPM() {
	dirs := []string{epm.EpmDir}
	for _, d := range dirs {
		err := os.RemoveAll(d)
		if err != nil {
			logger.Errorln("Error removing dir", d, err)
		}
	}
}

func installEPM() {
	cur, _ := os.Getwd()
	os.Chdir(path.Join(common.ErisLtd, "eris-pm", "cmd", "epm"))
	cmd := exec.Command("go", "install")
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Run()
	logger.Infoln(out.String())
	os.Chdir(cur)
}

func pullErisRepo(repo, branch string) {
	// pull changes
	os.Chdir(path.Join(common.ErisLtd, repo))
	cmd := exec.Command("git", "pull", "origin", branch)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Run()
	res := out.String()
	logger.Infoln(res)
}

func updateEPM() {
	cur, _ := os.Getwd()

	pullErisRepo("epm-go", "master")
	pullErisRepo("decerver-interfaces", "master")

	// install
	installEPM()

	// return to original dir
	os.Chdir(cur)
}

func cleanPullUpdate(clean, pull, update bool) {
	if clean && pull {
		cleanupEPM()
		updateEPM()
	} else if clean {
		cleanupEPM()
		if update {
			installEPM()
		}
	} else if pull {
		updateEPM()
	} else if update {
		installEPM()
	}
}

// looks for pkg-def file
// common.Exits if error (none or more than 1)
// returns dir of pkg, name of pkg (no extension) and whether or not there's a test file
func getPkgDefFile(pkgPath string) (string, string, bool) {
	logger.Infoln("Pkg path:", pkgPath)
	var pkgName string
	var test_ bool

	// if its not a directory, look for a corresponding test file
	f, err := os.Stat(pkgPath)
	common.IfExit(err)

	if !f.IsDir() {
		dir, fil := path.Split(pkgPath)
		spl := strings.Split(fil, ".")
		pkgName = spl[0]
		ext := spl[1]
		if ext != PkgExt {
			common.Exit(fmt.Errorf("Did not understand extension. Got %s, expected %s\n", ext, PkgExt))
		}

		_, err := os.Stat(path.Join(dir, pkgName) + "." + TestExt)
		if err != nil {
			logger.Errorf("There was no test found for package-definition %s. Deploying without test ...\n", pkgName)
			test_ = false
		} else {
			test_ = true
		}
		return dir, pkgName, test_
	}

	// read dir for files
	files, err := ioutil.ReadDir(pkgPath)
	common.IfExit(err)

	// find all package-defintion and package-definition-test files
	candidates := make(map[string]int)
	candidates_test := make(map[string]int)
	for _, f := range files {
		name := f.Name()
		spl := strings.Split(name, ".")
		if len(spl) < 2 {
			continue
		}
		name = spl[0]
		ext := spl[1]
		if ext == PkgExt {
			candidates[name] = 1
		} else if ext == TestExt {
			candidates_test[name] = 1
		}
	}
	// common.Exit if too many or no options
	if len(candidates) > 1 {
		common.Exit(fmt.Errorf("More than one package-definition file available. Please select with the '-p' flag"))
	} else if len(candidates) == 0 {
		common.Exit(fmt.Errorf("No package-definition files found for extensions %s, %s", PkgExt, TestExt))
	}
	// this should run once (there's only one candidate)
	for k, _ := range candidates {
		pkgName = k
		if candidates_test[pkgName] == 1 {
			test_ = true
		} else {
			logger.Infoln("There was no test found for package-definition %s. Deploying without test ...\n", pkgName)
			test_ = false
		}
	}
	return pkgPath, pkgName, test_
}

func checkInit() error {
	return nil
}

func confirm(message string) bool {
	fmt.Println(message, "Are you sure? (y/n)")
	var r string
	fmt.Scanln(&r)
	for ; ; fmt.Scanln(&r) {
		if r == "n" || r == "y" {
			break
		} else {
			fmt.Printf("Yes or no?", r)
		}
	}
	return r == "y"
}
