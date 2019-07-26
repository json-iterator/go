package jsoniter_test

import (
	"strings"

	jsoniter "github.com/json-iterator/go"

	"testing"
)

type keycollector []byte

func (kc *keycollector) readKeysString(iter *jsoniter.Iterator) {
	next := iter.WhatIsNext()
	switch next {
	case jsoniter.InvalidValue:
		iter.Skip()
	case jsoniter.StringValue:
		iter.Skip()
	case jsoniter.NumberValue:
		iter.Skip()
	case jsoniter.NilValue:
		iter.Skip()
	case jsoniter.BoolValue:
		iter.Skip()
	case jsoniter.ArrayValue:
		iter.ReadArrayCB(func(iter *jsoniter.Iterator) bool {
			kc.readKeysString(iter)
			return true
		})
	case jsoniter.ObjectValue:
		iter.ReadMapCB(func(iter *jsoniter.Iterator, key string) bool {
			*kc = append(*kc, []byte(key)...)
			kc.readKeysString(iter)
			return true
		})
	default:
		iter.Skip()
	}
}

func (kc *keycollector) readKeysBytes(iter *jsoniter.Iterator, buf []byte) []byte {
	next := iter.WhatIsNext()
	switch next {
	case jsoniter.InvalidValue:
		iter.Skip()
	case jsoniter.StringValue:
		iter.Skip()
	case jsoniter.NumberValue:
		iter.Skip()
	case jsoniter.NilValue:
		iter.Skip()
	case jsoniter.BoolValue:
		iter.Skip()
	case jsoniter.ArrayValue:
		iter.ReadArrayCB(func(iter *jsoniter.Iterator) bool {
			buf = kc.readKeysBytes(iter, buf)
			return true
		})
	case jsoniter.ObjectValue:
		iter.ReadMapCBFieldAsBytes(buf, func(iter *jsoniter.Iterator, key []byte) bool {
			*kc = append(*kc, key...)
			buf = kc.readKeysBytes(iter, key)
			return true
		})
	default:
		iter.Skip()
	}
	return buf
}

func TestReadKeys(t *testing.T) {
	str := `{
    "gravatar": {
      "handle": "buger",
      "urls": [
      ],
      "avatar": "http://1.gravatar.com/avatar/f7c8edd577d13b8930d5522f28123510",
      "avatars": [
        {
          "url": "http://1.gravatar.com/avatar/f7c8edd577d13b8930d5522f28123510",
          "type": "thumbnail"
        }
      ]
    }`

	want := "gravatarhandleurlsavataravatarsurltype"

	var keysString keycollector
	keysString.readKeysString(jsoniter.Parse(jsoniter.ConfigDefault, strings.NewReader(str), 4096))
	got := string(keysString)
	if got != want {
		t.Errorf("wanted %v, got %v", want, got)
	}

	var keysBytes keycollector
	keysBytes.readKeysBytes(jsoniter.Parse(jsoniter.ConfigDefault, strings.NewReader(str), 4096), nil)
	got = string(keysBytes)
	if got != want {
		t.Errorf("wanted %v, got %v", want, got)
	}

}
