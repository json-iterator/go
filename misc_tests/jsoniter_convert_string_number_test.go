// +build go1.8

package misc_tests

import (
	jsoniter "github.com/json-iterator/go"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_convert_read_uint64_invalid(t *testing.T) {
	should := require.New(t)
	iter := jsoniter.ParseString(jsoniter.ConfigConvertStringNumber, ",")
	iter.ReadUint64()
	should.NotNil(iter.Error)
}

func Test_convert_read_int32_with_quote(t *testing.T) {
	should := require.New(t)
	iter := jsoniter.ParseString(jsoniter.ConfigConvertStringNumber, "\"123\"")
	v := iter.ReadUint32()
	should.Nil(iter.Error)
	should.Equal(v, uint32(123))

	iter = jsoniter.ParseString(jsoniter.ConfigConvertStringNumber, "\"123\"")
	v2 := iter.ReadInt32()
	should.Nil(iter.Error)
	should.Equal(v2, int32(123))

	iter = jsoniter.ParseString(jsoniter.ConfigConvertStringNumber, "\"-123\"")
	v2 = iter.ReadInt32()
	should.Nil(iter.Error)
	should.Equal(v2, int32(-123))
}

func Test_convert_read_int64_with_quote(t *testing.T) {
	should := require.New(t)
	iter := jsoniter.ParseString(jsoniter.ConfigConvertStringNumber, "\"123456789098765\"")
	v := iter.ReadUint64()
	should.Nil(iter.Error)
	should.Equal(v, uint64(123456789098765))

	iter = jsoniter.ParseString(jsoniter.ConfigConvertStringNumber, "\"123456789098765L\"")
	v = iter.ReadUint64()
	should.Nil(iter.Error)
	should.Equal(v, uint64(123456789098765))

	iter = jsoniter.ParseString(jsoniter.ConfigConvertStringNumber, "\"123456789098765\"")
	v2 := iter.ReadInt64()
	should.Nil(iter.Error)
	should.Equal(v2, int64(123456789098765))

	iter = jsoniter.ParseString(jsoniter.ConfigConvertStringNumber, "\"123456789098765L\"")
	v2 = iter.ReadInt64()
	should.Nil(iter.Error)
	should.Equal(v2, int64(123456789098765))

	iter = jsoniter.ParseString(jsoniter.ConfigConvertStringNumber, "\"-123456789098765\"")
	v2 = iter.ReadInt64()
	should.Nil(iter.Error)
	should.Equal(v2, int64(-123456789098765))

	iter = jsoniter.ParseString(jsoniter.ConfigConvertStringNumber, "\"-123456789098765L\"")
	v2 = iter.ReadInt64()
	should.Nil(iter.Error)
	should.Equal(v2, int64(-123456789098765))
}

func Test_read_int32_array_with_quote(t *testing.T) {
	should := require.New(t)
	input := `["123",456,"789"]`
	val := make([]int32, 0)
	err := jsoniter.ConfigConvertStringNumber.UnmarshalFromString(input, &val)
	should.Nil(err)
	should.Equal(3, len(val))
}

func Test_read_int32_negative_with_quote(t *testing.T) {
	should := require.New(t)
	input := `-123456789`
	var val int32
	err := jsoniter.ConfigConvertStringNumber.UnmarshalFromString(input, &val)
	should.Nil(err)
	should.Equal(val, int32(-123456789))

	input = `"-123456789"`
	err = jsoniter.ConfigConvertStringNumber.UnmarshalFromString(input, &val)
	should.Nil(err)
	should.Equal(val, int32(-123456789))
}

func Test_read_int64_array_with_quote(t *testing.T) {
	should := require.New(t)
	input := `["123",456L,"789L"]`
	val := make([]int64, 0)
	err := jsoniter.ConfigConvertStringNumber.UnmarshalFromString(input, &val)
	should.Nil(err)
	should.Equal(3, len(val))
}

////////////////////////////////////////////////////////////////////////////////

func Test_read_big_float_with_quote(t *testing.T) {
	should := require.New(t)
	iter := jsoniter.ParseString(jsoniter.ConfigConvertStringNumber, `12.3`)
	val := iter.ReadBigFloat()
	val64, _ := val.Float64()
	should.Equal(12.3, val64)

	iter = jsoniter.ParseString(jsoniter.ConfigConvertStringNumber, `12.3`)
	f32 := iter.ReadFloat32()
	should.Equal(float32(12.3), f32)

	iter = jsoniter.ParseString(jsoniter.ConfigConvertStringNumber, `12.3`)
	f64 := iter.ReadFloat64()
	should.Equal(12.3, f64)

	iter = jsoniter.ParseString(jsoniter.ConfigConvertStringNumber, `"12.3"`)
	val = iter.ReadBigFloat()
	val64, _ = val.Float64()
	should.Equal(12.3, val64)

	iter = jsoniter.ParseString(jsoniter.ConfigConvertStringNumber, `"-12.3"`)
	val = iter.ReadBigFloat()
	val64, _ = val.Float64()
	should.Equal(-12.3, val64)

}

func Test_read_float_with_quote(t *testing.T) {
	should := require.New(t)
	iter := jsoniter.ParseString(jsoniter.ConfigConvertStringNumber, `"12.3"`)
	f32 := iter.ReadFloat32()
	should.Equal(float32(12.3), f32)

	iter = jsoniter.ParseString(jsoniter.ConfigConvertStringNumber, `"-12.3"`)
	f32 = iter.ReadFloat32()
	should.Equal(float32(-12.3), f32)

	iter = jsoniter.ParseString(jsoniter.ConfigConvertStringNumber, `"-0.1233"`)
	f32 = iter.ReadFloat32()
	should.Equal(float32(-0.1233), f32)

	iter = jsoniter.ParseString(jsoniter.ConfigConvertStringNumber, `"124567.3"`)
	f64 := iter.ReadFloat64()
	should.Equal(124567.3, f64)

	iter = jsoniter.ParseString(jsoniter.ConfigConvertStringNumber, `"-124567.3"`)
	f64 = iter.ReadFloat64()
	should.Equal(-124567.3, f64)

	iter = jsoniter.ParseString(jsoniter.ConfigConvertStringNumber, `"-0.0124567"`)
	f64 = iter.ReadFloat64()
	should.Equal(-0.0124567, f64)
}

func Test_read_big_int_with_quote(t *testing.T) {
	should := require.New(t)
	iter := jsoniter.ParseString(jsoniter.ConfigConvertStringNumber, `92233720368547758079223372036854775807`)
	val := iter.ReadBigInt()
	should.NotNil(val)
	should.Equal(`92233720368547758079223372036854775807`, val.String())

	iter = jsoniter.ParseString(jsoniter.ConfigConvertStringNumber, `"92233720368547758079223372036854775807"`)
	val = iter.ReadBigInt()
	should.NotNil(val)
	should.Equal(`92233720368547758079223372036854775807`, val.String())
}

func Test_read_float_as_interface_with_quote(t *testing.T) {
	should := require.New(t)
	iter := jsoniter.ParseString(jsoniter.ConfigConvertStringNumber, `12.3`)
	should.Equal(float64(12.3), iter.Read())

	//iter = jsoniter.ParseString(jsoniter.ConfigConvertStringNumber, `"12.3"`)
	//should.Equal(float64(12.3), iter.Read())
}

func Test_read_float64_cursor_with_quote(t *testing.T) {
	should := require.New(t)
	iter := jsoniter.ParseString(jsoniter.ConfigConvertStringNumber, "[1.23456789\n,2,3]")
	should.True(iter.ReadArray())
	should.Equal(1.23456789, iter.Read())
	should.True(iter.ReadArray())
	should.Equal(float64(2), iter.Read())
}

func Test_read_float64_array_with_quote(t *testing.T) {
	should := require.New(t)
	input := `["123.0",4560,"789.453"]`
	val := make([]float64, 0)
	err := jsoniter.ConfigConvertStringNumber.UnmarshalFromString(input, &val)
	should.Nil(err)
	should.Equal(3, len(val))
	should.Equal([]float64{123.0, 4560, 789.453}, val)
}

/////////////////////////////////////////////////////////////////////////////////
/// number convert to string

func Test_read_string_with_quote(t *testing.T) {
	should := require.New(t)
	input := `[123.0,4560, "789.453", 12L, "34567L"]`
	val := make([]string, 0)
	err := jsoniter.ConfigConvertStringNumber.UnmarshalFromString(input, &val)
	should.Nil(err)
	should.Equal(5, len(val))
	should.Equal([]string{"123.0", "4560", "789.453", "12", "34567L"}, val)

	var s string
	err = jsoniter.ConfigConvertStringNumber.UnmarshalFromString("12", &s)
	should.Nil(err)
	should.Equal(s, "12")
}

/////////////////////////////////////////////////////////////////////////////////
// struct test

func Test_convert_string_number(t *testing.T) {
	should := require.New(t)
	s := `{"ii": 1234, "si": "1234", "ll": 908374832, "sl": "908374832", "ll2": 908374832L,
			"ll3": "908374832L", "ff": 3.1415926, "sf": "3.1415926", "bf": false, "bt": true,
			"inner": {  "iii": 4321  , "isi": "4321", "ill": -108374832, "isl": "-108374832", "ill2": -108374832L,
				"ill3": "-108374832L " , "iff": -113.1415926, "isf": "-113.1415926"},
			"arr": ["1","-2", -3,"4",-5], "farr": ["-1.2", 2.309, "-4.5", 6.325]
			}`

	t1 := struct {
		II    string  `json:"ii"`
		Si    int     `json:"si"`
		Ll    string  `json:"ll"`
		Sl    uint    `json:"sl"`
		Ll2   string  `json:"ll2"`
		Ll3   int64   `json:"ll3"`
		Ff    string  `json:"ff"`
		Sf    float64 `json:"sf"`
		Bf    bool    `json:"bf"`
		Bt    bool    `json:"bt"`
		Inner struct {
			Iii  string  `json:"iii"`
			Isi  uint32  `json:"isi"`
			Ill  string  `json:"ill"`
			Isl  int32   `json:"isl"`
			Ill2 string  `json:"ill2"`
			Ill3 int64   `json:"ill3"`
			Iff  string  `json:"iff"`
			Isf  float64 `json:"isf"`
		} `json:"inner"`
	}{}

	t2 := struct {
		II    int     `json:"ii"`
		Si    string  `json:"si"`
		Ll    int64   `json:"ll"`
		Sl    string  `json:"sl"`
		Ll2   int64   `json:"ll2"`
		Ll3   string  `json:"ll3"`
		Ff    float64 `json:"ff"`
		Sf    string  `json:"sf"`
		Bf    bool    `json:"bf"`
		Bt    bool    `json:"bt"`
		Inner struct {
			Iii  uint32  `json:"iii"`
			Isi  string  `json:"isi"`
			Ill  int32   `json:"ill"`
			Isl  string  `json:"isl"`
			Ill2 int64   `json:"ill2"`
			Ill3 string  `json:"ill3"`
			Iff  string  `json:"iff"`
			Isf  float64 `json:"isf"`
		} `json:"inner"`
	}{}

	err := jsoniter.ConfigConvertStringNumber.UnmarshalFromString(s, &t1)
	should.Nil(err)

	err = jsoniter.ConfigConvertStringNumber.UnmarshalFromString(s, &t2)
	should.Nil(err)
}
