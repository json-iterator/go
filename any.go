package jsoniter

import (
	"fmt"
	"reflect"
)

type Any struct {
	Val   interface{}
	Error error
}

func (any *Any) GetObject(keys ...interface{}) interface{} {
	ret, err := getPath(any.Val, keys...)
	if err != nil {
		any.Error = err
		return "";
	}
	return ret
}

func (any *Any) GetString(keys ...interface{}) string {
	ret, err := getPath(any.Val, keys...)
	if err != nil {
		any.Error = err
		return "";
	}
	typedRet, ok := ret.(string)
	if !ok {
		any.Error = fmt.Errorf("%v is not string", ret);
		return "";
	}
	return typedRet
}

func (any *Any) GetUint8(keys ...interface{}) uint8 {
	ret, err := getPathAsInt64(any.Val, keys...)
	if err != nil {
		any.Error = err
		return 0;
	}
	return uint8(ret)
}

func (any *Any) GetInt8(keys ...interface{}) int8 {
	ret, err := getPathAsInt64(any.Val, keys...)
	if err != nil {
		any.Error = err
		return 0;
	}
	return int8(ret)
}

func (any *Any) GetUint16(keys ...interface{}) uint16 {
	ret, err := getPathAsInt64(any.Val, keys...)
	if err != nil {
		any.Error = err
		return 0;
	}
	return uint16(ret)
}

func (any *Any) GetInt16(keys ...interface{}) int16 {
	ret, err := getPathAsInt64(any.Val, keys...)
	if err != nil {
		any.Error = err
		return 0;
	}
	return int16(ret)
}

func (any *Any) GetUint32(keys ...interface{}) uint32 {
	ret, err := getPathAsInt64(any.Val, keys...)
	if err != nil {
		any.Error = err
		return 0;
	}
	return uint32(ret)
}

func (any *Any) GetInt32(keys ...interface{}) int32 {
	ret, err := getPathAsInt64(any.Val, keys...)
	if err != nil {
		any.Error = err
		return 0;
	}
	return int32(ret)
}

func (any *Any) GetUint64(keys ...interface{}) uint64 {
	ret, err := getPathAsUint64(any.Val, keys...)
	if err != nil {
		any.Error = err
		return 0;
	}
	return uint64(ret)
}

func (any *Any) GetInt64(keys ...interface{}) int64 {
	ret, err := getPathAsInt64(any.Val, keys...)
	if err != nil {
		any.Error = err
		return 0;
	}
	return int64(ret)
}

func (any *Any) GetInt(keys ...interface{}) int {
	ret, err := getPathAsInt64(any.Val, keys...)
	if err != nil {
		any.Error = err
		return 0;
	}
	return int(ret)
}

func (any *Any) GetUint(keys ...interface{}) uint {
	ret, err := getPathAsInt64(any.Val, keys...)
	if err != nil {
		any.Error = err
		return 0;
	}
	return uint(ret)
}

func (any *Any) GetFloat32(keys ...interface{}) float32 {
	ret, err := getPathAsFloat64(any.Val, keys...)
	if err != nil {
		any.Error = err
		return 0;
	}
	return float32(ret)
}

func (any *Any) GetFloat64(keys ...interface{}) float64 {
	ret, err := getPathAsFloat64(any.Val, keys...)
	if err != nil {
		any.Error = err
		return 0;
	}
	return ret
}

func (any *Any) GetBool(keys ...interface{}) bool {
	ret, err := getPath(any.Val, keys...)
	if err != nil {
		any.Error = err
		return false;
	}
	typedRet, ok := ret.(bool)
	if !ok {
		any.Error = fmt.Errorf("%v is not bool", ret)
		return false;
	}
	return typedRet
}

func (any *Any) IsNull(keys ...interface{}) bool {
	ret, err := getPath(any.Val, keys...)
	if err != nil {
		any.Error = err
		return false;
	}
	return reflect.ValueOf(ret).IsNil()
}

func getPathAsInt64(val interface{}, keys ...interface{}) (int64, error) {
	ret, err := getPath(val, keys...)
	if err != nil {
		return 0, err
	}
	switch ret := ret.(type) {
	case uint8:
		return int64(ret), nil;
	case int8:
		return int64(ret), nil;
	case uint16:
		return int64(ret), nil;
	case int16:
		return int64(ret), nil;
	case uint32:
		return int64(ret), nil;
	case int32:
		return int64(ret), nil;
	case uint64:
		return int64(ret), nil;
	case int64:
		return int64(ret), nil;
	case int:
		return int64(ret), nil;
	case uint:
		return int64(ret), nil;
	case float32:
		return int64(ret), nil;
	case float64:
		return int64(ret), nil;
	default:
		return 0, fmt.Errorf("%v is not number", ret)
	}
}

func getPathAsUint64(val interface{}, keys ...interface{}) (uint64, error) {
	ret, err := getPath(val, keys...)
	if err != nil {
		return 0, err
	}
	switch ret := ret.(type) {
	case uint8:
		return uint64(ret), nil;
	case int8:
		return uint64(ret), nil;
	case uint16:
		return uint64(ret), nil;
	case int16:
		return uint64(ret), nil;
	case uint32:
		return uint64(ret), nil;
	case int32:
		return uint64(ret), nil;
	case uint64:
		return uint64(ret), nil;
	case int64:
		return uint64(ret), nil;
	case int:
		return uint64(ret), nil;
	case uint:
		return uint64(ret), nil;
	case float32:
		return uint64(ret), nil;
	case float64:
		return uint64(ret), nil;
	default:
		return 0, fmt.Errorf("%v is not number", ret)
	}
}

func getPathAsFloat64(val interface{}, keys ...interface{}) (float64, error) {
	ret, err := getPath(val, keys...)
	if err != nil {
		return 0, err
	}
	switch ret := ret.(type) {
	case uint8:
		return float64(ret), nil;
	case int8:
		return float64(ret), nil;
	case uint16:
		return float64(ret), nil;
	case int16:
		return float64(ret), nil;
	case uint32:
		return float64(ret), nil;
	case int32:
		return float64(ret), nil;
	case uint64:
		return float64(ret), nil;
	case int64:
		return float64(ret), nil;
	case int:
		return float64(ret), nil;
	case uint:
		return float64(ret), nil;
	case float32:
		return float64(ret), nil;
	case float64:
		return float64(ret), nil;
	default:
		return 0, fmt.Errorf("%v is not number", ret)
	}
}

func getPath(val interface{}, keys ...interface{}) (interface{}, error) {
	if (len(keys) == 0) {
		return val, nil;
	}
	switch key := keys[0].(type) {
	case string:
		nextVal, err := getFromMap(val, key)
		if err != nil {
			return nil, err
		}
		nextKeys := make([]interface{}, len(keys) - 1)
		copy(nextKeys, keys[1:])
		return getPath(nextVal, nextKeys...)
	case int:
		nextVal, err := getFromArray(val, key)
		if err != nil {
			return nil, err
		}
		nextKeys := make([]interface{}, len(keys) - 1)
		copy(nextKeys, keys[1:])
		return getPath(nextVal, nextKeys...)
	default:
		return nil, fmt.Errorf("%v is not string or int", keys[0]);
	}
	return getPath(val, keys);
}

func getFromMap(val interface{}, key string) (interface{}, error) {
	mapVal, ok := val.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("%v is not map[string]interface{}", val)
	}
	ret, found := mapVal[key]
	if !found {
		return nil, fmt.Errorf("%v not found in %v", key, mapVal)
	}
	return ret, nil
}

func getFromArray(val interface{}, key int) (interface{}, error) {
	arrayVal, ok := val.([]interface{})
	if !ok {
		return nil, fmt.Errorf("%v is not []interface{}", val)
	}
	if key >= len(arrayVal) {
		return nil, fmt.Errorf("%v exceed %v", key, arrayVal)
	}
	if key < 0 {
		return nil, fmt.Errorf("%v exceed %v", key, arrayVal)
	}
	return arrayVal[key], nil
}
