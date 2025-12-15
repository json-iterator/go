package misc_tests

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	jsoniter "github.com/json-iterator/go"
)

func TestMaxMarshalledSize(t *testing.T) {
	testCases := []interface{}{
		nil,
		"",
		false,
		123,
		123.123,
		[]string{"foo", "bar"},
		map[string]int{"foo": 123},
	}

	for _, testCase := range testCases {
		t.Run(fmt.Sprintf("%#v", testCase), func(t *testing.T) {
			expectedBytes, err := jsoniter.Marshal(testCase)
			require.NoError(t, err)

			expectedLength := uint64(len(expectedBytes))

			expectSuccessfulMarshalling := func(t *testing.T, limit uint64) {
				cfg := jsoniter.Config{
					MaxMarshalledBytes: limit,
				}

				api := cfg.Froze()
				actualBytes, err := api.Marshal(testCase)
				require.NoError(t, err)
				require.Equal(t, expectedBytes, actualBytes)
			}

			expectFailedMarshalling := func(t *testing.T, limit uint64) {
				cfg := jsoniter.Config{
					MaxMarshalledBytes: limit,
				}

				api := cfg.Froze()
				actualBytes, err := api.Marshal(testCase)
				require.ErrorContains(t, err, fmt.Sprintf("marshalling produced a result over the configured limit of %d bytes", limit))
				require.Nil(t, actualBytes)
			}

			t.Run("limit set to 0 (unlimited)", func(t *testing.T) {
				expectSuccessfulMarshalling(t, 0)
			})

			t.Run("limit set to exact length of output", func(t *testing.T) {
				expectSuccessfulMarshalling(t, expectedLength)
			})

			t.Run("limit set to just under length of output", func(t *testing.T) {
				expectFailedMarshalling(t, expectedLength-1)
			})

			t.Run("limit set to well under length of output", func(t *testing.T) {
				expectFailedMarshalling(t, 1)
			})

			t.Run("limit set to just over length of output", func(t *testing.T) {
				expectSuccessfulMarshalling(t, expectedLength+1)
			})

			t.Run("limit set to well over length of output", func(t *testing.T) {
				expectSuccessfulMarshalling(t, expectedLength+100)
			})
		})
	}
}
