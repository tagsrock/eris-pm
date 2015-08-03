package commands

import (
	"fmt"
	"os"
	"path"

	//epm-binary-generator:IMPORT
	mod "github.com/eris-ltd/eris-pm/commands/modules/tendermint"

	"github.com/spf13/viper"

	"github.com/eris-ltd/eris-pm/epm"
)

// this needs to match the type of the chain we're trying to run
// it should be blank for the base epm (even though it includes thel ...)
//epm-binary-generator:CHAIN
const CHAIN = ""

func LoadChain(c *Context) (epm.ChainClient, error) {
	var chainExists bool
	rootDir, chainType, chainID, err := ResolveRootFlag(c)
	if err == nil {
		// if the dir hasnt been laid yet, the root wont resolve
		chainExists = true
		rootDir = ComposeRoot(chainType, chainID)
	}

	if chainExists {
		viper.SetConfigName("config") // name of config file (without extension)
		viper.AddConfigPath(rootDir)  // path to look for the config file in
		viper.ReadInConfig()          // Find and read the config file
	}

	var rpcAddr, signAddr string
	var pubkey string

	if c.IsSet("host") || c.IsSet("port") || !chainExists {
		rpcHost, rpcPort := c.String("host"), c.Int("port")
		rpcAddr = fmt.Sprintf("http://%s:%d/", rpcHost, rpcPort)
	} else {
		rpcAddr = viper.GetString("node_addr")
	}

	if c.IsSet("sign_host") || c.IsSet("sign_port") || !chainExists {
		signHost, signPort := c.String("sign_host"), c.Int("sign_port")
		signAddr = fmt.Sprintf("http://%s:%d", signHost, signPort)
		fmt.Println("GAHH", signHost, signPort, signAddr)
	} else {
		signAddr = viper.GetString("sign_addr")
	}

	if chainExists {
		chainID = viper.GetString("chain_id")
	}

	if c.IsSet("pubkey") || !chainExists {
		pubkey = c.String("pubkey")
	} else {
		pubkey = viper.GetString("pubkey")
	}

	logger.Debugln("ChainType, ChainID, RootDir", chainType, chainID, rootDir)

	chain, err := mod.NewChain(rootDir, chainID, rpcAddr, signAddr, pubkey)
	if err != nil {
		return nil, err
	}

	// TODO: set contract path?

	if _, err := os.Stat(path.Join(rootDir, "config.toml")); err != nil {
		if _, err := os.Stat(rootDir); err != nil {
			if err := os.MkdirAll(rootDir, 0777); err != nil {
				logger.Errorln("Failed to write config file", err)
			}
		}
		if err := chain.WriteConfig(); err != nil {
			logger.Errorln("Failed to write config file", err)
		}
	}

	return chain, nil
}
