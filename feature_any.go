package jsoniter

import (
	"fmt"
	"io"
	"reflect"
)

type Any interface {
	LastError() error
	ValueType() ValueType
	MustBeValid() Any
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
	ToVal(val interface{})
	Get(path ...interface{}) Any
	// TODO: add Set
	Size() int
	Keys() []string
	// TODO: remove me
	GetArray() []Any
	// TODO: remove me
	GetObject() map[string]Any
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

func (any *baseAny) GetArray() []Any {
	return []Any{}
}

func (any *baseAny) GetObject() map[string]Any {
	return map[string]Any{}
}

func (any *baseAny) ToVal(obj interface{}) {
	panic("not implemented")
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
	return &stringAny{baseAny{}, val}
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
		iter.unreadByte()
		return &stringAny{baseAny{}, iter.ReadString()}
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
	iter.skipNumber()
	lazyBuf := iter.stopCapture()
	return &numberLazyAny{baseAny{}, iter.cfg, lazyBuf, nil}
}

func (iter *Iterator) readObjectAny() Any {
	iter.startCapture(iter.head - 1)
	iter.skipObject()
	lazyBuf := iter.stopCapture()
	return &objectLazyAny{baseAny{}, iter.cfg, lazyBuf, nil}
}

func (iter *Iterator) readArrayAny() Any {
	iter.startCapture(iter.head - 1)
	iter.skipArray()
	lazyBuf := iter.stopCapture()
	return &arrayLazyAny{baseAny{}, iter.cfg, lazyBuf, nil}
}

func locateObjectField(iter *Iterator, target string) []byte {
	var found []byte
	iter.ReadObjectCB(func(iter *Iterator, field string) bool {
		if field == target {
			found = iter.SkipAndReturnBytes()
			return false
		}
		iter.Skip()
		return true
	})
	return found
}

func locateArrayElement(iter *Iterator, target int) []byte {
	var found []byte
	n := 0
	iter.ReadArrayCB(func(iter *Iterator) bool {
		if n == target {
			found = iter.SkipAndReturnBytes()
			return false
		}
		iter.Skip()
		n++
		return true
	})
	return found
}

func locatePath(iter *Iterator, path []interface{}) Any {
	for i, pathKeyObj := range path {
		switch pathKey := pathKeyObj.(type) {
		case string:
			valueBytes := locateObjectField(iter, pathKey)
			if valueBytes == nil {
				return newInvalidAny(path[i:])
			}
			iter.ResetBytes(valueBytes)
		case int:
			valueBytes := locateArrayElement(iter, pathKey)
			if valueBytes == nil {
				return newInvalidAny(path[i:])
			}
			iter.ResetBytes(valueBytes)
		case int32:
			if '*' == pathKey {
				return iter.readAny().Get(path[i:]...)
			} else {
				return newInvalidAny(path[i:])
			}
		default:
			return newInvalidAny(path[i:])
		}
	}
	if iter.Error != nil && iter.Error != io.EOF {
		return &invalidAny{baseAny{}, iter.Error}
	}
	return iter.readAny()
}
