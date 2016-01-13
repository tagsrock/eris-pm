package packages

import (
	"fmt"
	"path/filepath"

	"github.com/eris-ltd/eris-pm/definitions"

	log "github.com/eris-ltd/eris-pm/Godeps/_workspace/src/github.com/Sirupsen/logrus"
	"github.com/eris-ltd/eris-pm/Godeps/_workspace/src/github.com/spf13/viper"
)

func LoadPackage(fileName string) (*definitions.Package, error) {
	log.Info("Loading EPM Package Definition.")
	var pkg = definitions.BlankPackage()
	var epmJobs = viper.New()

	// setup file
	abs, err := filepath.Abs(fileName)
	if err != nil {
		return nil, fmt.Errorf("Sorry, the marmots were unable to find the absolute path to the epm jobs file.")
	}

	path := filepath.Dir(abs)
	file := filepath.Base(abs)
	extName := filepath.Ext(file)
	bName := file[:len(file)-len(extName)]
	log.WithField("=>", path).Debug("Config Path")
	log.WithField("=>", bName).Debug("Config FileBase")

	epmJobs.AddConfigPath(path)
	epmJobs.SetConfigName(bName)

	// load file
	if err := epmJobs.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("Sorry, the marmots were unable to load the epm jobs file. Please check your path.\nERROR =>\t\t\t%v", err)
	}

	// marshall file
	if err := epmJobs.Marshal(pkg); err != nil {
		return nil, fmt.Errorf("Sorry, the marmots could not figure that epm jobs file out.\nPlease check your epm.yaml is properly formatted.\n")
	}

	return pkg, nil
}
