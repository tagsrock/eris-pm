package commands

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"

	//epm-binary-generator:IMPORT
	//	mod "github.com/eris-ltd/eris-pm/commands/modules/tendermint"

	//	color "github.com/daviddengcn/go-colortext"
	"github.com/eris-ltd/eris-pm/Godeps/_workspace/src/github.com/eris-ltd/common/go/common"
	"github.com/eris-ltd/eris-pm/epm" // ed25519 key generation
)

func Set(c *Context) error {
	args := c.Args()
	kv := make(map[string]string)
	for _, a := range args {
		spl := strings.Split(a, ":")
		if len(spl) < 2 {
			common.Exit(fmt.Errorf("Arguments for set should be given as <key>:<value>. Got %s", a))
		}
		kv[spl[0]] = spl[1]
	}

	_, _, rootDir, err := ResolveRootFlag(c)
	if err != nil {
		return err
	}

	// setup EPM object with ChainInterface
	e := epm.NewEPM(nil)
	if err := e.ReadVars(path.Join(rootDir, EPMVars)); err != nil {
		return err
	}
	for k, v := range kv {
		e.StoreVar(k, v)
	}
	// write epm variables to file
	return e.WriteVars(path.Join(rootDir, EPMVars))
}

func Plop(c *Context) error {
	var root = epm.EpmDir // TODO: much better!!!

	var toPlop string
	if len(c.Args()) > 0 {
		toPlop = c.Args()[0]
	} else {
		toPlop = ""
	}
	switch toPlop {
	case "vars":
		b, err := ioutil.ReadFile(path.Join(root, EPMVars))
		if err != nil {
			return err
		}
		if len(c.Args()) > 1 {
			// plop only a single var
			spl := strings.Split(string(b), "\n")
			for _, s := range spl {
				ss := strings.Split(s, ":")
				if ss[0] == c.Args()[1] {
					fmt.Println(ss[1])
				}
			}
		} else {
			fmt.Println(string(b))
		}
	case "abi":
		if len(c.Args()) == 1 {
			return fmt.Errorf("Specify a contract to see its abi")
		}
		e := epm.NewEPM(nil)
		e.ReadVars(path.Join(root, EPMVars))
		addr := c.Args()[1]
		var err error
		if epm.IsVar(addr) {
			addr, err = e.VarSub(addr)
			if err != nil {
				return err
			}
		}
		b, err := ioutil.ReadFile(path.Join(root, "abi", common.StripHex(addr)))
		if err != nil {
			return err
		}
		fmt.Println(string(b))
	default:
		logger.Errorln("Plop options: abi, vars")
	}
	return nil
}

// deploy a pdx file on a chain
func Deploy(c *Context) error {
	packagePath := "."
	if len(c.Args()) > 0 {
		packagePath = c.Args()[0]
	}

	contractPath := c.String("c")
	dontClear := c.Bool("dont-clear")
	diffStorage := c.Bool("diff")

	if !c.IsSet("c") {
		contractPath = DefaultContractPath
	}
	var err error
	epm.ContractPath, err = filepath.Abs(contractPath)
	if err != nil {
		return err
	}

	logger.Debugln("Contract root:", epm.ContractPath)

	// clear the cache
	if !dontClear {
		/*
			err := os.RemoveAll(epm.EpmDir) //XXX: this is scratch!
			if err != nil {
				logger.Errorln("Error clearing cache: ", err)
			}
			common.InitDataDir(epm.EpmDir)
		*/
	}

	// Startup the chain
	chain, err := LoadChain(c)
	if err != nil {
		return err
	}

	// setup EPM object with ChainInterface
	e := epm.NewEPM(chain)
	e.ReadVars(path.Join(chain.RootDir(), EPMVars))

	// comb directory for package-definition file
	// exits on error
	dir, pkg, test_ := getPkgDefFile(packagePath)

	// epm parse the package definition file
	if err := e.Parse(path.Join(dir, pkg+"."+PkgExt)); err != nil {
		return err
	}

	if diffStorage {
		e.Diff = true
	}

	// epm execute jobs
	e.ExecuteJobs()
	// write epm variables to file
	e.WriteVars(path.Join(chain.RootDir(), EPMVars))
	// wait for a block
	e.Commit()
	// run tests
	if test_ {
		results, err := e.Test(path.Join(dir, pkg+"."+TestExt))
		if err != nil {
			logger.Errorln(err)
			if results != nil {
				logger.Errorln("Failed tests:", results.FailedTests)
			}
			fmt.Printf("Testing %s.pdt failed\n", pkg)
			os.Exit(1)
		}
	}
	return nil
}

/*
// run a single epm on-chain command (endow, deploy, etc.)
func Command(c *Context) {
	root, chainType, _, err := ResolveRootFlag(c)
	common.IfExit(err)

	chain := mod.NewChain(chainType, c.Bool("rpc"))

	args := c.Args()
	if len(args) < 3 {
		exit(fmt.Errorf("You must specify a command and at least 2 arguments"))
	}
	cmd := args[0]
	args = args[1:]

	// put the args into a string and parse them
	argString := ""
	for _, a := range args {
		argString += a + " "
	}
	job := epm.ParseArgs(cmd, argString)

	// set contract path
	contractPath := c.String("c")
	if !c.IsSet("c") {
		contractPath = DefaultContractPath
	}
	epm.ContractPath, err = filepath.Abs(contractPath)
	common.IfExit(err)
	logger.Infoln("Contract path:", epm.ContractPath)

	epm.ErrMode = epm.ReturnOnErr
	// load epm
	e, err := epm.NewEPM(chain, epm.LogFile)
	common.IfExit(err)
	e.ReadVars(path.Join(root, EPMVars))

	// we don't need to turn anything on for "set"
	if cmd != "set" {
		setupModule(c, chain, root)
	}

	// run job
	e.AddJob(job)
	err = e.ExecuteJobs()
	common.IfExit(err)
	e.WriteVars(path.Join(root, EPMVars))
	// not everything needs a new block
	if cmd != "call" && cmd != "assert" && cmd != "set" {
		e.Commit()
	}
}
*/

/*
func Test(c *Context) {
	packagePath := "."
	if len(c.Args()) > 0 {
		packagePath = c.Args()[0]
	}

	contractPath := c.String("contracts")
	dontClear := c.Bool("dont-clear")
	diffStorage := c.Bool("diff")

	chainRoot, chainType, _, err := ResolveRootFlag(c)
	common.IfExit(err)
	// hierarchy : name > chainId > db > config > HEAD > default

	if !c.IsSet("contracts") {
		contractPath = DefaultContractPath
	}
	epm.ContractPath, err = filepath.Abs(contractPath)
	common.IfExit(err)

	logger.Debugln("Contract root:", epm.ContractPath)

	// clear the cache
	if !dontClear {
		err := os.RemoveAll(utils.Epm)
		if err != nil {
			logger.Errorln("Error clearing cache: ", err)
		}
		utils.InitDataDir(utils.Epm)
	}

	// read all pdxs in the dir
	fs, err := ioutil.ReadDir(packagePath)
	common.IfExit(err)
	failed := make(map[string][]int)
	for _, f := range fs {
		fname := f.Name()
		if path.Ext(fname) != ".pdx" {
			continue
		}
		sp := strings.Split(fname, ".")
		pkg := sp[0]
		dir := packagePath
		if _, err := os.Stat(path.Join(dir, pkg+".pdt")); err != nil {
			continue
		}

		// setup EPM object with ChainInterface
		var chain epm.Blockchain
		chain = LoadChain(c, chainType, chainRoot)
		e, err := epm.NewEPM(chain, epm.LogFile)
		common.IfExit(err)
		e.ReadVars(path.Join(chainRoot, EPMVars))

		// epm parse the package definition file
		err = e.Parse(path.Join(dir, fname))
		common.IfExit(err)

		if diffStorage {
			e.Diff = true
		}

		// epm execute jobs
		e.ExecuteJobs()
		// write epm variables to file
		e.WriteVars(path.Join(chainRoot, EPMVars))
		// wait for a block
		e.Commit()
		// run tests
		results, err := e.Test(path.Join(dir, pkg+"."+TestExt))
		if err != nil {
			logger.Errorln(err)
			if results != nil {
				logger.Errorln("Failed tests:", results.FailedTests)
			}
		}
		chain.Shutdown()
		if results.Err != "" {
			log.Fatal(results.Err)
		}
		if results.Failed > 0 {
			failed[pkg] = results.FailedTests
		}
	}
	if len(failed) == 0 {
		fmt.Println("All tests passed")
	} else {
		fmt.Println("Failed:")
		for p, ns := range failed {
			fmt.Println(p, ns)
		}
	}
}
*/

/*
func Console(c *Context) {

	contractPath := c.String("c")
	dontClear := c.Bool("dont-clear")
	diffStorage := c.Bool("diff")

	chainRoot, chainType, _, err := ResolveRootFlag(c)
	common.IfExit(err)
	// hierarchy : name > chainId > db > config > HEAD > default

	// Startup the chain
	var chain epm.Blockchain
	chain = LoadChain(c, chainType, chainRoot)

	if !c.IsSet("c") {
		contractPath = DefaultContractPath
	}
	epm.ContractPath, err = filepath.Abs(contractPath)
	common.IfExit(err)

	logger.Debugln("Contract root:", epm.ContractPath)

	// clear the cache
	if !dontClear {
		err := os.RemoveAll(utils.Epm)
		if err != nil {
			logger.Errorln("Error clearing cache: ", err)
		}
		utils.InitDataDir(utils.Epm)
	}

	// setup EPM object with ChainInterface
	e, err := epm.NewEPM(chain, epm.LogFile)
	common.IfExit(err)
	e.ReadVars(path.Join(chainRoot, EPMVars))

	if diffStorage {
		e.Diff = true
	}
	//e.Repl()
}
*/
