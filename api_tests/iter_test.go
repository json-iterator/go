package test

import (
	"io"
	"strings"
	"testing"

	jsoniter "github.com/json-iterator/go"
	"github.com/stretchr/testify/require"
)

func Test_iterator_offsets(t *testing.T) {
	json := `{ "foo": "bar", "num": 123 }
	{ "num" : 27 }
	{ "arr": [1.1,2.2,3.3], "obj": { "key": "val"}, "num": 68 }
	{ "foo.name": "quiz", "num": "321"}`

	should := require.New(t)
	startOffsets := []int64{}
	for i, r := range json {
		if r == '{' {
			startOffsets = append(startOffsets, int64(i))
		}
	}
	should.Len(startOffsets, 5)

	iter := jsoniter.Parse(jsoniter.ConfigDefault, strings.NewReader(json), 8)
	should.NotNil(iter)

	should.EqualValues(0, iter.InputOffset())
	should.Equal(startOffsets[0], iter.InputOffset())

	iter.ReadObjectCB(func(iter *jsoniter.Iterator, key string) bool {
		switch key {
		case "foo":
			should.Equal(jsoniter.StringValue, iter.WhatIsNext())
			should.EqualValues(9, iter.InputOffset())
		case "num":
			should.Equal(jsoniter.NumberValue, iter.WhatIsNext())
			should.EqualValues(23, iter.InputOffset())
		default:
			should.NotNil(nil, "unexpected key: %s", key)
		}

		// skip the value
		iter.Skip()

		return true
	})
	should.NoError(iter.Error)
	should.EqualValues(28, iter.InputOffset())
	should.EqualValues('}', json[iter.InputOffset()-1])
	// there's still some whitespace to get to the next object
	should.NotEqual(startOffsets[1], iter.InputOffset())

	// read second line

	should.Equal(jsoniter.ObjectValue, iter.WhatIsNext())
	should.Equal(startOffsets[1], iter.InputOffset())
	iter.ReadObjectCB(func(iter *jsoniter.Iterator, key string) bool {
		switch key {
		case "num":
			should.Equal(jsoniter.NumberValue, iter.WhatIsNext())
			should.EqualValues(40, iter.InputOffset())
		default:
			should.NotNil(nil, "unexpected key: %s", key)
		}

		// skip the value
		iter.Skip()

		return true
	})
	should.NoError(iter.Error)

	// read third line
	should.Equal(jsoniter.ObjectValue, iter.WhatIsNext())
	should.Equal(startOffsets[2], iter.InputOffset())
	for iter.ReadObject() != "" {
		iter.Skip()
	}
	should.NoError(iter.Error)

	// read fourth line
	should.Equal(jsoniter.ObjectValue, iter.WhatIsNext())
	should.Equal(startOffsets[4], iter.InputOffset())
	for iter.ReadObject() != "" {
		iter.Skip()
	}
	should.NoError(iter.Error)

	should.Equal(jsoniter.InvalidValue, iter.WhatIsNext())
	should.EqualValues(len(json), iter.InputOffset())
	should.EqualError(iter.Error, io.EOF.Error())
}
