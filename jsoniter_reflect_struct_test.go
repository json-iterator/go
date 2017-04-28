package jsoniter

import (
	"testing"
	"github.com/json-iterator/go/require"
	"bytes"
)

func Test_decode_one_field_struct(t *testing.T) {
	should := require.New(t)
	type TestObject struct {
		field1 string
	}
	obj := TestObject{}
	should.Nil(UnmarshalFromString(`{}`, &obj))
	should.Equal("", obj.field1)
	should.Nil(UnmarshalFromString(`{"field1": "hello"}`, &obj))
	should.Equal("hello", obj.field1)
}

func Test_decode_two_fields_struct(t *testing.T) {
	should := require.New(t)
	type TestObject struct {
		field1 string
		field2 string
	}
	obj := TestObject{}
	should.Nil(UnmarshalFromString(`{}`, &obj))
	should.Equal("", obj.field1)
	should.Nil(UnmarshalFromString(`{"field1": "a", "field2": "b"}`, &obj))
	should.Equal("a", obj.field1)
	should.Equal("b", obj.field2)
}

func Test_decode_three_fields_struct(t *testing.T) {
	should := require.New(t)
	type TestObject struct {
		field1 string
		field2 string
		field3 string
	}
	obj := TestObject{}
	should.Nil(UnmarshalFromString(`{}`, &obj))
	should.Equal("", obj.field1)
	should.Nil(UnmarshalFromString(`{"field1": "a", "field2": "b", "field3": "c"}`, &obj))
	should.Equal("a", obj.field1)
	should.Equal("b", obj.field2)
	should.Equal("c", obj.field3)
}

func Test_decode_four_fields_struct(t *testing.T) {
	should := require.New(t)
	type TestObject struct {
		field1 string
		field2 string
		field3 string
		field4 string
	}
	obj := TestObject{}
	should.Nil(UnmarshalFromString(`{}`, &obj))
	should.Equal("", obj.field1)
	should.Nil(UnmarshalFromString(`{"field1": "a", "field2": "b", "field3": "c", "field4": "d"}`, &obj))
	should.Equal("a", obj.field1)
	should.Equal("b", obj.field2)
	should.Equal("c", obj.field3)
	should.Equal("d", obj.field4)
}

func Test_decode_five_fields_struct(t *testing.T) {
	should := require.New(t)
	type TestObject struct {
		field1 string
		field2 string
		field3 string
		field4 string
		field5 string
	}
	obj := TestObject{}
	should.Nil(UnmarshalFromString(`{}`, &obj))
	should.Equal("", obj.field1)
	should.Nil(UnmarshalFromString(`{"field1": "a", "field2": "b", "field3": "c", "field4": "d", "field5": "e"}`, &obj))
	should.Equal("a", obj.field1)
	should.Equal("b", obj.field2)
	should.Equal("c", obj.field3)
	should.Equal("d", obj.field4)
	should.Equal("e", obj.field5)
}

func Test_decode_ten_fields_struct(t *testing.T) {
	should := require.New(t)
	type TestObject struct {
		field1  string
		field2  string
		field3  string
		field4  string
		field5  string
		field6  string
		field7  string
		field8  string
		field9  string
		field10 string
	}
	obj := TestObject{}
	should.Nil(UnmarshalFromString(`{}`, &obj))
	should.Equal("", obj.field1)
	should.Nil(UnmarshalFromString(`{"field1": "a", "field2": "b", "field3": "c", "field4": "d", "field5": "e"}`, &obj))
	should.Equal("a", obj.field1)
	should.Equal("b", obj.field2)
	should.Equal("c", obj.field3)
	should.Equal("d", obj.field4)
	should.Equal("e", obj.field5)
}

func Test_decode_struct_field_with_tag(t *testing.T) {
	should := require.New(t)
	type TestObject struct {
		Field1 string `json:"field-1"`
		Field2 string `json:"-"`
		Field3 int    `json:",string"`
	}
	obj := TestObject{Field2: "world"}
	UnmarshalFromString(`{"field-1": "hello", "field2": "", "Field3": "100"}`, &obj)
	should.Equal("hello", obj.Field1)
	should.Equal("world", obj.Field2)
	should.Equal(100, obj.Field3)
}

func Test_write_val_zero_field_struct(t *testing.T) {
	should := require.New(t)
	type TestObject struct {
	}
	obj := TestObject{}
	str, err := MarshalToString(obj)
	should.Nil(err)
	should.Equal(`{}`, str)
}

func Test_write_val_one_field_struct(t *testing.T) {
	should := require.New(t)
	type TestObject struct {
		Field1 string `json:"field-1"`
	}
	obj := TestObject{"hello"}
	str, err := MarshalToString(obj)
	should.Nil(err)
	should.Equal(`{"field-1":"hello"}`, str)
}

func Test_mixed(t *testing.T) {
	should := require.New(t)
	type  AA struct {
		ID      int `json:"id"`
		Payload map[string]interface{} `json:"payload"`
		buf     *bytes.Buffer `json:"-"`
	}
	aa := AA{}
	err := UnmarshalFromString(` {"id":1, "payload":{"account":"123","password":"456"}}`, &aa)
	should.Nil(err)
	should.Equal(1, aa.ID)
	should.Equal("123", aa.Payload["account"])
}

func Test_omit_empty(t *testing.T) {
	should := require.New(t)
	type TestObject struct {
		Field1 string `json:"field-1,omitempty"`
		Field2 string `json:"field-2,omitempty"`
		Field3 string `json:"field-3,omitempty"`
	}
	obj := TestObject{}
	obj.Field2 = "hello"
	str, err := MarshalToString(&obj)
	should.Nil(err)
	should.Equal(`{"field-2":"hello"}`, str)
}

func Test_any_within_struct(t *testing.T) {
	should := require.New(t)
	type TestObject struct {
		Field1 Any
		Field2 Any
	}
	obj := TestObject{}
	err := UnmarshalFromString(`{"Field1": "hello", "Field2": [1,2,3]}`, &obj)
	should.Nil(err)
	should.Equal("hello", obj.Field1.ToString())
	should.Equal("[1,2,3]", obj.Field2.ToString())
}
