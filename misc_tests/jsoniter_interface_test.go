package misc_tests

import (
	"encoding/json"
	"fmt"
	"testing"
	"github.com/stretchr/testify/require"
	"github.com/json-iterator/go"
)

type MyInterface interface {
	Hello() string
}

type MyString string

func (ms MyString) Hello() string {
	return string(ms)
}

func Test_decode_object_contain_non_empty_interface(t *testing.T) {
	type TestObject struct {
		Field MyInterface
	}
	should := require.New(t)
	obj := TestObject{}
	obj.Field = MyString("abc")
	should.Nil(jsoniter.UnmarshalFromString(`{"Field": "hello"}`, &obj))
	should.Equal(MyString("hello"), obj.Field)
}

func Test_nil_non_empty_interface(t *testing.T) {
	type TestObject struct {
		Field []MyInterface
	}
	should := require.New(t)
	obj := TestObject{}
	b := []byte(`{"Field":["AAA"]}`)
	should.NotNil(json.Unmarshal(b, &obj))
	should.NotNil(jsoniter.Unmarshal(b, &obj))
}

func Test_read_large_number_as_interface(t *testing.T) {
	should := require.New(t)
	var val interface{}
	err := jsoniter.Config{UseNumber: true}.Froze().UnmarshalFromString(`123456789123456789123456789`, &val)
	should.Nil(err)
	output, err := jsoniter.MarshalToString(val)
	should.Nil(err)
	should.Equal(`123456789123456789123456789`, output)
}

func Test_unmarshal_ptr_to_interface(t *testing.T) {
	type TestData struct {
		Name string `json:"name"`
	}
	should := require.New(t)
	var obj interface{} = &TestData{}
	err := json.Unmarshal([]byte(`{"name":"value"}`), &obj)
	should.Nil(err)
	should.Equal("&{value}", fmt.Sprintf("%v", obj))
	obj = interface{}(&TestData{})
	err = jsoniter.Unmarshal([]byte(`{"name":"value"}`), &obj)
	should.Nil(err)
	should.Equal("&{value}", fmt.Sprintf("%v", obj))
}

func Test_nil_out_null_interface(t *testing.T) {
	type TestData struct {
		Field interface{} `json:"field"`
	}
	should := require.New(t)

	var boolVar bool
	obj := TestData{
		Field: &boolVar,
	}

	data1 := []byte(`{"field": true}`)

	err := jsoniter.Unmarshal(data1, &obj)
	should.NoError(err)
	should.Equal(true, *(obj.Field.(*bool)))

	data2 := []byte(`{"field": null}`)

	err = jsoniter.Unmarshal(data2, &obj)
	should.NoError(err)
	should.Equal(nil, obj.Field)

	// Checking stdlib behavior matches.
	obj2 := TestData{
		Field: &boolVar,
	}

	err = json.Unmarshal(data1, &obj2)
	should.NoError(err)
	should.Equal(true, *(obj2.Field.(*bool)))

	err = json.Unmarshal(data2, &obj2)
	should.NoError(err)
	should.Equal(nil, obj2.Field)
}

func Test_overwrite_interface_ptr_value_with_nil(t *testing.T) {
	type Wrapper struct {
		Payload interface{} `json:"payload,omitempty"`
	}
	type Payload struct {
		Value int `json:"val,omitempty"`
	}

	should := require.New(t)

	payload := &Payload{}
	wrapper := &Wrapper{
		Payload: &payload,
	}

	err := json.Unmarshal([]byte(`{"payload": {"val": 42}}`), &wrapper)
	should.Equal(nil, err)
	should.Equal(&payload, wrapper.Payload)
	should.Equal(42, (*(wrapper.Payload.(**Payload))).Value)

	err = json.Unmarshal([]byte(`{"payload": null}`), &wrapper)
	should.Equal(nil, err)
	should.Equal(&payload, wrapper.Payload)
	should.Equal((*Payload)(nil), payload)

	payload = &Payload{}
	wrapper = &Wrapper{
		Payload: &payload,
	}

	err = jsoniter.Unmarshal([]byte(`{"payload": {"val": 42}}`), &wrapper)
	should.Equal(nil, err)
	should.Equal(&payload, wrapper.Payload)
	should.Equal(42, (*(wrapper.Payload.(**Payload))).Value)

	err = jsoniter.Unmarshal([]byte(`{"payload": null}`), &wrapper)
	should.Equal(nil, err)
	should.Equal(&payload, wrapper.Payload)
	should.Equal((*Payload)(nil), payload)
}

func Test_overwrite_interface_value_with_nil(t *testing.T) {
	type Wrapper struct {
		Payload interface{} `json:"payload,omitempty"`
	}
	type Payload struct {
		Value int `json:"val,omitempty"`
	}

	should := require.New(t)

	payload := &Payload{}
	wrapper := &Wrapper{
		Payload: payload,
	}

	err := json.Unmarshal([]byte(`{"payload": {"val": 42}}`), &wrapper)
	should.Equal(nil, err)
	should.Equal(42, (*(wrapper.Payload.(*Payload))).Value)

	err = json.Unmarshal([]byte(`{"payload": null}`), &wrapper)
	should.Equal(nil, err)
	should.Equal(nil, wrapper.Payload)
	should.Equal(42, payload.Value)

	payload = &Payload{}
	wrapper = &Wrapper{
		Payload: payload,
	}

	err = jsoniter.Unmarshal([]byte(`{"payload": {"val": 42}}`), &wrapper)
	should.Equal(nil, err)
	should.Equal(42, (*(wrapper.Payload.(*Payload))).Value)

	err = jsoniter.Unmarshal([]byte(`{"payload": null}`), &wrapper)
	should.Equal(nil, err)
	should.Equal(nil, wrapper.Payload)
	should.Equal(42, payload.Value)
}

func Test_unmarshal_into_nil(t *testing.T) {
	type Payload struct {
		Value int `json:"val,omitempty"`
	}
	type Wrapper struct {
		Payload interface{} `json:"payload,omitempty"`
	}

	should := require.New(t)

	var payload *Payload
	wrapper := &Wrapper{
		Payload: payload,
	}

	err := json.Unmarshal([]byte(`{"payload": {"val": 42}}`), &wrapper)
	should.Nil(err)
	should.NotNil(wrapper.Payload)
	should.Nil(payload)

	err = json.Unmarshal([]byte(`{"payload": null}`), &wrapper)
	should.Nil(err)
	should.Nil(wrapper.Payload)
	should.Nil(payload)

	payload = nil
	wrapper = &Wrapper{
		Payload: payload,
	}

	err = jsoniter.Unmarshal([]byte(`{"payload": {"val": 42}}`), &wrapper)
	should.Nil(err)
	should.NotNil(wrapper.Payload)
	should.Nil(payload)

	err = jsoniter.Unmarshal([]byte(`{"payload": null}`), &wrapper)
	should.Nil(err)
	should.Nil(wrapper.Payload)
	should.Nil(payload)
}
