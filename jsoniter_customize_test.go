package jsoniter

import (
	"testing"
	"time"
	"unsafe"
	"strconv"
	"reflect"
)

func Test_customize_type_decoder(t *testing.T) {
	RegisterTypeDecoder("time.Time", func(ptr unsafe.Pointer, iter *Iterator) {
		t, err := time.ParseInLocation("2006-01-02 15:04:05", iter.ReadString(), time.UTC)
		if err != nil {
			iter.Error = err
			return
		}
		*((*time.Time)(ptr)) = t
	})
	defer ClearDecoders()
	val := time.Time{}
	err := Unmarshal([]byte(`"2016-12-05 08:43:28"`), &val)
	if err != nil {
		t.Fatal(err)
	}
	year, month, day := val.Date()
	if year != 2016 || month != 12 || day != 5 {
		t.Fatal(val)
	}
}

type Tom struct {
	field1 string
}

func Test_customize_field_decoder(t *testing.T) {
	RegisterFieldDecoder("jsoniter.Tom", "field1", func(ptr unsafe.Pointer, iter *Iterator) {
		*((*string)(ptr)) = strconv.Itoa(iter.ReadInt())
	})
	defer ClearDecoders()
	tom := Tom{}
	err := Unmarshal([]byte(`{"field1": 100}`), &tom)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_customize_field_by_extension(t *testing.T) {
	RegisterExtension(func(type_ reflect.Type, field *reflect.StructField) ([]string, DecoderFunc) {
		if (type_.String() == "jsoniter.Tom" && field.Name == "field1") {
			return []string{"field-1"}, func(ptr unsafe.Pointer, iter *Iterator) {
				*((*string)(ptr)) = strconv.Itoa(iter.ReadInt())
			}
		}
		return nil, nil
	})
	tom := Tom{}
	err := Unmarshal([]byte(`{"field-1": 100}`), &tom)
	if err != nil {
		t.Fatal(err)
	}
	if tom.field1 != "100" {
		t.Fatal(tom.field1)
	}
}

type Jerry struct {
	field1 string
}

func Test_customize_type_by_extension(t *testing.T) {
	RegisterExtension(func(type_ reflect.Type, field *reflect.StructField) ([]string, DecoderFunc) {
		if (type_.String() == "jsoniter.Jerry" && field == nil) {
			return nil, func(ptr unsafe.Pointer, iter *Iterator) {
				obj := (*Jerry)(ptr)
				obj.field1 = iter.ReadString()
			}
		}
		return nil, nil
	})
	jerry := Jerry{}
	err := Unmarshal([]byte(`"100"`), &jerry)
	if err != nil {
		t.Fatal(err)
	}
	if jerry.field1 != "100" {
		t.Fatal(jerry.field1)
	}
}