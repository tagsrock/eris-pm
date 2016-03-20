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

	log "github.com/eris-ltd/eris-pm/Godeps/_workspace/src/github.com/Sirupsen/logrus"
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

// Populates the parameters for packing according to the abi specification
type PackType struct {
	Name       string
	Type       string
	Raw        string
	Dynamic    bool
	Data       []byte
	DataLength []byte
	ArgNumber  int
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

	if len(data) == 0 {
		log.Debug("Nothing to pack")
	}

	var packer []*PackType
	var arguments []byte
	for i, a := range data {
		input := method.Inputs[i]
		thisPacked := &PackType{}

		thisPacked.Name = input.Name
		thisPacked.Type = input.Type.String()
		thisPacked.Raw = a
		thisPacked.ArgNumber = i

		log.WithFields(log.Fields{
			"name":   thisPacked.Name,
			"type":   thisPacked.Type,
			"val":    thisPacked.Raw,
			"argNum": thisPacked.ArgNumber,
		}).Debug("ABI Pack")
		err := PackProcessType(thisPacked)
		if err != nil {
			return nil, err
		}

		packer = append(packer, thisPacked)
	}

	arguments = ProcessPackedTypes(packer)

	// Set function id; final formulation of call
	packed := method.Id()
	packed = append(packed, arguments...)

	return packed, nil
}

// Order all the arguments as they should be ordered
func ProcessPackedTypes(packer []*PackType) []byte {
	var arguments []byte
	var argumentsData []byte

	// first we loop through and assemble the dynamic types
	for _, thisPacked := range packer {
		if !thisPacked.Dynamic {
			continue // these will get populated in the next range through
		} else {
			argumentsData = append(argumentsData, thisPacked.DataLength...)
			argumentsData = append(argumentsData, thisPacked.Data...)
		}
	}

	// second we loop through and find the length pointers (dynamic types) and raw data (static types)
	for i, thisPacked := range packer {
		if !thisPacked.Dynamic {
			arguments = append(arguments, thisPacked.Data...)
		} else {
			arguments = append(arguments, findOffset(packer, i)...)
		}
	}

	return append(arguments, argumentsData...)
}

// Conversion to []byte based on "Type"
// https://github.com/ethereum/wiki/wiki/Ethereum-Contract-ABI#formal-specification-of-the-encoding
func PackProcessType(thisPacked *PackType) error {
	t := getMajorType(thisPacked.Type)
	switch t {
	case "byte":
		thisPacked.Dynamic = false
		thisPacked.Data = common.RightPadBytes([]byte(thisPacked.Raw), lengths["retBlock"])
		return nil
	case "string":
		thisPacked.Dynamic = true
		thisPacked.DataLength = U2U256(uint64(len(thisPacked.Raw)))
		thisPacked.Data = common.RightPadBytes([]byte(thisPacked.Raw), lengths["retBlock"])
		return nil
	case "uint":
		thisPacked.Dynamic = false
		val, err := strconv.ParseUint(thisPacked.Raw, 10, 64)
		if err != nil {
			return err
		}
		thisPacked.Data = U2U256(val)
		return nil
	case "int":
		thisPacked.Dynamic = false
		val, err := strconv.ParseInt(thisPacked.Raw, 10, 64)
		if err != nil {
			return err
		}
		thisPacked.Data = S2S256(val)
		return nil
	case "address":
		thisPacked.Dynamic = false
		thisPacked.Data = common.LeftPadBytes(common.AddressStringToBytes(thisPacked.Raw), lengths["retBlock"])
		return nil
	case "bool":
		thisPacked.Dynamic = false
		if thisPacked.Raw == "1" || thisPacked.Raw == "true" {
			thisPacked.Data = common.LeftPadBytes(common.Big1.Bytes(), lengths["retBlock"])
			return nil
		} else {
			thisPacked.Data = common.LeftPadBytes(common.Big0.Bytes(), lengths["retBlock"])
			return nil
		}
	default:
		return fmt.Errorf("Unknown type. Cannot pack.")
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
			return nil, fmt.Errorf("Unrecognized return type (%s)", method.Outputs[i].Type.String())
		}

		next = start + lengths["retBlock"]
		if next > end {
			log.WithFields(log.Fields{
				"name":   ret[i].Name,
				"type":   ret[i].Type,
				"val":    ret[i].Value,
				"len":    lengths[method.Outputs[i].Type.String()],
				"retBlk": lengths["retBlock"],
				"start":  start,
				"next":   next,
				"end":    end,
			}).Error("Too little data")
			return nil, fmt.Errorf("Too little data; usually means the wrong abi was used")
		}

		ret[i].Name = method.Outputs[i].Name
		ret[i].Type = method.Outputs[i].Type.String()
		ret[i].Value, next = UnpackProcessType(ret[i].Type, data[start:next], start)
		log.WithFields(log.Fields{
			"name": ret[i].Name,
			"type": ret[i].Type,
			"val":  ret[i].Value,
		}).Debug("ABI Unpack")

		start = next
	}

	if start != end {
		log.WithFields(log.Fields{
			"name":   ret[len(ret)-1].Name,
			"type":   ret[len(ret)-1].Type,
			"val":    ret[len(ret)-1].Value,
			"len":    lengths[method.Outputs[len(ret)-1].Type.String()],
			"retBlk": lengths["retBlock"],
			"start":  start,
			"next":   next,
			"end":    end,
		}).Error("Too much data")
		return nil, fmt.Errorf("Too much data; usually this means difficulty in rendering the return from the contract")
	}

	retbytes, err := json.Marshal(ret)
	if err != nil {
		return nil, err
	}

	return retbytes, nil

}

//utility Functions

//Conversion to string based on "Type"
func UnpackProcessType(typ string, value []byte, start int) (string, int) {
	t := getMajorType(typ)
	switch t {
	case "byte":
		return string(common.UnRightPadBytes(value)), (start + lengths["retBlock"])
	case "string":
		return unpackByteArray(value, start)
	case "uint":
		val := common.StripZeros(common.BigD(value).String())
		if val == "" {
			return "0", (start + lengths["retBlock"])
		}
		return val, (start + lengths["retBlock"])
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
			return ("-" + val), (start + lengths["retBlock"])
		}

		val := common.StripZeros(common.BigD(value).String())

		// blank strings will be zero after all the decoding finishes
		if val == "" {
			return "0", (start + lengths["retBlock"])
		}

		return val, (start + lengths["retBlock"])
	case "address":
		return strings.ToUpper(hex.EncodeToString(common.Address(value))), (start + lengths["retBlock"])
	case "bool":
		return new(big.Int).SetBytes(value).String(), (start + lengths["retBlock"])
	default:
		return hex.EncodeToString(value), (start + lengths["retBlock"])
	}
}

func unpackByteArray(value []byte, start int) (string, int) {
	var next int

	// first we get the bytes delimiter
	next = start + lengths["retBlock"]
	v1 := value[start:next]
	delim, _ := strconv.Atoi(common.StripZeros(common.BigD(v1).String()))
	start = next

	// next we get the length of string
	next = start + lengths["retBlock"]
	v2 := value[start:next]
	length, _ := strconv.Atoi(common.StripZeros(common.BigD(v2).String()))
	start = next

	// now we unmarshall
	blocks := ((length - 1) / delim) + 1          // number of "chunks to use"
	next = start + (lengths["retBlock"] * blocks) // how much of the byte array to use
	var val string
	if length%delim == 0 {
		val = string(value[start:next]) // strings exactly length%32==0 do not marshall with UnRightPad
	} else {
		val = string(common.UnRightPadBytes(value[start:next]))
	}
	log.WithFields(log.Fields{
		"len":      length,
		"delimter": delim,
		"val":      val,
	}).Debug("Strings unpack")
	return val, next
}

func findOffset(packer []*PackType, i int) []byte {
	offset := len(packer) // number of arguments

	if i == 0 {
		return U2U256(uint64(32 * offset))
	}
	for i = i - 1; i >= 0; i-- { // loop thru backwards
		if packer[i].Dynamic {
			offset++                                     // data length slot
			offset = offset + (len(packer[i].Data) / 32) // length of additional slots
		}
	}

	return U2U256(uint64(32 * offset))
}

// ________ to we used these? ---------------------|

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
	if typ == "bytes" {
		return "string"
	}
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
