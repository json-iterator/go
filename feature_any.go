package jsoniter

import (
	"fmt"
	"reflect"
)

type Any interface {
	LastError() error
	ValueType() ValueType
	ToBool() bool
	ToInt() int
	ToInt32() int32
	ToInt64() int64
	ToUint() uint
	ToUint32() uint32
	ToUint64() uint64
	ToFloat32() float32
	ToFloat64() float64
	ToString() string
	Get(path ...interface{}) Any
	Size() int
	Keys() []string
	IterateObject() (func() (string, Any, bool), bool)
	IterateArray() (func() (Any, bool), bool)
	GetArray() []Any
	SetArray(newList []Any) bool
	GetObject() map[string]Any
	SetObject(map[string]Any) bool
	GetInterface() interface{}
	WriteTo(stream *Stream)
}

type baseAny struct{}

func (any *baseAny) Get(path ...interface{}) Any {
	return &invalidAny{baseAny{}, fmt.Errorf("Get %v from simple value", path)}
}

func (any *baseAny) Size() int {
	return 0
}

func (any *baseAny) Keys() []string {
	return []string{}
}

func (any *baseAny) IterateObject() (func() (string, Any, bool), bool) {
	return nil, false
}

func (any *baseAny) IterateArray() (func() (Any, bool), bool) {
	return nil, false
}

func (any *baseAny) GetArray() []Any {
	return []Any{}
}

func (any *baseAny) SetArray(newList []Any) bool {
	return false
}

func (any *baseAny) GetObject() map[string]Any {
	return map[string]Any{}
}

func (any *baseAny) SetObject(map[string]Any) bool {
	return false
}

func WrapInt32(val int32) Any {
	return &int32Any{baseAny{}, val}
}

func WrapInt64(val int64) Any {
	return &int64Any{baseAny{}, val}
}

func WrapUint32(val uint32) Any {
	return &uint32Any{baseAny{}, val}
}

func WrapUint64(val uint64) Any {
	return &uint64Any{baseAny{}, val}
}

func WrapFloat64(val float64) Any {
	return &floatAny{baseAny{}, val}
}

func WrapString(val string) Any {
	return &stringAny{baseAny{}, nil, val}
}

func Wrap(val interface{}) Any {
	if val == nil {
		return &nilAny{}
	}
	asAny, isAny := val.(Any)
	if isAny {
		return asAny
	}
	type_ := reflect.TypeOf(val)
	switch type_.Kind() {
	case reflect.Slice:
		return wrapArray(val)
	case reflect.Struct:
		return wrapStruct(val)
	case reflect.Map:
		return wrapMap(val)
	case reflect.String:
		return WrapString(val.(string))
	case reflect.Int:
		return WrapInt64(int64(val.(int)))
	case reflect.Int8:
		return WrapInt32(int32(val.(int8)))
	case reflect.Int16:
		return WrapInt32(int32(val.(int16)))
	case reflect.Int32:
		return WrapInt32(val.(int32))
	case reflect.Int64:
		return WrapInt64(val.(int64))
	case reflect.Uint:
		return WrapUint64(uint64(val.(uint)))
	case reflect.Uint8:
		return WrapUint32(uint32(val.(uint8)))
	case reflect.Uint16:
		return WrapUint32(uint32(val.(uint16)))
	case reflect.Uint32:
		return WrapUint32(uint32(val.(uint32)))
	case reflect.Uint64:
		return WrapUint64(val.(uint64))
	case reflect.Float32:
		return WrapFloat64(float64(val.(float32)))
	case reflect.Float64:
		return WrapFloat64(val.(float64))
	case reflect.Bool:
		if val.(bool) == true {
			return &trueAny{}
		} else {
			return &falseAny{}
		}
	}
	return &invalidAny{baseAny{}, fmt.Errorf("unsupported type: %v", type_)}
}

func (iter *Iterator) ReadAny() Any {
	return iter.readAny()
}

func (iter *Iterator) readAny() Any {
	c := iter.nextToken()
	switch c {
	case '"':
		return iter.readStringAny()
	case 'n':
		iter.skipFixedBytes(3) // null
		return &nilAny{}
	case 't':
		iter.skipFixedBytes(3) // true
		return &trueAny{}
	case 'f':
		iter.skipFixedBytes(4) // false
		return &falseAny{}
	case '{':
		return iter.readObjectAny()
	case '[':
		return iter.readArrayAny()
	case '-':
		return iter.readNumberAny(false)
	default:
		return iter.readNumberAny(true)
	}
}

func (iter *Iterator) readNumberAny(positive bool) Any {
	iter.startCapture(iter.head - 1)
	dotFound := iter.skipNumberAndTellDotFoundOrNot()
	lazyBuf := iter.stopCapture()
	if dotFound {
		return &float64LazyAny{baseAny{}, iter.cfg, lazyBuf, nil, 0}
	} else {
		if positive {
			return &uint64LazyAny{baseAny{}, iter.cfg, lazyBuf, nil, 0}
		} else {
			return &int64LazyAny{baseAny{}, iter.cfg, lazyBuf, nil, 0}
		}
	}
}

func (iter *Iterator) readStringAny() Any {
	iter.startCapture(iter.head - 1)
	iter.skipString()
	lazyBuf := iter.stopCapture()
	return &stringLazyAny{baseAny{}, iter.cfg, lazyBuf, nil, ""}
}

func (iter *Iterator) readObjectAny() Any {
	iter.startCapture(iter.head - 1)
	iter.skipObject()
	lazyBuf := iter.stopCapture()
	return &objectLazyAny{baseAny{}, iter.cfg, lazyBuf, nil, nil, lazyBuf}
}

func (iter *Iterator) readArrayAny() Any {
	iter.startCapture(iter.head - 1)
	iter.skipArray()
	lazyBuf := iter.stopCapture()
	return &arrayLazyAny{baseAny{}, iter.cfg, lazyBuf, nil, nil, lazyBuf}
}
