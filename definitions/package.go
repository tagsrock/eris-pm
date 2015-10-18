package definitions

type Package struct {
	// Sets/Resets the primary account to use
	Account        *Account       `mapstructure:"account" json:"account" yaml:"account" toml:"account"`

	// Set an arbitrary value
	Set            *Set           `mapstructure:"set" json:"set" yaml:"set" toml:"set"`

	// Contract compile and send to the chain functions
	Deploy         *Deploy        `mapstructure:"deploy" json:"deploy" yaml:"deploy" toml:"deploy"`
	Include        *Include       `mapstructure:"include" json:"include" yaml:"include" toml:"include"`
	ModifyDeploy   *ModifyDeploy  `mapstructure:"modify-deploy" json:"modify-deploy" yaml:"modify-deploy" toml:"modify-deploy"`

	// @dennismckinnon working on this
	// PackageDeploy  *PackageDeploy `mapstructure:"package-deploy" json:"package-deploy" yaml:"package-deploy" toml:"package-deploy"`

	// Wrapper for mintx
	Send					 *Send          `mapstructure:"send" json:"send" yaml:"send" toml:"send"`
	Bond           *Bond          `mapstructure:"bond" json:"bond" yaml:"bond" toml:"bond"`
	Unbond         *Unbond        `mapstructure:"unbond" json:"unbond" yaml:"unbond" toml:"unbond"`
	Rebond         *Rebond        `mapstructure:"rebond" json:"rebond" yaml:"rebond" toml:"rebond"`
	Permission     *Permission    `mapstructure:"permission" json:"permission" yaml:"permission" toml:"permission"`
	RegisterName   *RegisterName  `mapstructure:"register" json:"register" yaml:"register" toml:"register"`
	Call           *Call					`mapstructure:"call" json:"call" yaml:"call" toml:"call"`

	// Wrapper for mintdump
	DumpState      *DumpState     `mapstructure:"dump-state" json:"dump-state" yaml:"dump-state" toml:"dump-state"`
	RestoreState   *RestoreState  `mapstructure:"restore-state" json:"restore-state" yaml:"restore-state" toml:"restore-state"`

	// Used for Tests Only
	Query          *Query         `mapstructure:"query" json:"query" yaml:"query" toml:"query"`
	GetNameEntry   *GetNameEntry  `mapstructure:"query-name" json:"query-name" yaml:"query-name" toml:"query-name"`
	Assert         *Assert        `mapstructure:"assert" json:"assert" yaml:"assert" toml:"assert"`
}
