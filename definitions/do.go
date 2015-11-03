package definitions

type Do struct {
	Debug         bool     `mapstructure:"," json:"," yaml:"," toml:","`
	Verbose       bool     `mapstructure:"," json:"," yaml:"," toml:","`
	YAMLPath      string   `mapstructure:"," json:"," yaml:"," toml:","`
	ContractsPath string   `mapstructure:"," json:"," yaml:"," toml:","`
	ABIPath       string   `mapstructure:"," json:"," yaml:"," toml:","`
	Chain         string   `mapstructure:"," json:"," yaml:"," toml:","`
	Signer        string   `mapstructure:"," json:"," yaml:"," toml:","`
	Compiler      string   `mapstructure:"," json:"," yaml:"," toml:","`
	PublicKey     string   `mapstructure:"," json:"," yaml:"," toml:","`
	ChainID       string   `mapstructure:"," json:"," yaml:"," toml:","`
	DefaultGas    uint     `mapstructure:"," json:"," yaml:"," toml:","`
	DefaultAddr   string   `mapstructure:"," json:"," yaml:"," toml:","`
	DefaultSets   []string `mapstructure:"," json:"," yaml:"," toml:","`

	Package *Package
	Result  string
}

func NowDo() *Do {
	return &Do{}
}
