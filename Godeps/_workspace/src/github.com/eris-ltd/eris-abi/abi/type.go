package abi

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
)

const (
	IntTy byte = iota
	UintTy
	BoolTy
	SliceTy
	AddressTy
	RealTy
)

// Type is the reflection of the supported argument type
type Type struct {
	Kind       reflect.Kind
	Type       reflect.Type
	Size       int
	T          byte   // Our own type checking
	isSlice    bool
	stringKind string // holds the unparsed string for deriving signatures
	baseType   string // holds our base type for slices
}

func (t Type) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`"%s"`, t.stringKind)), nil
}

// New type returns a fully parsed Type given by the input string or an error if it  can't be parsed.
//
// Strings can be in the format of:
//
// 	Input  = Type [ "[" [ Number ] "]" ] Name .
// 	Type   = [ "u" ] "int" [ Number ] .
//
// Examples:
//
//      string     int       uint       fixed
//      string32   int8      uint8      uint[]
//      address    int256    uint256    fixed[2]
func NewType(t string) (typ Type, err error) {
	// parse eg. uint32 || uint32[] || uint32[20]
	// 1. full string 2. type 3. (opt.) is slice 4. (opt.) size
	freg, err := regexp.Compile("([a-zA-Z0-9]+)(\\[([0-9]*)?\\])?")
	if err != nil {
		return Type{}, err
	}
	res := freg.FindAllStringSubmatch(t, -1)[0]
	var (
		isslice bool
		size    int
	)
	switch {
	case res[3] != "":
		// err is ignored. Already checked for number through the regexp
		size, _ = strconv.Atoi(res[3])
		isslice = true
	case res[2] != "":
		isslice = true
		size = -1
	case res[0] == "":
		return Type{}, fmt.Errorf("type parse error for `%s`", t)
	}

	// parse eg. uint32 || uint
	treg, err := regexp.Compile("([a-zA-Z]+)([0-9]*)?")
	if err != nil {
		return Type{}, err
	}

	parsedType := treg.FindAllStringSubmatch(res[1], -1)[0]
	vsize, _ := strconv.Atoi(parsedType[2])
	vtype := parsedType[1]
	// substitute canonical representation
	if vsize == 0 && (vtype == "int" || vtype == "uint") {
		vsize = 256
		t += "256"
	}

	if isslice {
		typ.Kind = reflect.Slice
		typ.Size = size
		typ.isSlice = true
		switch vtype {
		case "int", "bool":
			typ.Type = big_ts
		case "uint":
			typ.Type = ubig_ts
		case "address", "bytes":
			typ.Type = byte_ts
		default:
			return Type{}, fmt.Errorf("unsupported arg slice type: %s", t)
		}
	} else {
		typ.isSlice = false
		typ.Size = 1
		switch vtype {
		case "int":
			typ.Kind = reflect.Ptr
			typ.Type = big_t
			typ.T = IntTy
		case "uint":
			typ.Kind = reflect.Ptr
			typ.Type = ubig_t
			typ.T = UintTy
		case "bool":
			typ.Kind = reflect.Bool
		case "fixed": // TODO
			typ.Kind = reflect.Invalid
		case "address":
			typ.Kind = reflect.Slice
			typ.Type = byte_ts
			typ.T = AddressTy
		case "string", "bytes":
			typ.Kind = reflect.String
		default:
			return Type{}, fmt.Errorf("unsupported arg type: %s", t)
		}
	}
	typ.stringKind = t
	typ.baseType = res[1]

	return
}

func (t Type) String() (out string) {
	return t.stringKind
}

func (t Type) BaseType() (out string) {
	return t.baseType
}
