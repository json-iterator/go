package jsoniter

import (
	"bytes"
	"github.com/json-iterator/go/require"
	"testing"
)

func Test_decode_one_field_struct(t *testing.T) {
	should := require.New(t)
	type TestObject struct {
		Field1 string
	}
	obj := TestObject{}
	should.Nil(UnmarshalFromString(`{}`, &obj))
	should.Equal("", obj.Field1)
	should.Nil(UnmarshalFromString(`{"field1": "hello"}`, &obj))
	should.Equal("hello", obj.Field1)
}

func Test_decode_two_fields_struct(t *testing.T) {
	should := require.New(t)
	type TestObject struct {
		Field1 string
		Field2 string
	}
	obj := TestObject{}
	should.Nil(UnmarshalFromString(`{}`, &obj))
	should.Equal("", obj.Field1)
	should.Nil(UnmarshalFromString(`{"Field1": "a", "Field2": "b"}`, &obj))
	should.Equal("a", obj.Field1)
	should.Equal("b", obj.Field2)
}

func Test_decode_three_fields_struct(t *testing.T) {
	should := require.New(t)
	type TestObject struct {
		Field1 string
		Field2 string
		Field3 string
	}
	obj := TestObject{}
	should.Nil(UnmarshalFromString(`{}`, &obj))
	should.Equal("", obj.Field1)
	should.Nil(UnmarshalFromString(`{"Field1": "a", "Field2": "b", "Field3": "c"}`, &obj))
	should.Equal("a", obj.Field1)
	should.Equal("b", obj.Field2)
	should.Equal("c", obj.Field3)
}

func Test_decode_four_fields_struct(t *testing.T) {
	should := require.New(t)
	type TestObject struct {
		Field1 string
		Field2 string
		Field3 string
		Field4 string
	}
	obj := TestObject{}
	should.Nil(UnmarshalFromString(`{}`, &obj))
	should.Equal("", obj.Field1)
	should.Nil(UnmarshalFromString(`{"Field1": "a", "Field2": "b", "Field3": "c", "Field4": "d"}`, &obj))
	should.Equal("a", obj.Field1)
	should.Equal("b", obj.Field2)
	should.Equal("c", obj.Field3)
	should.Equal("d", obj.Field4)
}

func Test_decode_five_fields_struct(t *testing.T) {
	should := require.New(t)
	type TestObject struct {
		Field1 string
		Field2 string
		Field3 string
		Field4 string
		Field5 string
	}
	obj := TestObject{}
	should.Nil(UnmarshalFromString(`{}`, &obj))
	should.Equal("", obj.Field1)
	should.Nil(UnmarshalFromString(`{"Field1": "a", "Field2": "b", "Field3": "c", "Field4": "d", "Field5": "e"}`, &obj))
	should.Equal("a", obj.Field1)
	should.Equal("b", obj.Field2)
	should.Equal("c", obj.Field3)
	should.Equal("d", obj.Field4)
	should.Equal("e", obj.Field5)
}

func Test_decode_ten_fields_struct(t *testing.T) {
	should := require.New(t)
	type TestObject struct {
		Field1  string
		Field2  string
		Field3  string
		Field4  string
		Field5  string
		Field6  string
		Field7  string
		Field8  string
		Field9  string
		Field10 string
	}
	obj := TestObject{}
	should.Nil(UnmarshalFromString(`{}`, &obj))
	should.Equal("", obj.Field1)
	should.Nil(UnmarshalFromString(`{"Field1": "a", "Field2": "b", "Field3": "c", "Field4": "d", "Field5": "e"}`, &obj))
	should.Equal("a", obj.Field1)
	should.Equal("b", obj.Field2)
	should.Equal("c", obj.Field3)
	should.Equal("d", obj.Field4)
	should.Equal("e", obj.Field5)
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
	type AA struct {
		ID      int                    `json:"id"`
		Payload map[string]interface{} `json:"payload"`
		buf     *bytes.Buffer          `json:"-"`
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

func Test_recursive_struct(t *testing.T) {
	should := require.New(t)
	type TestObject struct {
		Field1 string
		Me     *TestObject
	}
	obj := TestObject{}
	str, err := MarshalToString(obj)
	should.Nil(err)
	should.Contains(str, `"Field1":""`)
	should.Contains(str, `"Me":null`)
	err = UnmarshalFromString(str, &obj)
	should.Nil(err)
}

func Test_one_field_struct(t *testing.T) {
	should := require.New(t)
	type YetYetAnotherObject struct {
		Field string
	}
	type YetAnotherObject struct {
		Field *YetYetAnotherObject
	}
	type AnotherObject struct {
		Field *YetAnotherObject
	}
	type TestObject struct {
		Me *AnotherObject
	}
	obj := TestObject{&AnotherObject{&YetAnotherObject{&YetYetAnotherObject{"abc"}}}}
	str, err := MarshalToString(obj)
	should.Nil(err)
	should.Equal(`{"Me":{"Field":{"Field":{"Field":"abc"}}}}`, str)
	str, err = MarshalToString(&obj)
	should.Nil(err)
	should.Equal(`{"Me":{"Field":{"Field":{"Field":"abc"}}}}`, str)
}

func Test_anonymous_struct_marshal(t *testing.T) {
	should := require.New(t)
	type TestObject struct {
		Field string
	}
	str, err := MarshalToString(struct {
		TestObject
		Field int
	}{
		Field: 100,
	})
	should.Nil(err)
	should.Equal(`{"Field":100}`, str)
}
