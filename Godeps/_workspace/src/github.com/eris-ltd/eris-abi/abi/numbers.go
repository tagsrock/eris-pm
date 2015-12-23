package abi

import (
	"math/big"
	"reflect"

	"github.com/eris-ltd/eris-pm/Godeps/_workspace/src/github.com/eris-ltd/common/go/common"
)

var big_t = reflect.TypeOf(&big.Int{})
var ubig_t = reflect.TypeOf(&big.Int{})
var byte_t = reflect.TypeOf(byte(0))
var byte_ts = reflect.TypeOf([]byte(nil))
var uint_t = reflect.TypeOf(uint(0))
var uint8_t = reflect.TypeOf(uint8(0))
var uint16_t = reflect.TypeOf(uint16(0))
var uint32_t = reflect.TypeOf(uint32(0))
var uint64_t = reflect.TypeOf(uint64(0))
var int_t = reflect.TypeOf(int(0))
var int8_t = reflect.TypeOf(int8(0))
var int16_t = reflect.TypeOf(int16(0))
var int32_t = reflect.TypeOf(int32(0))
var int64_t = reflect.TypeOf(int64(0))

var uint_ts = reflect.TypeOf([]uint(nil))
var uint8_ts = reflect.TypeOf([]uint8(nil))
var uint16_ts = reflect.TypeOf([]uint16(nil))
var uint32_ts = reflect.TypeOf([]uint32(nil))
var uint64_ts = reflect.TypeOf([]uint64(nil))
var ubig_ts = reflect.TypeOf([]*big.Int(nil))

var int_ts = reflect.TypeOf([]int(nil))
var int8_ts = reflect.TypeOf([]int8(nil))
var int16_ts = reflect.TypeOf([]int16(nil))
var int32_ts = reflect.TypeOf([]int32(nil))
var int64_ts = reflect.TypeOf([]int64(nil))
var big_ts = reflect.TypeOf([]*big.Int(nil))

func U2U256(n uint64) []byte {
	return U256(big.NewInt(int64(n)))
}

// U256 will ensure unsigned 256bit on big nums
func U256(n *big.Int) []byte {
	return common.LeftPadBytes(common.U256(n).Bytes(), 32)
}

// S256 will ensure signed 256bit on big nums
func S2S256(n int64) []byte {
	return S256(big.NewInt(n))
}

func S256(n *big.Int) []byte {
	sint := common.S256(n)
	ret := common.LeftPadBytes(sint.Bytes(), 32)
	if sint.Cmp(common.Big0) < 0 {
		for i, b := range ret {
			if b == 0 {
				ret[i] = 1
				continue
			}
			break
		}
	}

	return ret
}
