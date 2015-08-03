package commands

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/eris-ltd/common/go/common"
)

//-------------------------------------------------------
// resolve chains by HEAD, ref, or type/id

// Returns root, type, id, err
func ResolveRootFlag(c *Context) (string, string, string, error) {
	ref := c.String("chain")
	if ref == "" {
		ref = fmt.Sprintf("%s/%s", c.String("type"), c.String("chain_id"))
	}
	return resolveRoot(ref)
}

// Returns root, type, id, err
func ResolveRootArg(c *Context) (string, string, string, error) {
	args := c.Args()
	ref := ""
	if len(args) > 0 {
		ref = args[0]
	}
	return resolveRoot(ref)
}

// Returns root, type, id, err
func resolveRoot(ref string) (string, string, string, error) {
	chainType, chainId, err := ResolveChain(ref)
	root := ComposeRoot(chainType, chainId)
	return root, chainType, chainId, err
}

// Compose path to this chain's root, including multichain ref
func ComposeRoot(chainType, chainId string) string {
	return path.Join(common.BlockchainsPath, chainType, chainId)
}

// Get ChainId from a reference name by reading the ref file
func ChainFromName(name string) (string, string, error) {
	refsPath := path.Join(common.BlockchainsPath, "refs", name)
	b, err := ioutil.ReadFile(refsPath)
	if err != nil {
		return "", "", err
	}
	s := string(b)
	return SplitRef(s)
}

// Allow chain types to be specified in shorter form (ie. 'eth' for 'ethereum')
func ResolveChainType(chainType string) (string, error) {
	switch chainType {
	/*case "thel", "thelonious", "monk":
		return "thelonious", nil
	case "btc", "bitcoin":
		return "bitcoin", nil
	case "eth", "ethereum":
		return "ethereum", nil
	*/
	case "mint", "tendermint":
		return "tendermint", nil
	}
	return "", fmt.Errorf("Unknown chain type: %s", chainType)
}

// Determines the chainId from a chainId prefix and a type
func ResolveChainId(chainType, chainId string) (string, error) {
	if chainId == "" {
		return "", fmt.Errorf("Empty chainId")
	}

	var err error

	chainType, err = ResolveChainType(chainType)
	if err != nil {
		return chainId, err
	}

	p := ComposeRoot(chainType, chainId)
	if _, err := os.Stat(p); err != nil {
		// see if its a prefix of a chainId
		id, err := findPrefixMatch(path.Join(common.BlockchainsPath, chainType), chainId)
		if err != nil {
			return chainId, err
		}
		p = path.Join(common.BlockchainsPath, chainType, id)
		chainId = id
	}
	if _, err := os.Stat(p); err != nil {
		return chainId, fmt.Errorf("Could not locate %s chain by id %s", chainType, chainId)
	}

	return chainId, nil
}

// Resolve a chain's type and id from a a reference
// Reference is either blank (head), a ref name, or <type>/<id>
func ResolveChain(ref string) (chainType string, chainId string, err error) {
	if ref == "" {
		return GetHead()
	}

	chainType, chainId, err = SplitRef(ref)
	if err == nil {
		if chainType, err = ResolveChainType(chainType); err != nil {
			return
		}
		chainId, err = ResolveChainId(chainType, chainId)
		return
	}

	return ChainFromName(ref)
}

// Return full path to a blockchain's directory. Wraps ResolveChain.
func ResolveChainDir(chainType, name, chainId string) (string, error) {
	var err error
	if name != "" {
		chainType, chainId, err = ChainFromName(name)
	} else if chainId != "" {
		chainType, chainId, err = ResolveChain(path.Join(chainType, chainId))
	}
	if err != nil {
		return "", err
	}
	return ComposeRoot(chainType, chainId), nil
}

// lookup chainIds by prefix match
func findPrefixMatch(dirPath, prefix string) (string, error) {
	fs, _ := ioutil.ReadDir(dirPath)
	found := false
	var p string
	for _, f := range fs {
		if strings.HasPrefix(f.Name(), prefix) {
			if found {
				return "", fmt.Errorf("ChainId collision! Multiple chains begin with %s. Please be more specific", prefix)
			}
			p = f.Name() //path.Join(Blockchains, chainType, f.Name())
			found = true
		}
	}
	if !found {
		return "", fmt.Errorf("ChainId %s did not match any known chains.", prefix)
	}
	return p, nil
}

// Maximum entries in the HEAD file
var MaxHead = 100

// Write a new entry (type/chainId) to the HEAD file
// Expects the chain type and head (id) to be full (already resolved)
func changeHead(typ, id string) error {
	s := fmt.Sprintf("%s/%s", typ, id)
	return ioutil.WriteFile(common.HEAD, []byte(s), 0666)
}

// Change the head to null (no head)
func NullHead() error {
	return changeHead("", "")
}

// Write a new entry to the HEAD file.
// Arguments are chain type and new head (chainId or ref name)
func ChangeHead(typ, id string) error {
	var err error
	typ, err = ResolveChainType(typ)
	if err != nil {
		return err
	}
	id, err = ResolveChainId(typ, id)
	if err != nil {
		return err
	}
	return changeHead(typ, id)
}

func addRef(typ, id, ref string) error {
	typ, err := ResolveChainType(typ)
	if err != nil {
		return err
	}

	dataDir := path.Join(common.BlockchainsPath, typ)
	_, err = os.Stat(path.Join(dataDir, id))
	if err != nil {
		id, err = findPrefixMatch(dataDir, id)
		if err != nil {
			return err
		}
	}

	refid := path.Join(typ, id)
	return ioutil.WriteFile(path.Join(common.Refs, ref), []byte(refid), 0644)
}

// Add a reference name to a chainId
func AddRef(typ, id, ref string) error {
	_, err := os.Stat(path.Join(common.Refs, ref))
	if err == nil {
		return fmt.Errorf("Ref %s already exists", ref)
	}
	return addRef(typ, id, ref)
}

func AddRefForce(typ, id, ref string) error {
	return addRef(typ, id, ref)
}

// Return a list of chain references
func GetRefs() (map[string]string, error) {
	fs, err := ioutil.ReadDir(common.Refs)
	if err != nil {
		return nil, err
	}
	m := make(map[string]string)
	for _, f := range fs {
		name := f.Name()
		b, err := ioutil.ReadFile(path.Join(common.Refs, name))
		if err != nil {
			return nil, err
		}
		m[name] = string(b)
	}
	return m, nil
}

// Get the current active chain (top of the HEAD file)
// Returns chain type and chain id
func GetHead() (string, string, error) {
	// TODO: only read the one line!
	f, err := ioutil.ReadFile(common.HEAD)
	if err != nil {
		return "", "", err
	}
	fspl := strings.Split(string(f), "\n")
	head := fspl[0]
	if head == "" {
		return "", "", fmt.Errorf("There is no chain checked out")
	}
	return SplitRef(head)
}

type ErrBadRef struct {
	ref string
}

func (e ErrBadRef) Error() string {
	return fmt.Sprintf("Improperly formatted ref: %s", e.ref)
}

func SplitRef(ref string) (string, string, error) {
	sp := strings.Split(ref, "/")
	if len(sp) != 2 {
		return "", "", ErrBadRef{ref}
	}
	return sp[0], sp[1], nil
}
