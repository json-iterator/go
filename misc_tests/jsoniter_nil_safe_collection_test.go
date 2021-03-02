package misc_tests

import (
	"testing"

	jsoniter "github.com/json-iterator/go"
	"github.com/stretchr/testify/require"
)

func Test_marshal_nil_safe_collection(t *testing.T) {
	type B struct {
		Arr  []string
		Arr2 []string
		Map  map[string]string
	}

	type A struct {
		Arr              []string
		Arr2             []string
		Map              map[string]string
		Struct           B
		StructPoint      *B
		StructSlice      []B
		StructSlicePoint []*B
		StructMap        map[string]B
		StructMapPoint   map[string]*B
	}

	a := A{}
	a.Arr2 = append(a.Arr2, "aa")
	should := require.New(t)
	out := `{"Arr":[],"Arr2":["aa"],"Map":{},"Struct":{"Arr":[],"Arr2":[],"Map":{}},"StructPoint":null,"StructSlice":[],"StructSlicePoint":[],"StructMap":{},"StructMapPoint":{}}`
	aout, err := jsoniter.Config{NilSafeCollection: true, SortMapKeys: true}.Froze().Marshal(&a)
	should.Nil(err)
	should.Equal(out, string(aout))
	bout, err := jsoniter.Config{NilSafeCollection: true, SortMapKeys: false}.Froze().Marshal(&a)
	should.Nil(err)
	should.Equal(out, string(bout))
}
