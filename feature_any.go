package jsoniter

import (
	"fmt"
	"reflect"
	"strconv"
)

type Any struct {
	val   interface{}
	Error error
	LastAccessed interface{}
}

func MakeAny(val interface{}) *Any {
	return &Any{val, nil, nil}
}

func (any *Any) Get(keys ...interface{}) interface{} {
	ret, err := getPath(any.val, keys...)
	any.LastAccessed = ret
	if err != nil {
		any.Error = err
		return "";
	}
	return ret
}

func (any *Any) GetValueType(keys ...interface{}) ValueType {
	ret, err := getPath(any.val, keys...)
	any.LastAccessed = ret
	if err != nil {
		any.Error = err
		return Invalid;
	}

	switch reflect.TypeOf(ret).Kind() {
	case reflect.Uint8:
		return Number;
	case reflect.Int8:
		return Number;
	case reflect.Uint16:
		return Number;
	case reflect.Int16:
		return Number;
	case reflect.Uint32:
		return Number;
	case reflect.Int32:
		return Number;
	case reflect.Uint64:
		return Number;
	case reflect.Int64:
		return Number;
	case reflect.Int:
		return Number;
	case reflect.Uint:
		return Number;
	case reflect.Float32:
		return Number;
	case reflect.Float64:
		return Number;
	case reflect.String:
		return String;
	case reflect.Bool:
		return Bool;
	case reflect.Array:
		return Array;
	case reflect.Struct:
		return Object;
	default:
		return Invalid
	}
}

func (any *Any) ToString(keys ...interface{}) string {
	ret, err := getPath(any.val, keys...)
	any.LastAccessed = ret
	if err != nil {
		any.Error = err
		return "";
	}
	switch ret := ret.(type) {
	case uint8:
		return strconv.FormatInt(int64(ret), 10);
	case int8:
		return strconv.FormatInt(int64(ret), 10);
	case uint16:
		return strconv.FormatInt(int64(ret), 10);
	case int16:
		return strconv.FormatInt(int64(ret), 10);
	case uint32:
		return strconv.FormatInt(int64(ret), 10);
	case int32:
		return strconv.FormatInt(int64(ret), 10);
	case uint64:
		return strconv.FormatUint(uint64(ret), 10);
	case int64:
		return strconv.FormatInt(int64(ret), 10);
	case int:
		return strconv.FormatInt(int64(ret), 10);
	case uint:
		return strconv.FormatInt(int64(ret), 10);
	case float32:
		return strconv.FormatFloat(float64(ret), 'E', -1, 32);
	case float64:
		return strconv.FormatFloat(ret, 'E', -1, 64);
	case string:
		return ret
	default:
		return fmt.Sprintf("%v", ret)
	}
}

func (any *Any) ToUint8(keys ...interface{}) uint8 {
	ret, err := getPathAsInt64(any, keys...)
	if err != nil {
		any.Error = err
		return 0;
	}
	return uint8(ret)
}

func (any *Any) ToInt8(keys ...interface{}) int8 {
	ret, err := getPathAsInt64(any, keys...)
	if err != nil {
		any.Error = err
		return 0;
	}
	return int8(ret)
}

func (any *Any) ToUint16(keys ...interface{}) uint16 {
	ret, err := getPathAsInt64(any, keys...)
	if err != nil {
		any.Error = err
		return 0;
	}
	return uint16(ret)
}

func (any *Any) ToInt16(keys ...interface{}) int16 {
	ret, err := getPathAsInt64(any, keys...)
	if err != nil {
		any.Error = err
		return 0;
	}
	return int16(ret)
}

func (any *Any) ToUint32(keys ...interface{}) uint32 {
	ret, err := getPathAsInt64(any, keys...)
	if err != nil {
		any.Error = err
		return 0;
	}
	return uint32(ret)
}

func (any *Any) ToInt32(keys ...interface{}) int32 {
	ret, err := getPathAsInt64(any, keys...)
	if err != nil {
		any.Error = err
		return 0;
	}
	return int32(ret)
}

func (any *Any) ToUint64(keys ...interface{}) uint64 {
	ret, err := getPathAsUint64(any, keys...)
	if err != nil {
		any.Error = err
		return 0;
	}
	return uint64(ret)
}

func (any *Any) ToInt64(keys ...interface{}) int64 {
	ret, err := getPathAsInt64(any, keys...)
	if err != nil {
		any.Error = err
		return 0;
	}
	return int64(ret)
}

func (any *Any) ToInt(keys ...interface{}) int {
	ret, err := getPathAsInt64(any, keys...)
	if err != nil {
		any.Error = err
		return 0;
	}
	return int(ret)
}

func (any *Any) ToUint(keys ...interface{}) uint {
	ret, err := getPathAsInt64(any, keys...)
	if err != nil {
		any.Error = err
		return 0;
	}
	return uint(ret)
}

func (any *Any) ToFloat32(keys ...interface{}) float32 {
	ret, err := getPathAsFloat64(any, keys...)
	if err != nil {
		any.Error = err
		return 0;
	}
	return float32(ret)
}

func (any *Any) ToFloat64(keys ...interface{}) float64 {
	ret, err := getPathAsFloat64(any, keys...)
	if err != nil {
		any.Error = err
		return 0;
	}
	return ret
}

func (any *Any) ToBool(keys ...interface{}) bool {
	ret, err := getPath(any.val, keys...)
	any.LastAccessed = ret
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
	ret, err := getPath(any.val, keys...)
	any.LastAccessed = ret
	if err != nil {
		any.Error = err
		return false;
	}
	return reflect.ValueOf(ret).IsNil()
}

func getPathAsInt64(any *Any, keys ...interface{}) (int64, error) {
	ret, err := getPath(any.val, keys...)
	any.LastAccessed = ret
	if err != nil {
		any.Error = err
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
	case string:
		intVal, err := strconv.ParseInt(ret, 10, 64)
		if err != nil {
			return 0, err
		}
		return intVal, nil;
	default:
		return 0, fmt.Errorf("%v is not number", ret)
	}
}

func getPathAsUint64(any *Any, keys ...interface{}) (uint64, error) {
	ret, err := getPath(any.val, keys...)
	any.LastAccessed = ret
	if err != nil {
		any.Error = err
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
	case string:
		intVal, err := strconv.ParseUint(ret, 10, 64)
		if err != nil {
			return 0, err
		}
		return intVal, nil;
	default:
		return 0, fmt.Errorf("%v is not number", ret)
	}
}

func getPathAsFloat64(any *Any, keys ...interface{}) (float64, error) {
	ret, err := getPath(any.val, keys...)
	any.LastAccessed = ret
	if err != nil {
		any.Error = err
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
	case string:
		floatVal, err := strconv.ParseFloat(ret, 64)
		if err != nil {
			return 0, err
		}
		return floatVal, nil;
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
