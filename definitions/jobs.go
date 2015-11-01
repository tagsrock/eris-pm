package definitions

// ------------------------------------------------------------------------
// Util Jobs
// ------------------------------------------------------------------------

type Account struct {
	Address string `mapstructure:"address" json:"address" yaml:"address" toml:"address"`
}

type Set struct {
	Value string `mapstructure:"val" json:"val" yaml:"val" toml:"val"`
}

// ------------------------------------------------------------------------
// Transaction Jobs
// ------------------------------------------------------------------------

type Send struct {
	Source      string `mapstructure:"source" json:"source" yaml:"source" toml:"source"`
	Destination string `mapstructure:"destination" json:"destination" yaml:"destination" toml:"destination"`
	Amount      string `mapstructure:"amount" json:"amount" yaml:"amount" toml:"amount"`
	Nonce       string `mapstructure:"nonce" json:"nonce" yaml:"nonce" toml:"nonce"`
	Wait        bool   `mapstructure:"wait" json:"wait" yaml:"wait" toml:"wait"`
}

type RegisterName struct {
	Source   string `mapstructure:"source" json:"source" yaml:"source" toml:"source"`
	Name     string `mapstructure:"name" json:"name" yaml:"name" toml:"name"`
	Data     string `mapstructure:"data" json:"data" yaml:"data" toml:"data"`
	DataFile string `mapstructure:"data_file" json:"data_file" yaml:"data_file" toml:"data_file"`
	Amount   string `mapstructure:"amount" json:"amount" yaml:"amount" toml:"amount"`
	Fee      string `mapstructure:"fee" json:"fee" yaml:"fee" toml:"fee"`
	Nonce    string `mapstructure:"nonce" json:"nonce" yaml:"nonce" toml:"nonce"`
	Wait     bool   `mapstructure:"wait" json:"wait" yaml:"wait" toml:"wait"`
}

// Action: "set_base", "unset_base", "set_global", "add_role" "rm_role"
type Permission struct {
	Source         string `mapstructure:"source" json:"source" yaml:"source" toml:"source"`
	Action         string `mapstructure:"action" json:"action" yaml:"action" toml:"action"`
	PermissionFlag string `mapstructure:"permission" json:"permission" yaml:"permission" toml:"permission"`
	Value          string `mapstructure:"value" json:"value" yaml:"value" toml:"value"`
	Target         string `mapstructure:"target" json:"target" yaml:"target" toml:"target"`
	Role           string `mapstructure:"role" json:"role" yaml:"role" toml:"role"`
	Nonce          string `mapstructure:"nonce" json:"nonce" yaml:"nonce" toml:"nonce"`
	Wait           bool   `mapstructure:"wait" json:"wait" yaml:"wait" toml:"wait"`
}

type Bond struct {
	PublicKey string `mapstructure:"pub_key" json:"pub_key" yaml:"pub_key" toml:"pub_key"`
	Account   string `mapstructure:"account" json:"account" yaml:"account" toml:"account"`
	Amount    string `mapstructure:"amount" json:"amount" yaml:"amount" toml:"amount"`
	Nonce     string `mapstructure:"nonce" json:"nonce" yaml:"nonce" toml:"nonce"`
	Wait      bool   `mapstructure:"wait" json:"wait" yaml:"wait" toml:"wait"`
}

type Unbond struct {
	Account string `mapstructure:"account" json:"account" yaml:"account" toml:"account"`
	Height  string `mapstructure:"height" json:"height" yaml:"height" toml:"height"`
	Wait    bool   `mapstructure:"wait" json:"wait" yaml:"wait" toml:"wait"`
}

type Rebond struct {
	Account string `mapstructure:"account" json:"account" yaml:"account" toml:"account"`
	Height  string `mapstructure:"height" json:"height" yaml:"height" toml:"height"`
	Wait    bool   `mapstructure:"wait" json:"wait" yaml:"wait" toml:"wait"`
}

// ------------------------------------------------------------------------
// Contracts Jobs
// ------------------------------------------------------------------------

type Deploy struct {
	// TODO
}

type PackageDeploy struct {
	// TODO
}

type Call struct {
	Source      string `mapstructure:"source" json:"source" yaml:"source" toml:"source"`
	Destination string `mapstructure:"destination" json:"destination" yaml:"destination" toml:"destination"`
	Data        string `mapstructure:"data" json:"data" yaml:"data" toml:"data"`
	Amount      string `mapstructure:"amount" json:"amount" yaml:"amount" toml:"amount"`
	Nonce       string `mapstructure:"nonce" json:"nonce" yaml:"nonce" toml:"nonce"`
	Fee         string `mapstructure:"fee" json:"fee" yaml:"fee" toml:"fee"`
	Gas         string `mapstructure:"gas" json:"gas" yaml:"gas" toml:"gas"`
	Wait        bool   `mapstructure:"wait" json:"wait" yaml:"wait" toml:"wait"`
}

// ------------------------------------------------------------------------
// State Jobs
// ------------------------------------------------------------------------

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

// ------------------------------------------------------------------------
// Testing Jobs
// ------------------------------------------------------------------------

// aka. Simulated Call. Only exposed for testing
type QueryContract struct {
	Source      string   `mapstructure:"source" json:"source" yaml:"source" toml:"source"`
	Destination string   `mapstructure:"destination" json:"destination" yaml:"destination" toml:"destination"`
	Data        []string `mapstructure:"data" json:"data" yaml:"data" toml:"data"`
}

// Only exposed for testing
type QueryAccount struct {
	Account string `mapstructure:"account" json:"account" yaml:"account" toml:"account"`
	Field   string `mapstructure:"field" json:"field" yaml:"field" toml:"field"`
}

// Only exposed for testing
type QueryName struct {
	Name  string `mapstructure:"name" json:"name" yaml:"name" toml:"name"`
	Field string `mapstructure:"field" json:"field" yaml:"field" toml:"field"`
}

// Only exposed for testing
type QueryVals struct {
	Field string `mapstructure:"field" json:"field" yaml:"field" toml:"field"`
}

// Only exposed for testing
type Assert struct {
	Key      string `mapstructure:"key" json:"key" yaml:"key" toml:"key"`
	Relation string `mapstructure:"relation" json:"relation" yaml:"relation" toml:"relation"`
	Value    string `mapstructure:"val" json:"val" yaml:"val" toml:"val"`
}
