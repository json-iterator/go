package jsoniter

import (
	"testing"

	"github.com/json-iterator/go/require"
)

var input2 = `{"a":["b",{"c":"d","f":[1,2,3,4]},"e",["2",3,{"g":"h"}]],"k":2,"l":{"m":"n"},"i":"j"}`

func TestNode(t *testing.T) {
	cReader := &KeyString{Key: Field("c")}
	eReader := &KeyString{Key: Index(2)}
	nReader := &KeyInt64{Key: Index(3)}
	gReader := &KeyString{Key: Field("g")}

	iter := ParseString(input2)
	many := NewExtractMany(iter,
		NewExtract(NewPath(Field("a"), Index(1)), cReader),
		NewExtract(NewPath(Field("a")), eReader),
		NewExtract(NewPath(Field("a"), Index(1), Field("f")), nReader),
		NewExtract(NewPath(Field("a"), Index(3), Index(2)), gReader),
	)

	err := many.ExtractObject()
	require.NoError(t, iter.Error)
	require.NoError(t, err)
	require.Equal(t, eReader.String, "e")
	require.Equal(t, cReader.String, "d")
	require.Equal(t, nReader.Int64, int64(4))
	require.Equal(t, gReader.String, "h")
}

func BenchmarkExtract(b *testing.B) {
	input := []byte(`
{
    "_shards":{
        "total" : 5,
        "successful" : 5,
        "failed" : 0
    },
    "hits":{
        "total" : 1,
        "hits" : [
            {
                "_index" : "twitter",
                "_type" : "tweet",
                "_id" : "1",
                "_source" : {
                    "user" : "kimchy",
                    "postDate" : "2009-11-15T14:12:12",
                    "message" : "trying out Elasticsearch"
                }
            }
        ]
    },
    "code": 200
}`)

	cReader := &KeyString{Key: Field("code")}

	for i := 0; i < b.N; i++ {
		iter := ParseBytes(input)
		many := NewExtractMany(iter,
			NewExtract(nil, cReader),
		)

		many.ExtractObject()
	}
}

var input = `{"type":"xxx", "payload":{"a":true,"b":3,"c":{"d":"e"}},"f":"g"}`

func TestExtractFirstLayer(t *testing.T) {
	stringReader := &KeyString{Key: Field("type")}
	iter := ParseString(input)
	many := NewExtractMany(iter, NewExtract(nil, stringReader))
	err := many.ExtractObject()
	require.NoError(t, iter.Error)
	require.NoError(t, err)
	require.Equal(t, stringReader.String, "xxx")

	stringReader = &KeyString{Key: Field("f")}
	iter = ParseString(input)
	many = NewExtractMany(iter, NewExtract(nil, stringReader))
	err = many.ExtractObject()
	require.NoError(t, iter.Error)
	require.NoError(t, err)
	require.Equal(t, stringReader.String, "g")
}

func TestExtractSecondLayer(t *testing.T) {
	boolReader := &KeyBool{Key: Field("a")}
	iter := ParseString(input)
	many := NewExtractMany(iter, NewExtract(NewPath(Field("payload")), boolReader))
	err := many.ExtractObject()
	require.NoError(t, iter.Error)
	require.NoError(t, err)
	require.True(t, boolReader.Bool)
}

func TestExtractThirdLayer(t *testing.T) {
	strReader := &KeyString{Key: Field("d")}
	iter := ParseString(input)
	many := NewExtractMany(iter, NewExtract(NewPath(Field("payload"), Field("c")), strReader))
	err := many.ExtractObject()
	require.NoError(t, iter.Error)
	require.NoError(t, err)
	require.Equal(t, strReader.String, "e")
}

func TestExtractMany(t *testing.T) {
	typeReader, aReader, bReader, dReader, fReader :=
		&KeyString{Key: Field("type")}, &KeyBool{Key: Field("a")}, &KeyInt64{Key: Field("b")}, &KeyString{Key: Field("d")}, &KeyString{Key: Field("f")}
	iter := ParseString(input)
	extracts := []*Extract{
		NewExtract(nil, typeReader, fReader),
		NewExtract(NewPath(Field("payload")), aReader, bReader),
		NewExtract(NewPath(Field("payload"), Field("c")), dReader),
	}
	many := NewExtractMany(iter, extracts...)
	err := many.ExtractObject()
	require.NoError(t, iter.Error)
	require.NoError(t, err)

	require.Equal(t, typeReader.String, "xxx")
	require.True(t, aReader.Bool)
	require.Equal(t, bReader.Int64, int64(3))
	require.Equal(t, dReader.String, "e")
	require.Equal(t, fReader.String, "g")
}

func TestExtractManyWithMayNull(t *testing.T) {
	typeReader, fakeReader := &KeyString{Key: Field("type")}, &KeyString{Key: Field("fake"), MayNull: true}

	iter := ParseString(input)
	many := NewExtractMany(iter, NewExtract(nil, typeReader, fakeReader))
	err := many.ExtractObject()

	require.NoError(t, iter.Error)
	require.NoError(t, err)
	require.Equal(t, typeReader.String, "xxx")
	require.False(t, fakeReader.HasRead())
}

func TestExtractWithInvalidPath(t *testing.T) {
	typeReader := &KeyString{Key: Field("type")}
	iter := ParseString(input)
	many := NewExtractMany(iter, NewExtract(NewPath(Field("type")), typeReader))
	err := many.ExtractObject()
	require.EqualError(t, err, ErrInvalidPath.Error())
}

func TestExtractWithInvalidReader(t *testing.T) {
	typeReader := &KeyBool{Key: Field("type")}
	iter := ParseString(input)
	many := NewExtractMany(iter, NewExtract(nil, typeReader))
	err := many.ExtractObject()
	require.EqualError(t, err, ErrInvalidKeyReader.Error())
}

func TestExtractMustRead(t *testing.T) {
	typeReader, mustReader := &KeyString{Key: Field("type")}, &KeyString{Key: Field("must")}
	iter := ParseString(input)
	many := NewExtractMany(iter, NewExtract(nil, typeReader, mustReader))
	err := many.ExtractObject()
	require.EqualError(t, err, "1 not read")
}
