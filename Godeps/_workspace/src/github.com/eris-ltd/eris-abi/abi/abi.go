package abi

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"regexp"
	"strconv"
	"strings"

	"github.com/eris-ltd/eris-pm/Godeps/_workspace/src/github.com/eris-ltd/common/go/common"
	"github.com/eris-ltd/eris-pm/Godeps/_workspace/src/github.com/ethereum/go-ethereum/crypto/sha3"
)

var NullABI = ABI{}

// Callable method given a `Name` and whether the method is a constant.
// If the method is `Const` no transaction needs to be created for this
// particular Method call. It can easily be simulated using a local VM.
// For example a `Balance()` method only needs to retrieve something
// from the storage and therefor requires no Tx to be send to the
// network. A method such as `Transact` does require a Tx and thus will
// be flagged `true`.
// Inputs specifies the required input parameters for this gives method.
type Method struct {
	Name     string // `json:"name"`
	Constant bool
	Inputs   []Argument // `json:"inputs"`
	Outputs  []Argument // `json:"outputs"`
	Type     string
}

type Argpairs struct {
	Name  string
	Type  string
	Value string
	VByte []byte
}

// Argument holds the name of the argument and the corresponding type.
// Types are used when packing and testing arguments.
type Argument struct {
	Name string `json:"name"`
	Type Type   `json:"type"`
}

// The ABI holds information about a contract's context and available
// invokable methods. It will allow you to type check function calls and
// packs data accordingly.
type ABI struct {
	Methods map[string]Method
}

// Returns the methods string signature according to the ABI spec.
//
// Example
//
//     function foo(uint32 a, int b)    =    "foo(uint32,int256)"
//
// Please note that "int" is substitute for its canonical representation "int256"
func (m Method) String() (out string) {
	if strings.Contains(m.Name, "(") && strings.Contains(m.Name, ")") {
		return m.Name
	}
	out += m.Name
	types := make([]string, len(m.Inputs))
	i := 0
	for _, input := range m.Inputs {
		types[i] = input.Type.String()
		i++
	}
	out += "(" + strings.Join(types, ",") + ")"

	return
}

func (m Method) Id() []byte {
	return Sha3([]byte(m.String()))[:4]
}

// Pack the given method name to conform the ABI. Method call's data
// will consist of method_id, args0, arg1, ... argN. Method id consists
// of 4 bytes and arguments are all 32 bytes.
// Method ids are created from the first 4 bytes of the hash of the
// methods string signature. (signature = baz(uint32,string32))
func (abi ABI) Pack(name string, data []string) ([]byte, error) {
	method, exist := abi.Methods[name]
	if !exist {
		return nil, fmt.Errorf("method '%s' not found", name)
	}

	// start with argument count match
	if len(data) != len(method.Inputs) {
		return nil, fmt.Errorf("argument count mismatch: %d for %d", len(data), len(method.Inputs))
	}

	var arguments []byte
	for i, a := range data {
		input := method.Inputs[i]

		logger.Debugf("ABI Pack. Name =>\t\t%s\n", input.Name)
		logger.Debugf("ABI Pack. Type =>\t\t%s\n", input.Type.String())
		logger.Debugf("ABI Pack. Value =>\t\t%s\n", a)
		ret, err := PackProcessType(input.Type.String(), a)
		if err != nil {
			return nil, err
		}

		arguments = append(arguments, ret...)
	}

	// Set function id; final formulation of call
	packed := method.Id()
	packed = append(packed, arguments...)

	return packed, nil
}

//Conversion to []byte based on "Type"
func PackProcessType(typ string, value string) ([]byte, error) {
	t := getMajorType(typ)
	switch t {
	case "byte", "string":
		return common.RightPadBytes([]byte(value), lengths["retBlock"]), nil
	case "uint":
		val, err := strconv.ParseUint(value, 10, 64)
		if err != nil {
			return nil, err
		}
		return U2U256(val), nil
	case "int":
		val, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return nil, err
		}
		return S2S256(val), nil
	case "address":
		return common.LeftPadBytes(common.AddressStringToBytes(value), lengths["retBlock"]), nil
	case "bool":
		if value == "1" {
			return common.LeftPadBytes(common.Big1.Bytes(), lengths["retBlock"]), nil
		} else {
			return common.LeftPadBytes(common.Big0.Bytes(), lengths["retBlock"]), nil
		}
	default:
		return nil, fmt.Errorf("Unknown type. Cannot pack.")
	}
}

//Unpacking function
func (abi ABI) UnPack(name string, data []byte) ([]byte, error) {
	method, exist := abi.Methods[name]
	if !exist {
		return nil, fmt.Errorf("method '%s' not found", name)
	}

	ret := make([]Argpairs, len(method.Outputs))

	//Note this assumes all return values are 32 bytes (If this is not correct, process type should return number of bytes consumed?)
	start := 0
	var next int
	end := len(data)

	for i := range method.Outputs {

		_, ok := lengths[method.Outputs[i].Type.String()]
		if !ok {
			return nil, fmt.Errorf("Unrecognized return type")
		}

		next = start + lengths["retBlock"]
		if next > end {
			return nil, fmt.Errorf("Too little data")
		}

		ret[i].Name = method.Outputs[i].Name
		ret[i].Type = method.Outputs[i].Type.String()
		ret[i].Value = UnpackProcessType(ret[i].Type, data[start:next])
		logger.Debugf("ABI Unpack. Name =>\t\t%s\n", ret[i].Name)
		logger.Debugf("ABI Unpack. Type =>\t\t%s\n", ret[i].Type)
		logger.Debugf("ABI Unpack. Value =>\t\t%s\n", ret[i].Value)

		start = next
	}

	if start != end {
		return nil, fmt.Errorf("Too much data")
	}

	retbytes, err := json.Marshal(ret)
	if err != nil {
		return nil, err
	}

	return retbytes, nil

}

//utility Functions

//Conversion to string based on "Type"
func UnpackProcessType(typ string, value []byte) string {
	t := getMajorType(typ)
	switch t {
	case "byte", "string":
		return string(common.UnRightPadBytes(value))
	case "uint":
		val := common.StripZeros(common.BigD(value).String())
		if val == "" {
			return "0"
		}
		return val
	case "int":
		// there is weird encoding solidity does with negative ints
		//   it prepends with 01 instead of 00.
		if value[0] == 1 {
			i := 0
			for ; i < len(value); i++ {
				if value[i] != 1 {
					break
				}
			}
			val := common.BigD(value[i:]).String()
			return ("-" + val)
		}

		val := common.StripZeros(common.BigD(value).String())

		// blank strings will be zero after all the decoding finishes
		if val == "" {
			return "0"
		}

		return val
	case "address":
		return strings.ToUpper(hex.EncodeToString(common.Address(value)))
	case "bool":
		return new(big.Int).SetBytes(value).String()
	default:
		return hex.EncodeToString(value)
	}
}

func UnpackPrettyPrint(injson []byte) (string, error) {
	var ret []Argpairs

	err := json.Unmarshal(injson, &ret)
	if err != nil {
		return "", err
	}

	//Pretty print time
	pps := ""
	unc := int(1)
	for _, A := range ret {
		if A.Name == "" {
			tname := "UVar" + strconv.Itoa(unc)
			pps = pps + tname + " : " + A.Value
			unc = unc + 1
		} else {
			pps = pps + A.Name + " : " + A.Value
		}
		pps = pps + "\n"
	}

	return pps, nil
}

//Fills an ABI object with umarshalled data.
func (abi *ABI) UnmarshalJSON(data []byte) error {
	var methods []Method
	if err := json.Unmarshal(data, &methods); err != nil {
		return err
	}

	abi.Methods = make(map[string]Method)
	for _, method := range methods {
		abi.Methods[method.Name] = method
	}

	return nil
}

func (a *Argument) UnmarshalJSON(data []byte) error {
	var extarg struct {
		Name string
		Type string
	}
	err := json.Unmarshal(data, &extarg)
	if err != nil {
		return fmt.Errorf("argument json err: %v", err)
	}

	a.Type, err = NewType(extarg.Type)
	if err != nil {
		return err
	}
	a.Name = extarg.Name

	return nil
}

func JSON(reader io.Reader) (ABI, error) {
	dec := json.NewDecoder(reader)

	var abi ABI
	if err := dec.Decode(&abi); err != nil {
		return ABI{}, err
	}

	return abi, nil
}

func Sha3(data []byte) []byte {
	d := sha3.NewKeccak256()
	d.Write(data)

	return d.Sum(nil)
}

func getMajorType(typ string) string {
	var t bool
	t, _ = regexp.MatchString("byte", typ)
	if t {
		return "byte"
	}
	t, _ = regexp.MatchString("string", typ)
	if t {
		return "string"
	}
	t, _ = regexp.MatchString("uint", typ) // Test uint first because int will also match uint
	if t {
		return "uint"
	}
	t, _ = regexp.MatchString("int", typ)
	if t {
		return "int"
	}
	t, _ = regexp.MatchString("address", typ)
	if t {
		return "address"
	}
	t, _ = regexp.MatchString("bool", typ)
	if t {
		return "bool"
	}
	return "unknown"
}
