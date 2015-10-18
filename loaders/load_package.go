package loaders

import (
	"fmt"
	"path/filepath"

	"github.com/eris-ltd/eris-pm/definitions"

	"github.com/eris-ltd/eris-pm/Godeps/_workspace/src/github.com/spf13/viper"
)

func LoadPackage(fileName string) (*definitions.Package, error) {
	var pkg = definitions.BlankPackage()
	var epmJobs = viper.New()

	// setup file
	abs, err := filepath.Abs(fileName)
	if err != nil {
		return nil, fmt.Errorf("Sorry, the marmots were unable to find the absolute path to the epm jobs file.")
	}
	epmJobs.AddConfigPath(filepath.Dir(abs))
	epmJobs.SetConfigName(filepath.Base(abs))

	// load file
	if err := epmJobs.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("Sorry, the marmots were unable to load the epm jobs file. Please check your path.")
	}

	// marshall file
	if err := epmJobs.Marshal(pkg); err != nil {
		return nil, fmt.Errorf("Sorry, the marmots could not figure that epm jobs file out.\nPlease check your epm.yaml is properly formatted.\n")
	}

	return pkg, nil
}
