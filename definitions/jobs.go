package definitions

type Account struct {
	Name    string `mapstructure:"name" json:"name" yaml:"name" toml:"name"`
	Address string `mapstructure:"address" json:"address" yaml:"address" toml:"address"`
}

type Set struct {
	Key   string `mapstructure:"key" json:"key" yaml:"key" toml:"key"`
	Value string `mapstructure:"val" json:"val" yaml:"val" toml:"val"`
}

type Deploy struct {
	// TODO
}

type Include struct {
	// TODO
}

type ModifyDeploy struct {
	// TODO
}

type PackageDeploy struct {
	// TODO
}

type Send struct {
	Source      string `mapstructure:"source" json:"source" yaml:"source" toml:"source"`
	Destination string `mapstructure:"destination" json:"destination" yaml:"destination" toml:"destination"`
	Amount      uint   `mapstructure:"amount" json:"amount" yaml:"amount" toml:"amount"`
	Wait        bool   `mapstructure:"wait" json:"wait" yaml:"wait" toml:"wait"`
}

type Bond struct {
	Account     string `mapstructure:"account" json:"account" yaml:"account" toml:"account"`
	Amount      uint   `mapstructure:"amount" json:"amount" yaml:"amount" toml:"amount"`
	Destination string `mapstructure:"destination" json:"destination" yaml:"destination" toml:"destination"`
}

type Unbond struct {
	Account string `mapstructure:"account" json:"account" yaml:"account" toml:"account"`
	Height  uint   `mapstructure:"height" json:"height" yaml:"height" toml:"height"`
}

type Rebond struct {
	Account string `mapstructure:"account" json:"account" yaml:"account" toml:"account"`
	Height  uint   `mapstructure:"height" json:"height" yaml:"height" toml:"height"`
}

type RegisterName struct {
	Source string `mapstructure:"source" json:"source" yaml:"source" toml:"source"`
	Name 	 string `mapstructure:"name" json:"name" yaml:"name" toml:"name"`
	Data   string `mapstructure:"data" json:"data" yaml:"data" toml:"data"`
	Amount uint   `mapstructure:"amount" json:"amount" yaml:"amount" toml:"amount"`
	Fee    uint   `mapstructure:"fee" json:"fee" yaml:"fee" toml:"fee"`
	Wait   bool   `mapstructure:"wait" json:"wait" yaml:"wait" toml:"wait"`
}

type Call struct {
	Source      string   `mapstructure:"source" json:"source" yaml:"source" toml:"source"`
	Destination string   `mapstructure:"destination" json:"destination" yaml:"destination" toml:"destination"`
	Data        []string `mapstructure:"data" json:"data" yaml:"data" toml:"data"`
	Amount      uint     `mapstructure:"amount" json:"amount" yaml:"amount" toml:"amount"`
	Fee         uint     `mapstructure:"fee" json:"fee" yaml:"fee" toml:"fee"`
	Gas         uint     `mapstructure:"gas" json:"gas" yaml:"gas" toml:"gas"`
	Wait        bool     `mapstructure:"wait" json:"wait" yaml:"wait" toml:"wait"`
}

type Permission struct {
	// todo
}

type DumpState struct {
	WithValidators bool   `mapstructure:"include-validators" json:"include-validators" yaml:"include-validators" toml:"include-validators"`
	ToIPFS         bool   `mapstructure:"to-ipfs" json:"to-ipfs" yaml:"to-ipfs" toml:"to-ipfs"`
	ToFile         bool   `mapstructure:"to-file" json:"to-file" yaml:"to-file" toml:"to-file"`
	IPFSHost       string `mapstructure:"ipfs-host" json:"ipfs-host" yaml:"ipfs-host" toml:"ipfs-host"`
	FilePath       string `mapstructure:"file" json:"file" yaml:"file" toml:"file"`
}

type RestoreState struct {
	FromIPFS bool   `mapstructure:"from-ipfs" json:"from-ipfs" yaml:"from-ipfs" toml:"from-ipfs"`
	FromFile bool   `mapstructure:"from-file" json:"from-file" yaml:"from-file" toml:"from-file"`
	IPFSHost string `mapstructure:"ipfs-host" json:"ipfs-host" yaml:"ipfs-host" toml:"ipfs-host"`
	FilePath string `mapstructure:"file" json:"file" yaml:"file" toml:"file"`
}

type GetNameEntry struct {
	Name string `mapstructure:"name" json:"name" yaml:"name" toml:"name"`
	Data string `mapstructure:"data" json:"data" yaml:"data" toml:"data"`
}

type Query struct {
	Account string `mapstructure:"account" json:"account" yaml:"account" toml:"account"`
	Storage string `mapstructure:"storage" json:"storage" yaml:"storage" toml:"storage"` // this probably isnt right. HELP!
}

type Assert struct {
	Key      string `mapstructure:"key" json:"key" yaml:"key" toml:"key"`
	Relation string `mapstructure:"relation" json:"relation" yaml:"relation" toml:"relation"`
	Value    string `mapstructure:"val" json:"val" yaml:"val" toml:"val"`
}