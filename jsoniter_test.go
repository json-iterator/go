package jsoniter

import (
	"testing"
)

func Test_int_1(t *testing.T) {
	iter := NewIterator([]byte("1"))
	val := iter.ReadInt()
	if val != 1 {
		t.Fatal(val)
	}
}

func Test_int_minus_1(t *testing.T) {
	iter := NewIterator([]byte("-1"))
	val := iter.ReadInt()
	if val != -1 {
		t.Fatal(val)
	}
}

func Test_int_100(t *testing.T) {
	iter := NewIterator([]byte("100,"))
	val := iter.ReadInt()
	if val != 100 {
		t.Fatal(val)
	}
}

func Test_int_0(t *testing.T) {
	iter := NewIterator([]byte("0"))
	val := iter.ReadInt()
	if val != 0 {
		t.Fatal(val)
	}
}

//func Test_single_element(t *testing.T) {
//	iter := NewIterator([]byte("[1]"))
//	val := 0
//	iter.ReadArray(func(iter Iterator, index int) {
//		val = iter.ReadInt()
//	})
//	if val != 1 {
//		t.Fatal(val)
//	}
//}
//
//func Test_multiple_elements(t *testing.T) {
//	iter := NewIterator([]byte("[1, 2]"))
//	result := []int{0, 0}
//	iter.ReadArray(func(iter Iterator, index int) {
//		result[index] = iter.ReadInt()
//	})
//	if !reflect.DeepEqual([]int{1, 2}, result) {
//		t.Fatal(result)
//	}
//}
//
//func Test_invalid_array(t *testing.T) {
//	iter := NewIterator([]byte("[1, ]"))
//	result := []int{0, 0}
//	var foundErr error
//	iter.ErrorHandler = func(err error) {
//		foundErr = err
//	}
//	iter.ReadArray(func(iter Iterator, index int) {
//		result[index] = iter.ReadInt()
//	})
//	if foundErr == nil {
//		t.FailNow()
//	}
//}
//
//func Test_single_field(t *testing.T) {
//	iter := NewIterator([]byte(`{"a": 1}`))
//	result := map[string]int{}
//	iter.ReadObject(func(iter Iterator, field string) {
//		result[field] = iter.ReadInt()
//	})
//	if !reflect.DeepEqual(map[string]int{"a": 1}, result) {
//		t.Fatal(result)
//	}
//}
//
//func Test_multiple_fields(t *testing.T) {
//	iter := NewIterator([]byte(`{"a": 1, "b": 2}`))
//	result := map[string]int{}
//	iter.ReadObject(func(iter Iterator, field string) {
//		result[field] = iter.ReadInt()
//	})
//	if !reflect.DeepEqual(map[string]int{"a": 1, "b": 2}, result) {
//		t.Fatal(result)
//	}
//}
//
//func Test_nested_object(t *testing.T) {
//	iter := NewIterator([]byte(`{"a": [{"b": 2}, {"b": 1}]}`))
//	obj := map[string][]map[string]int{}
//	iter.ReadObject(func(iter Iterator, field string) {
//		array := []map[string]int{}
//		iter.ReadArray(func(iter Iterator, index int) {
//			nestedObj := map[string]int{}
//			iter.ReadObject(func(iter Iterator, field string) {
//				nestedObj[field] = iter.ReadInt()
//			})
//			array = append(array, nestedObj)
//		})
//		obj[field] = array
//	})
//	if !reflect.DeepEqual(obj, map[string][]map[string]int{
//		"a": {{"b": 2}, {"b": 1}},
//	}) {
//		t.Fatal(obj)
//	}
//}
//
//func Test_skip(t *testing.T) {
//	iter := NewIterator([]byte(`{"a": [{"b": 2}, {"b": 1}], "c": 3}`))
//	val := 0
//	iter.ReadObject(func(iter Iterator, field string) {
//		if ("c" == field) {
//			val = iter.ReadInt()
//		} else {
//			iter.Skip()
//		}
//	})
//	if val != 3 {
//		t.Fatal(val)
//	}
//}
//
//func Test_int8(t *testing.T) {
//	iter := NewIterator([]byte("[1]"))
//	val := int8(0)
//	iter.ReadArray(func(iter Iterator, index int) {
//		val = iter.ReadInt8()
//	})
//	if val != int8(1) {
//		t.Fatal(val)
//	}
//}
//
//func Test_int16(t *testing.T) {
//	iter := NewIterator([]byte("[1]"))
//	val := int16(0)
//	iter.ReadArray(func(iter Iterator, index int) {
//		val = iter.ReadInt16()
//	})
//	if val != int16(1) {
//		t.Fatal(val)
//	}
//}
//
//func Test_int32(t *testing.T) {
//	iter := NewIterator([]byte("[1]"))
//	val := int32(0)
//	iter.ReadArray(func(iter Iterator, index int) {
//		val = iter.ReadInt32()
//	})
//	if val != int32(1) {
//		t.Fatal(val)
//	}
//}
//
//func Test_int64(t *testing.T) {
//	iter := NewIterator([]byte("[1]"))
//	val := int64(0)
//	iter.ReadArray(func(iter Iterator, index int) {
//		val = iter.ReadInt64()
//	})
//	if val != int64(1) {
//		t.Fatal(val)
//	}
//}
//
//func Test_uint(t *testing.T) {
//	iter := NewIterator([]byte("[1]"))
//	val := uint(0)
//	iter.ReadArray(func(iter Iterator, index int) {
//		val = iter.ReadUint()
//	})
//	if val != uint(1) {
//		t.Fatal(val)
//	}
//}
//
//func Test_uint8(t *testing.T) {
//	iter := NewIterator([]byte("[1]"))
//	val := uint8(0)
//	iter.ReadArray(func(iter Iterator, index int) {
//		val = iter.ReadUint8()
//	})
//	if val != uint8(1) {
//		t.Fatal(val)
//	}
//}
//
//func Test_uint16(t *testing.T) {
//	iter := NewIterator([]byte("[1]"))
//	val := uint16(0)
//	iter.ReadArray(func(iter Iterator, index int) {
//		val = iter.ReadUint16()
//	})
//	if val != uint16(1) {
//		t.Fatal(val)
//	}
//}
//
//func Test_uint32(t *testing.T) {
//	iter := NewIterator([]byte("[1]"))
//	val := uint32(0)
//	iter.ReadArray(func(iter Iterator, index int) {
//		val = iter.ReadUint32()
//	})
//	if val != uint32(1) {
//		t.Fatal(val)
//	}
//}
//
//func Test_uint64(t *testing.T) {
//	iter := NewIterator([]byte("[1]"))
//	val := uint64(0)
//	iter.ReadArray(func(iter Iterator, index int) {
//		val = iter.ReadUint64()
//	})
//	if val != uint64(1) {
//		t.Fatal(val)
//	}
//}
//
//func Test_float32(t *testing.T) {
//	iter := NewIterator([]byte("[1.1]"))
//	val := float32(0)
//	iter.ReadArray(func(iter Iterator, index int) {
//		val = iter.ReadFloat32()
//	})
//	if val != float32(1.1) {
//		t.Fatal(val)
//	}
//}
//
//func Test_float64(t *testing.T) {
//	iter := NewIterator([]byte("[1.1]"))
//	val := float64(0)
//	iter.ReadArray(func(iter Iterator, index int) {
//		val = iter.ReadFloat64()
//	})
//	if val != float64(1.1) {
//		t.Fatal(val)
//	}
//}
//