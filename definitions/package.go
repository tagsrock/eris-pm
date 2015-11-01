package definitions

type Package struct {
	Account string
	Jobs    []*Jobs
}

type Jobs struct {
	JobName   string `mapstructure:"name" json:"name" yaml:"name" toml:"name"`
	Job       *Job   `mapstructure:"job" json:"job" yaml:"job" toml:"job"`
	JobResult string
}

type Job struct {
	// Sets/Resets the primary account to use
	Account *Account `mapstructure:"account" json:"account" yaml:"account" toml:"account"`

	// Set an arbitrary value
	Set *Set `mapstructure:"set" json:"set" yaml:"set" toml:"set"`

	// Contract compile and send to the chain functions
	// @dennismckinnon working on these
	Deploy        *Deploy        `mapstructure:"deploy" json:"deploy" yaml:"deploy" toml:"deploy"`
	PackageDeploy *PackageDeploy `mapstructure:"package-deploy" json:"package-deploy" yaml:"package-deploy" toml:"package-deploy"`

	// Wrapper for mintx
	Send         *Send         `mapstructure:"send" json:"send" yaml:"send" toml:"send"`
	Bond         *Bond         `mapstructure:"bond" json:"bond" yaml:"bond" toml:"bond"`
	Unbond       *Unbond       `mapstructure:"unbond" json:"unbond" yaml:"unbond" toml:"unbond"`
	Rebond       *Rebond       `mapstructure:"rebond" json:"rebond" yaml:"rebond" toml:"rebond"`
	Permission   *Permission   `mapstructure:"permission" json:"permission" yaml:"permission" toml:"permission"`
	RegisterName *RegisterName `mapstructure:"register" json:"register" yaml:"register" toml:"register"`
	Call         *Call         `mapstructure:"call" json:"call" yaml:"call" toml:"call"`

	// Wrapper for mintdump
	DumpState    *DumpState    `mapstructure:"dump-state" json:"dump-state" yaml:"dump-state" toml:"dump-state"`
	RestoreState *RestoreState `mapstructure:"restore-state" json:"restore-state" yaml:"restore-state" toml:"restore-state"`

	// Used for Tests Only
	QueryContract *QueryContract `mapstructure:"query-contract" json:"query-contract" yaml:"query-contract" toml:"query-contract"`
	QueryAccount  *QueryAccount  `mapstructure:"query-account" json:"query-account" yaml:"query-account" toml:"query-account"`
	QueryName     *QueryName     `mapstructure:"query-name" json:"query-name" yaml:"query-name" toml:"query-name"`
	QueryVals     *QueryVals     `mapstructure:"query-vals" json:"query-vals" yaml:"query-vals" toml:"query-vals"`
	Assert        *Assert        `mapstructure:"assert" json:"assert" yaml:"assert" toml:"assert"`
}

func BlankPackage() *Package {
	return &Package{}
}
