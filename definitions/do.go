package definitions

type Do struct {
	Debug         bool     `mapstructure:"," json:"," yaml:"," toml:","`
	Verbose       bool     `mapstructure:"," json:"," yaml:"," toml:","`
	YAMLPath      string   `mapstructure:"," json:"," yaml:"," toml:","`
	ContractsPath string   `mapstructure:"," json:"," yaml:"," toml:","`
	ChainHost     string   `mapstructure:"," json:"," yaml:"," toml:","`
	ChainPort     string   `mapstructure:"," json:"," yaml:"," toml:","`
	SignHost      string   `mapstructure:"," json:"," yaml:"," toml:","`
	SignPort      string   `mapstructure:"," json:"," yaml:"," toml:","`
	CompilerHost  string   `mapstructure:"," json:"," yaml:"," toml:","`
	CompilerPort  string   `mapstructure:"," json:"," yaml:"," toml:","`
	DefaultGas    uint     `mapstructure:"," json:"," yaml:"," toml:","`

	Result string
}

func NowDo() *Do {
	return &Do{}
}
