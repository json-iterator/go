package jsoniter

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_writeByte_should_grow_buffer(t *testing.T) {
	should := require.New(t)
	stream := NewStream(ConfigDefault, nil, 1)
	stream.writeByte('1')
	should.Equal("1", string(stream.Buffer()))
	should.Equal(1, len(stream.buf))
	stream.writeByte('2')
	should.Equal("12", string(stream.Buffer()))
	should.Equal(2, len(stream.buf))
	stream.writeThreeBytes('3', '4', '5')
	should.Equal("12345", string(stream.Buffer()))
}

func Test_writeBytes_should_grow_buffer(t *testing.T) {
	should := require.New(t)
	stream := NewStream(ConfigDefault, nil, 1)
	stream.Write([]byte{'1', '2'})
	should.Equal("12", string(stream.Buffer()))
	should.Equal(2, len(stream.buf))
	stream.Write([]byte{'3', '4', '5', '6', '7'})
	should.Equal("1234567", string(stream.Buffer()))
	should.Equal(7, len(stream.buf))
}

func Test_writeIndention_should_grow_buffer(t *testing.T) {
	should := require.New(t)
	stream := NewStream(Config{IndentionStep: 2}.Froze(), nil, 1)
	stream.WriteVal([]int{1, 2, 3})
	should.Equal("[\n  1,\n  2,\n  3\n]", string(stream.Buffer()))
}

func Test_writeRaw_should_grow_buffer(t *testing.T) {
	should := require.New(t)
	stream := NewStream(ConfigDefault, nil, 1)
	stream.WriteRaw("123")
	should.Nil(stream.Error)
	should.Equal("123", string(stream.Buffer()))
}

func Test_writeString_should_grow_buffer(t *testing.T) {
	should := require.New(t)
	stream := NewStream(ConfigDefault, nil, 0)
	stream.WriteString("123")
	should.Nil(stream.Error)
	should.Equal(`"123"`, string(stream.Buffer()))
}

type NopWriter struct {
	bufferSize int
}

func (w *NopWriter) Write(p []byte) (n int, err error) {
	w.bufferSize = cap(p)
	return len(p), nil
}

func Test_flush_buffer_should_stop_grow_buffer(t *testing.T) {
	writer := new(NopWriter)
	NewEncoder(writer).Encode(make([]int, 10000000))
	should := require.New(t)
	should.Equal(8, writer.bufferSize)
}

func TestStreamIndentionStep(t *testing.T) {
	should := require.New(t)

	config := (Config{
		EscapeHTML:    true,
		IndentionStep: 4,
	}).Froze()

	// {
	// 	"_source": {
	// 		"@timestamp": "2018-10-07T15:01:48.442Z",
	// 		"source": {
	// 			"geoip": {
	// 				"city_name": "Mosman",
	// 				"country_name": "Australia"
	// 			}
	// 		},
	// 		"user": {
	// 			"service": "wOot",
	// 			"uid": "1"
	// 		}
	// 	}
	// }

	res := `{
    "_source": {
        "@timestamp": "2018-10-07T15:01:48.442Z",
        "source": {
            "geoip": {
                "city_name": "Mosman",
                "country_name": "Australia"
            }
        },
        "user": {
            "service": "wOot",
            "uid": "1"
        }
    }
}`

	stream := NewStream(config, nil, 0)

	stream.WriteObjectStart()

	stream.WriteObjectField("_source")
	stream.WriteObjectStart()

	stream.WriteObjectField("@timestamp")
	stream.WriteString("2018-10-07T15:01:48.442Z")
	stream.WriteMore()
	// ,
	stream.WriteObjectField("source")
	stream.WriteObjectStart()
	stream.WriteObjectField("geoip")
	stream.WriteObjectStart()
	stream.WriteObjectField("city_name")
	stream.WriteString("Mosman")
	stream.WriteMore()
	stream.WriteObjectField("country_name")
	stream.WriteString("Australia")
	stream.WriteObjectEnd()
	stream.WriteObjectEnd()
	stream.WriteMore()
	// ,
	stream.WriteObjectField("user")
	stream.WriteObjectStart()
	stream.WriteObjectField("service")
	stream.WriteString("wOot")
	stream.WriteMore()
	stream.WriteObjectField("uid")
	stream.WriteString("1")
	stream.WriteObjectEnd()

	stream.WriteObjectEnd()
	stream.WriteObjectEnd()

	should.Equal(res, string(stream.Buffer()))
}
