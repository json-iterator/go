package jsoniter

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_empty_object(t *testing.T) {
	should := require.New(t)
	iter := ParseString(ConfigDefault, `{}`)
	field := iter.ReadObject()
	should.Equal("", field)
	iter = ParseString(ConfigDefault, `{}`)
	iter.ReadObjectCB(func(iter *Iterator, field string) bool {
		should.FailNow("should not call")
		return true
	})
}

func Test_one_field(t *testing.T) {
	should := require.New(t)
	iter := ParseString(ConfigDefault, `{"a": "stream"}`)
	field := iter.ReadObject()
	should.Equal("a", field)
	value := iter.ReadString()
	should.Equal("stream", value)
	field = iter.ReadObject()
	should.Equal("", field)
	iter = ParseString(ConfigDefault, `{"a": "stream"}`)
	should.True(iter.ReadObjectCB(func(iter *Iterator, field string) bool {
		should.Equal("a", field)
		return true
	}))

}

func Test_two_field(t *testing.T) {
	should := require.New(t)
	iter := ParseString(ConfigDefault, `{ "a": "stream" , "c": "d" }`)
	field := iter.ReadObject()
	should.Equal("a", field)
	value := iter.ReadString()
	should.Equal("stream", value)
	field = iter.ReadObject()
	should.Equal("c", field)
	value = iter.ReadString()
	should.Equal("d", value)
	field = iter.ReadObject()
	should.Equal("", field)
	iter = ParseString(ConfigDefault, `{"field1": "1", "field2": 2}`)
	for field := iter.ReadObject(); field != ""; field = iter.ReadObject() {
		switch field {
		case "field1":
			iter.ReadString()
		case "field2":
			iter.ReadInt64()
		default:
			iter.ReportError("bind object", "unexpected field")
		}
	}
}

func Test_object_wrapper_any_get_all(t *testing.T) {
	should := require.New(t)
	type TestObject struct {
		Field1 []int
		Field2 []int
	}
	any := Wrap(TestObject{[]int{1, 2}, []int{3, 4}})
	should.Contains(any.Get('*', 0).ToString(), `"Field2":3`)
	should.Contains(any.Keys(), "Field1")
	should.Contains(any.Keys(), "Field2")
	should.NotContains(any.Keys(), "Field3")

	//should.Contains(any.GetObject()["Field1"].GetArray()[0], 1)
}

func Test_write_object(t *testing.T) {
	should := require.New(t)
	buf := &bytes.Buffer{}
	stream := NewStream(Config{IndentionStep: 2}.Froze(), buf, 4096)
	stream.WriteObjectStart()
	stream.WriteObjectField("hello")
	stream.WriteInt(1)
	stream.WriteMore()
	stream.WriteObjectField("world")
	stream.WriteInt(2)
	stream.WriteObjectEnd()
	stream.Flush()
	should.Nil(stream.Error)
	should.Equal("{\n  \"hello\": 1,\n  \"world\": 2\n}", buf.String())
}

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
	should.Nil(UnmarshalFromString(`{"Field1": "a", "Field2": "stream"}`, &obj))
	should.Equal("a", obj.Field1)
	should.Equal("stream", obj.Field2)
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
	should.Nil(UnmarshalFromString(`{"Field1": "a", "Field2": "stream", "Field3": "c"}`, &obj))
	should.Equal("a", obj.Field1)
	should.Equal("stream", obj.Field2)
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
	should.Nil(UnmarshalFromString(`{"Field1": "a", "Field2": "stream", "Field3": "c", "Field4": "d"}`, &obj))
	should.Equal("a", obj.Field1)
	should.Equal("stream", obj.Field2)
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
	should.Nil(UnmarshalFromString(`{"Field1": "a", "Field2": "stream", "Field3": "c", "Field4": "d", "Field5": "e"}`, &obj))
	should.Equal("a", obj.Field1)
	should.Equal("stream", obj.Field2)
	should.Equal("c", obj.Field3)
	should.Equal("d", obj.Field4)
	should.Equal("e", obj.Field5)
}

func Test_decode_six_fields_struct(t *testing.T) {
	should := require.New(t)
	type TestObject struct {
		Field1 string
		Field2 string
		Field3 string
		Field4 string
		Field5 string
		Field6 string
	}
	obj := TestObject{}
	should.Nil(UnmarshalFromString(`{}`, &obj))
	should.Equal("", obj.Field1)
	should.Nil(UnmarshalFromString(`{"Field1": "a", "Field2": "stream", "Field3": "c", "Field4": "d", "Field5": "e", "Field6": "x"}`, &obj))
	should.Equal("a", obj.Field1)
	should.Equal("stream", obj.Field2)
	should.Equal("c", obj.Field3)
	should.Equal("d", obj.Field4)
	should.Equal("e", obj.Field5)
	should.Equal("x", obj.Field6)
}

func Test_decode_seven_fields_struct(t *testing.T) {
	should := require.New(t)
	type TestObject struct {
		Field1 string
		Field2 string
		Field3 string
		Field4 string
		Field5 string
		Field6 string
		Field7 string
	}
	obj := TestObject{}
	should.Nil(UnmarshalFromString(`{}`, &obj))
	should.Equal("", obj.Field1)
	should.Nil(UnmarshalFromString(`{"Field1": "a", "Field2": "stream", "Field3": "c", "Field4": "d", "Field5": "e", "Field6": "x", "Field7":"y"}`, &obj))
	should.Equal("a", obj.Field1)
	should.Equal("stream", obj.Field2)
	should.Equal("c", obj.Field3)
	should.Equal("d", obj.Field4)
	should.Equal("e", obj.Field5)
	should.Equal("x", obj.Field6)
	should.Equal("y", obj.Field7)
}

func Test_decode_eight_fields_struct(t *testing.T) {
	should := require.New(t)
	type TestObject struct {
		Field1 string
		Field2 string
		Field3 string
		Field4 string
		Field5 string
		Field6 string
		Field7 string
		Field8 string
	}
	obj := TestObject{}
	should.Nil(UnmarshalFromString(`{}`, &obj))
	should.Equal("", obj.Field1)
	should.Nil(UnmarshalFromString(`{"Field8":"1", "Field1": "a", "Field2": "stream", "Field3": "c", "Field4": "d", "Field5": "e", "Field6": "x", "Field7":"y"}`, &obj))
	should.Equal("a", obj.Field1)
	should.Equal("stream", obj.Field2)
	should.Equal("c", obj.Field3)
	should.Equal("d", obj.Field4)
	should.Equal("e", obj.Field5)
	should.Equal("x", obj.Field6)
	should.Equal("y", obj.Field7)
	should.Equal("1", obj.Field8)
}

func Test_decode_nine_fields_struct(t *testing.T) {
	should := require.New(t)
	type TestObject struct {
		Field1 string
		Field2 string
		Field3 string
		Field4 string
		Field5 string
		Field6 string
		Field7 string
		Field8 string
		Field9 string
	}
	obj := TestObject{}
	should.Nil(UnmarshalFromString(`{}`, &obj))
	should.Equal("", obj.Field1)
	should.Nil(UnmarshalFromString(`{"Field8" : "zzzzzzzzzzz", "Field7": "zz", "Field6" : "xx", "Field1": "a", "Field2": "stream", "Field3": "c", "Field4": "d", "Field5": "e", "Field9":"f"}`, &obj))
	should.Equal("a", obj.Field1)
	should.Equal("stream", obj.Field2)
	should.Equal("c", obj.Field3)
	should.Equal("d", obj.Field4)
	should.Equal("e", obj.Field5)
	should.Equal("xx", obj.Field6)
	should.Equal("zz", obj.Field7)
	should.Equal("zzzzzzzzzzz", obj.Field8)
	should.Equal("f", obj.Field9)
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
	should.Nil(UnmarshalFromString(`{"Field10":"x", "Field9": "x", "Field8":"x", "Field7":"x", "Field6":"x", "Field1": "a", "Field2": "stream", "Field3": "c", "Field4": "d", "Field5": "e"}`, &obj))
	should.Equal("a", obj.Field1)
	should.Equal("stream", obj.Field2)
	should.Equal("c", obj.Field3)
	should.Equal("d", obj.Field4)
	should.Equal("e", obj.Field5)
	should.Equal("x", obj.Field6)
	should.Equal("x", obj.Field7)
	should.Equal("x", obj.Field8)
	should.Equal("x", obj.Field9)
	should.Equal("x", obj.Field10)
}

func Test_decode_more_than_ten_fields_struct(t *testing.T) {
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
		Field11 int
	}
	obj := TestObject{}
	should.Nil(UnmarshalFromString(`{}`, &obj))
	should.Equal("", obj.Field1)
	should.Nil(UnmarshalFromString(`{"Field11":1, "Field1": "a", "Field2": "stream", "Field3": "c", "Field4": "d", "Field5": "e"}`, &obj))
	should.Equal("a", obj.Field1)
	should.Equal("stream", obj.Field2)
	should.Equal("c", obj.Field3)
	should.Equal("d", obj.Field4)
	should.Equal("e", obj.Field5)
	should.Equal(1, obj.Field11)
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

func Test_decode_struct_field_with_tag_string(t *testing.T) {
	should := require.New(t)
	type TestObject struct {
		Field1 int `json:",string"`
	}
	obj := TestObject{Field1: 100}
	should.Nil(UnmarshalFromString(`{"Field1": "100"}`, &obj))
	should.Equal(100, obj.Field1)
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
		buf     *bytes.Buffer
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

func Test_ignore_field_on_not_valid_type(t *testing.T) {
	should := require.New(t)
	type TestObject struct {
		Field1 string `json:"field-1,omitempty"`
		Field2 func() `json:"-"`
	}
	obj := TestObject{}
	obj.Field1 = "hello world"
	obj.Field2 = func() {}
	str, err := MarshalToString(&obj)
	should.Nil(err)
	should.Equal(`{"field-1":"hello world"}`, str)
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

func Test_encode_anonymous_struct(t *testing.T) {
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

func Test_decode_anonymous_struct(t *testing.T) {
	should := require.New(t)
	type Inner struct {
		Key string `json:"key"`
	}

	type Outer struct {
		Inner
	}
	var outer Outer
	j := []byte("{\"key\":\"value\"}")
	should.Nil(Unmarshal(j, &outer))
	should.Equal("value", outer.Key)
}

func Test_multiple_level_anonymous_struct(t *testing.T) {
	type Level1 struct {
		Field1 string
	}
	type Level2 struct {
		Level1
		Field2 string
	}
	type Level3 struct {
		Level2
		Field3 string
	}
	should := require.New(t)
	obj := Level3{Level2{Level1{"1"}, "2"}, "3"}
	output, err := MarshalToString(obj)
	should.Nil(err)
	should.Equal(`{"Field1":"1","Field2":"2","Field3":"3"}`, output)
}

func Test_multiple_level_anonymous_struct_with_ptr(t *testing.T) {
	type Level1 struct {
		Field1 string
		Field2 string
		Field4 string
	}
	type Level2 struct {
		*Level1
		Field2 string
		Field3 string
	}
	type Level3 struct {
		*Level2
		Field3 string
	}
	should := require.New(t)
	obj := Level3{&Level2{&Level1{"1", "", "4"}, "2", ""}, "3"}
	output, err := MarshalToString(obj)
	should.Nil(err)
	should.Contains(output, `"Field1":"1"`)
	should.Contains(output, `"Field2":"2"`)
	should.Contains(output, `"Field3":"3"`)
	should.Contains(output, `"Field4":"4"`)
}

func Test_shadow_struct_field(t *testing.T) {
	should := require.New(t)
	type omit *struct{}
	type CacheItem struct {
		Key    string `json:"key"`
		MaxAge int    `json:"cacheAge"`
	}
	output, err := MarshalToString(struct {
		*CacheItem

		// Omit bad keys
		OmitMaxAge omit `json:"cacheAge,omitempty"`

		// Add nice keys
		MaxAge int `json:"max_age"`
	}{
		CacheItem: &CacheItem{
			Key:    "value",
			MaxAge: 100,
		},
		MaxAge: 20,
	})
	should.Nil(err)
	should.Contains(output, `"key":"value"`)
	should.Contains(output, `"max_age":20`)
}

func Test_embedded_order(t *testing.T) {
	type A struct {
		Field2 string
	}

	type C struct {
		Field5 string
	}

	type B struct {
		Field4 string
		C
		Field6 string
	}

	type TestObject struct {
		Field1 string
		A
		Field3 string
		B
		Field7 string
	}
	should := require.New(t)
	s := TestObject{}
	output, err := MarshalToString(s)
	should.Nil(err)
	should.Equal(`{"Field1":"","Field2":"","Field3":"","Field4":"","Field5":"","Field6":"","Field7":""}`, output)
}

func Test_decode_nested(t *testing.T) {
	type StructOfString struct {
		Field1 string
		Field2 string
	}
	iter := ParseString(ConfigDefault, `[{"field1": "hello"}, null, {"field2": "world"}]`)
	slice := []*StructOfString{}
	iter.ReadVal(&slice)
	if len(slice) != 3 {
		fmt.Println(iter.Error)
		t.Fatal(len(slice))
	}
	if slice[0].Field1 != "hello" {
		fmt.Println(iter.Error)
		t.Fatal(slice[0])
	}
	if slice[1] != nil {
		fmt.Println(iter.Error)
		t.Fatal(slice[1])
	}
	if slice[2].Field2 != "world" {
		fmt.Println(iter.Error)
		t.Fatal(slice[2])
	}
}
