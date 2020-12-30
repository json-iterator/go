package misc_tests

import (
	"encoding/json"
	"fmt"
	"github.com/json-iterator/go"
	"testing"
)

type JsonTest struct {
	Data json.RawMessage `json:"data"`
}

func (jt *JsonTest) Bytes() []byte {
	return jt.Data
}

type JsonIterTest struct {
	Data jsoniter.RawMessage `json:"data"`
}

func (jit *JsonIterTest) Bytes() []byte {
	return jit.Data
}

func TestEncodingToAndFromJsonIter(t *testing.T) {
	encodeTypes := map[string]encoderFunc{
		"json":     encodeJson,
		"jsoniter": encodeJsonIter,
	}

	rawTypes := map[string]rawFunc{
		"json":     func() HasData { return &JsonTest{json.RawMessage("\"hello-world\"")} },
		"jsoniter": func() HasData { return &JsonIterTest{jsoniter.RawMessage("\"hello-world\"")} },
	}

	decodeTypes := map[string]decoderFunc{
		"json":     decodeJson,
		"jsoniter": decodeJsonIter,
	}

	var err error
	for encoderType, encoder := range encodeTypes {
		for encodeRawType, encodeRawFactory := range rawTypes {
			for decodeType, decoder := range decodeTypes {
				for decodeRawType, decodeRawFactory := range rawTypes {
					desc := fmt.Sprintf("encoder:%s source:%s.RawMessage decoder:%s destination:%s.RawMessage\n",
						encoderType, encodeRawType, decodeType, decodeRawType)

					if err = runTest(encoder, encodeRawFactory, decoder, decodeRawFactory); err != nil {
						t.Errorf("%s: %v", desc, err)
					}
				}
			}
		}
	}
}

type encoderFunc func(interface{}) ([]byte, error)
type decoderFunc func([]byte, interface{}) error
type rawFunc func() HasData
type HasData interface {
	Bytes() []byte
}

func encodeJson(raw interface{}) ([]byte, error) {
	return json.Marshal(raw)
}

func encodeJsonIter(raw interface{}) ([]byte, error) {
	return jsoniter.Marshal(raw)
}

func decodeJson(data []byte, raw interface{}) error {
	return json.Unmarshal(data, &raw)
}

func decodeJsonIter(data []byte, raw interface{}) error {
	return jsoniter.Unmarshal(data, &raw)
}

func runTest(encoder encoderFunc, encoderFactory rawFunc, decoder decoderFunc, decoderFactory rawFunc) error {
	raw := encoderFactory()
	data, err := encoder(raw)
	if err != nil {
		return fmt.Errorf("failed to encode: %v", err)
	}

	newRaw := decoderFactory()
	if err = decoder(data, &newRaw); err != nil {
		return fmt.Errorf("failed to decode %q: %v", string(data), err)
	}

	if !equals(raw.Bytes(), newRaw.Bytes()) {
		return fmt.Errorf("%q != %q", string(raw.Bytes()), string(newRaw.Bytes()))
	}

	return nil
}

func equals(data1 []byte, data2 []byte) bool {
	if len(data1) != len(data2) {
		return false
	}

	for idx := 0; idx < len(data1); idx++ {
		if data1[idx] != data2[idx] {
			return false
		}
	}

	return true
}
