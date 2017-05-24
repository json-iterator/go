package jsoniter

import (
	"reflect"
	"strconv"
	"testing"
	"time"
	"unsafe"
	"github.com/json-iterator/go/require"
	"encoding/json"
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

func Test_customize_byte_array_encoder(t *testing.T) {
	should := require.New(t)
	RegisterTypeEncoder("[]uint8", func(ptr unsafe.Pointer, stream *Stream) {
		t := *((*[]byte)(ptr))
		stream.WriteString(string(t))
	})
	defer CleanEncoders()
	val := []byte("abc")
	str, err := MarshalToString(val)
	should.Nil(err)
	should.Equal(`"abc"`, str)
}

func Test_customize_float_marshal(t *testing.T) {
	should := require.New(t)
	EnableLossyFloatMarshalling()
	defer CleanEncoders()
	str, err := MarshalToString(float32(1.23456789))
	should.Nil(err)
	should.Equal("1.234568", str)
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
	should := require.New(t)
	RegisterExtension(func(type_ reflect.Type, field *reflect.StructField) ([]string, EncoderFunc, DecoderFunc) {
		if type_.String() == "jsoniter.TestObject1" && field.Name == "field1" {
			encode := func(ptr unsafe.Pointer, stream *Stream) {
				str := *((*string)(ptr))
				val, _ := strconv.Atoi(str)
				stream.WriteInt(val)
			}
			decode := func(ptr unsafe.Pointer, iter *Iterator) {
				*((*string)(ptr)) = strconv.Itoa(iter.ReadInt())
			}
			return []string{"field-1"}, encode, decode
		}
		return nil, nil, nil
	})
	obj := TestObject1{}
	err := UnmarshalFromString(`{"field-1": 100}`, &obj)
	should.Nil(err)
	should.Equal("100", obj.field1)
	str, err := MarshalToString(obj)
	should.Nil(err)
	should.Equal(`{"field-1":100}`, str)
}

func Test_unexported_fields(t *testing.T) {
	EnableUnexportedStructFieldsSupport()
	should := require.New(t)
	type TestObject struct {
		field1 string
		field2 string `json:"field-2"`
	}
	obj := TestObject{}
	obj.field1 = "hello"
	should.Nil(UnmarshalFromString(`{}`, &obj))
	should.Equal("hello", obj.field1)
	should.Nil(UnmarshalFromString(`{"field1": "world", "field-2": "abc"}`, &obj))
	should.Equal("world", obj.field1)
	should.Equal("abc", obj.field2)
	str, err := MarshalToString(obj)
	should.Nil(err)
	should.Equal(`{"field1":"world","field-2":"abc"}`, str)
}

type ObjectImplementedMarshaler int

func (obj *ObjectImplementedMarshaler) MarshalJSON() ([]byte, error) {
	return []byte(`"hello"`), nil
}

func Test_marshaler(t *testing.T) {
	type TestObject struct {
		Field *ObjectImplementedMarshaler
	}
	should := require.New(t)
	val := ObjectImplementedMarshaler(100)
	obj := TestObject{&val}
	bytes, err := json.Marshal(obj)
	should.Nil(err)
	should.Equal(`{"Field":"hello"}`, string(bytes))
	str, err := MarshalToString(obj)
	should.Nil(err)
	should.Equal(`{"Field":"hello"}`, str)
}

func Test_marshaler_and_encoder(t *testing.T) {
	type TestObject struct {
		Field *ObjectImplementedMarshaler
	}
	should := require.New(t)
	RegisterTypeEncoder("jsoniter.ObjectImplementedMarshaler", func(ptr unsafe.Pointer, stream *Stream) {
		stream.WriteString("hello from encoder")
	})
	val := ObjectImplementedMarshaler(100)
	obj := TestObject{&val}
	bytes, err := json.Marshal(obj)
	should.Nil(err)
	should.Equal(`{"Field":"hello"}`, string(bytes))
	str, err := MarshalToString(obj)
	should.Nil(err)
	should.Equal(`{"Field":"hello from encoder"}`, str)
}

type ObjectImplementedUnmarshaler int

func (obj *ObjectImplementedUnmarshaler) UnmarshalJSON([]byte) error {
	*obj = 100
	return nil
}

func Test_unmarshaler(t *testing.T) {
	type TestObject struct {
		Field *ObjectImplementedUnmarshaler
		Field2 string
	}
	should := require.New(t)
	obj := TestObject{}
	val := ObjectImplementedUnmarshaler(0)
	obj.Field = &val
	err := json.Unmarshal([]byte(`{"Field":"hello"}`), &obj)
	should.Nil(err)
	should.Equal(100, int(*obj.Field))
	err = Unmarshal([]byte(`{"Field":"hello"}`), &obj)
	should.Nil(err)
	should.Equal(100, int(*obj.Field))
}

func Test_unmarshaler_and_decoder(t *testing.T) {
	type TestObject struct {
		Field *ObjectImplementedUnmarshaler
		Field2 string
	}
	should := require.New(t)
	RegisterTypeDecoder("jsoniter.ObjectImplementedUnmarshaler", func(ptr unsafe.Pointer, iter *Iterator) {
		*(*ObjectImplementedUnmarshaler)(ptr) = 10
		iter.Skip()
	})
	obj := TestObject{}
	val := ObjectImplementedUnmarshaler(0)
	obj.Field = &val
	err := json.Unmarshal([]byte(`{"Field":"hello"}`), &obj)
	should.Nil(err)
	should.Equal(100, int(*obj.Field))
	err = Unmarshal([]byte(`{"Field":"hello"}`), &obj)
	should.Nil(err)
	should.Equal(10, int(*obj.Field))
}
