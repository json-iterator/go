package jsoniter

import (
	"testing"

	"github.com/json-iterator/go/require"
)

var input2 = `{"a":["b",{"c":"d","f":[1,2,3,4]},"e",["2",3,{"g":"h"}]]}`

func TestNode(t *testing.T) {
	cReader := &KeyString{Key: STRING("c")}
	eReader := &KeyString{Key: INT(2)}
	nReader := &KeyInt64{Key: INT(3)}
	gReader := &KeyString{Key: STRING("g")}

	iter := ParseString(input2)
	many := NewExtractMany(iter,
		NewExtract(NewPath(ArrayNode("a"), ArrayIndex(1)), cReader),
		NewExtract(NewPath(ArrayNode("a")), eReader),
		NewExtract(NewPath(ArrayNode("a"), ArrayIndex(1), ArrayNode("f")), nReader),
		NewExtract(NewPath(ArrayNode("a"), ArrayIndex(3), ArrayIndex(2)), gReader),
	)

	err := many.ExtractObject()
	require.NoError(t, iter.Error)
	require.NoError(t, err)
	require.Equal(t, eReader.String, "e")
	require.Equal(t, cReader.String, "d")
	require.Equal(t, nReader.Int64, int64(4))
	require.Equal(t, gReader.String, "h")
}
