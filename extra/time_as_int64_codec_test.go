package extra

import (
	"testing"
	"time"
	"github.com/json-iterator/go/require"
	"github.com/json-iterator/go"
)

func Test_time_as_int64(t *testing.T) {
	should := require.New(t)
	RegisterTimeAsInt64Codec(time.Nanosecond)
	output, err := jsoniter.Marshal(time.Unix(1, 1002))
	should.Nil(err)
	should.Equal("1000001002", string(output))
}

func Test_time_as_int64_keep_microsecond(t *testing.T) {
	should := require.New(t)
	RegisterTimeAsInt64Codec(time.Microsecond)
	output, err := jsoniter.Marshal(time.Unix(1, 1002))
	should.Nil(err)
	should.Equal("1000001", string(output))
}
