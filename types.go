package schema

import (
	"fmt"
	"reflect"
	"strconv"
)

func BuiltinTypes() []Type {
	return []Type{
		boolType{},
		intType{},
		int8Type{},
		int16Type{},
		int32Type{},
		int64Type{},
		uintType{},
		uint8Type{},
		uint16Type{},
		uint32Type{},
		uint64Type{},
		float32Type{},
		float64Type{},
		stringType{},
	}
}

type boolType struct{}

func (boolType) DataType() interface{} { return false }

func (boolType) Decode(s string) (val interface{}, err error) {
	v, err := strconv.ParseBool(s)
	if err != nil {
		return nil, fmt.Errorf("invalid value(bool): %s", s)
	}
	val = v
	return
}
func (boolType) Encode(val interface{}) (s string, err error) {
	v, ok := val.(bool)
	if !ok {
		return "", fmt.Errorf("invalid data type, expect bool, but got %s", reflect.TypeOf(val))
	}
	return strconv.FormatBool(v), nil
}

type intType struct{}

func (intType) DataType() interface{} { return int(0) }

func (intType) Decode(s string) (val interface{}, err error) {
	v, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid value(int): %s", s)
	}
	val = int(v)
	return
}

func (intType) Encode(val interface{}) (s string, err error) {
	v, ok := val.(int)
	if !ok {
		return "", fmt.Errorf("invalid data type, expe ctint, but got %s", reflect.TypeOf(val))
	}
	return strconv.FormatInt(int64(v), 10), nil
}

type int8Type struct{}

func (int8Type) DataType() interface{} { return int8(0) }

func (int8Type) Decode(s string) (val interface{}, err error) {
	v, err := strconv.ParseInt(s, 10, 8)
	if err != nil {
		return nil, fmt.Errorf("invalid value(int8): %s", s)
	}
	val = int8(v)
	return
}

func (int8Type) Encode(val interface{}) (s string, err error) {
	v, ok := val.(int8)
	if !ok {
		return "", fmt.Errorf("invalid data type, expec tint8, but got %s", reflect.TypeOf(val))
	}
	return strconv.FormatInt(int64(v), 10), nil
}

type int16Type struct{}

func (int16Type) DataType() interface{} { return int16(0) }

func (int16Type) Decode(s string) (val interface{}, err error) {
	v, err := strconv.ParseInt(s, 10, 16)
	if err != nil {
		return nil, fmt.Errorf("invalid value(int16): %s", s)
	}
	val = int16(v)
	return
}

func (int16Type) Encode(val interface{}) (s string, err error) {
	v, ok := val.(int16)
	if !ok {
		return "", fmt.Errorf("invalid data type, expect int16, but got %s", reflect.TypeOf(val))
	}
	return strconv.FormatInt(int64(v), 10), nil
}

type int32Type struct{}

func (int32Type) DataType() interface{} { return int32(0) }

func (int32Type) Decode(s string) (val interface{}, err error) {
	v, err := strconv.ParseInt(s, 10, 32)
	if err != nil {
		return nil, fmt.Errorf("invalid value(int32): %s", s)
	}
	val = int32(v)
	return
}

func (int32Type) Encode(val interface{}) (s string, err error) {
	v, ok := val.(int32)
	if !ok {
		return "", fmt.Errorf("invalid data type, expect int32, but got %s", reflect.TypeOf(val))
	}
	return strconv.FormatInt(int64(v), 10), nil
}

type int64Type struct{}

func (int64Type) DataType() interface{} { return int64(0) }

func (int64Type) Decode(s string) (val interface{}, err error) {
	v, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid value(int64): %s", s)
	}
	val = int64(v)
	return
}

func (int64Type) Encode(val interface{}) (s string, err error) {
	v, ok := val.(int64)
	if !ok {
		return "", fmt.Errorf("invalid data type, expect int64, but got %s", reflect.TypeOf(val))
	}
	return strconv.FormatInt(v, 10), nil
}

type uintType struct{}

func (uintType) DataType() interface{} { return uint(0) }

func (uintType) Decode(s string) (val interface{}, err error) {
	v, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid value(uint): %s", s)
	}
	val = uint(v)
	return
}

func (uintType) Encode(val interface{}) (s string, err error) {
	v, ok := val.(uint)
	if !ok {
		return "", fmt.Errorf("invalid data type, expect  uint, but got %s", reflect.TypeOf(val))
	}
	return strconv.FormatUint(uint64(v), 10), nil
}

type uint8Type struct{}

func (uint8Type) DataType() interface{} { return uint8(0) }

func (uint8Type) Decode(s string) (val interface{}, err error) {
	v, err := strconv.ParseUint(s, 10, 8)
	if err != nil {
		return nil, fmt.Errorf("invalid value(uint8): %s", s)
	}
	val = uint8(v)
	return
}
func (uint8Type) Encode(val interface{}) (s string, err error) {
	v, ok := val.(uint8)
	if !ok {
		return "", fmt.Errorf("invalid data type, expect uint8, but got %s", reflect.TypeOf(val))
	}
	return strconv.FormatUint(uint64(v), 10), nil
}

type uint16Type struct{}

func (uint16Type) DataType() interface{} { return uint16(0) }

func (uint16Type) Decode(s string) (val interface{}, err error) {
	v, err := strconv.ParseUint(s, 10, 16)
	if err != nil {
		return nil, fmt.Errorf("invalid value(uint16): %s", s)
	}
	val = uint16(v)
	return
}
func (uint16Type) Encode(val interface{}) (s string, err error) {
	v, ok := val.(uint16)
	if !ok {
		return "", fmt.Errorf("invalid data type, expect uint16, but got %s", reflect.TypeOf(val))
	}
	return strconv.FormatUint(uint64(v), 10), nil
}

type uint32Type struct{}

func (uint32Type) DataType() interface{} { return uint32(0) }

func (uint32Type) Decode(s string) (val interface{}, err error) {
	v, err := strconv.ParseUint(s, 10, 32)
	if err != nil {
		return nil, fmt.Errorf("invalid value(uint32): %s", s)
	}
	val = uint32(v)
	return
}
func (uint32Type) Encode(val interface{}) (s string, err error) {
	v, ok := val.(uint32)
	if !ok {
		return "", fmt.Errorf("invalid data type, expect uint32, but got %s", reflect.TypeOf(val))
	}
	return strconv.FormatUint(uint64(v), 10), nil
}

type uint64Type struct{}

func (uint64Type) DataType() interface{} { return uint64(0) }

func (uint64Type) Decode(s string) (val interface{}, err error) {
	v, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid value(uint64): %s", s)
	}
	val = uint64(v)
	return
}
func (uint64Type) Encode(val interface{}) (s string, err error) {
	v, ok := val.(uint64)
	if !ok {
		return "", fmt.Errorf("invalid data type, expect uint64, but got %s", reflect.TypeOf(val))
	}
	return strconv.FormatUint(uint64(v), 10), nil
}

type float32Type struct{}

func (float32Type) DataType() interface{} { return float32(0) }

func (float32Type) Decode(s string) (val interface{}, err error) {
	v, err := strconv.ParseFloat(s, 32)
	if err != nil {
		return nil, fmt.Errorf("invalid value(float32): %s", s)
	}
	val = float32(v)
	return
}
func (float32Type) Encode(val interface{}) (s string, err error) {
	v, ok := val.(float32)
	if !ok {
		return "", fmt.Errorf("invalid data type, expect float32, but got %s", reflect.TypeOf(val))
	}
	return strconv.FormatFloat(float64(v), 'f', -1, 32), nil
}

type float64Type struct{}

func (float64Type) DataType() interface{} { return float64(0) }

func (float64Type) Decode(s string) (val interface{}, err error) {
	v, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid value(float64): %s", s)
	}
	val = float64(v)
	return
}
func (float64Type) Encode(val interface{}) (s string, err error) {
	v, ok := val.(float64)
	if !ok {
		return "", fmt.Errorf("invalid data type, expect float64, but got %s", reflect.TypeOf(val))
	}
	return strconv.FormatFloat(float64(v), 'f', -1, 64), nil
}

type stringType struct{}

func (stringType) DataType() interface{} { return "" }

func (stringType) Decode(v string) (val interface{}, err error) {
	val = v
	return
}
func (stringType) Encode(val interface{}) (s string, err error) {
	v, ok := val.(string)
	if !ok {
		return "", fmt.Errorf("invalid data type, expect string, but got %s", reflect.TypeOf(val))
	}
	return v, nil
}
