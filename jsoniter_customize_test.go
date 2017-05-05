package jsoniter

import (
	"reflect"
	"strconv"
	"testing"
	"time"
	"unsafe"
	"github.com/json-iterator/go/require"
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
	defer CleanDecoders()
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

func Test_customize_type_encoder(t *testing.T) {
	should := require.New(t)
	RegisterTypeEncoder("time.Time", func(ptr unsafe.Pointer, stream *Stream) {
		t := *((*time.Time)(ptr))
		stream.WriteString(t.UTC().Format("2006-01-02 15:04:05"))
	})
	defer CleanEncoders()
	val := time.Unix(0, 0)
	str, err := MarshalToString(val)
	should.Nil(err)
	should.Equal(`"1970-01-01 00:00:00"`, str)
}

type Tom struct {
	field1 string
}

func Test_customize_field_decoder(t *testing.T) {
	RegisterFieldDecoder("jsoniter.Tom", "field1", func(ptr unsafe.Pointer, iter *Iterator) {
		*((*string)(ptr)) = strconv.Itoa(iter.ReadInt())
	})
	defer CleanDecoders()
	tom := Tom{}
	err := Unmarshal([]byte(`{"field1": 100}`), &tom)
	if err != nil {
		t.Fatal(err)
	}
}

type TestObject1 struct {
	field1 string
}

func Test_customize_field_by_extension(t *testing.T) {
	RegisterExtension(func(type_ reflect.Type, field *reflect.StructField) ([]string, DecoderFunc) {
		if type_.String() == "jsoniter.TestObject1" && field.Name == "field1" {
			return []string{"field-1"}, func(ptr unsafe.Pointer, iter *Iterator) {
				*((*string)(ptr)) = strconv.Itoa(iter.ReadInt())
			}
		}
		return nil, nil
	})
	obj := TestObject1{}
	err := Unmarshal([]byte(`{"field-1": 100}`), &obj)
	if err != nil {
		t.Fatal(err)
	}
	if obj.field1 != "100" {
		t.Fatal(obj.field1)
	}
}
