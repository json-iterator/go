package jsoniter

import (
	"unsafe"
)

type mapStringStringDecoder struct {
}

func (decoder *mapStringStringDecoder) Decode(ptr unsafe.Pointer, iter *Iterator) {
	mapPtr := (*map[string]string)(ptr)
	c := iter.nextToken()
	if c == 'n' {
		iter.skipThreeBytes('u', 'l', 'l')
		*mapPtr = nil
		return
	}
	if *mapPtr == nil {
		*mapPtr = make(map[string]string)
	}
	if c != '{' {
		iter.ReportError("ReadMapCB", `expect { or n, but found `+string([]byte{c}))
		return
	}
	c = iter.nextToken()
	if c == '}' {
		return
	}
	if c != '"' {
		iter.ReportError("ReadMapCB", `expect " after }, but found `+string([]byte{c}))
		return
	}
	iter.unreadByte()
	for c = ','; c == ','; c = iter.nextToken() {
		key := iter.ReadString()
		c = iter.nextToken()
		if c != ':' {
			iter.ReportError("ReadMapCB", "expect : after object field, but found "+string([]byte{c}))
			return
		}
		elem := iter.ReadString()
		(*mapPtr)[key] = elem
	}
	if c != '}' {
		iter.ReportError("ReadMapCB", `expect }, but found `+string([]byte{c}))
	}
}
