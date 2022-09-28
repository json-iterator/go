package jsoniter

import (
	"fmt"
	"io"
	"strings"
	"unsafe"

	"github.com/modern-go/reflect2"
)

func decoderOfStruct(ctx *ctx, typ reflect2.Type) ValDecoder {
	bindings := map[string]*Binding{}
	structDescriptor := describeStruct(ctx, typ)
	for _, binding := range structDescriptor.Fields {
		for _, fromName := range binding.FromNames {
			old := bindings[fromName]
			if old == nil {
				bindings[fromName] = binding
				continue
			}
			ignoreOld, ignoreNew := resolveConflictBinding(ctx.frozenConfig, old, binding)
			if ignoreOld {
				delete(bindings, fromName)
			}
			if !ignoreNew {
				bindings[fromName] = binding
			}
		}
	}
	fields := map[string]*StructFieldDecoder{}
	for k, binding := range bindings {
		fields[k] = binding.Decoder.(*StructFieldDecoder)
	}

	if !ctx.caseSensitive() {
		for k, binding := range bindings {
			if _, found := fields[strings.ToLower(k)]; !found {
				fields[strings.ToLower(k)] = binding.Decoder.(*StructFieldDecoder)
			}
		}
	}

	return createStructDecoder(ctx, typ, fields)
}

func createStructDecoder(ctx *ctx, typ reflect2.Type, fields map[string]*StructFieldDecoder) ValDecoder {
	if ctx.disallowUnknownFields {
		return &generalStructDecoder{typ: typ, fields: fields, disallowUnknownFields: true}
	}
	knownHash := map[int64]struct{}{
		0: {},
	}

	switch len(fields) {
	case 0:
		return &skipObjectDecoder{typ}
	case 1:
		for fieldName, fieldDecoder := range fields {
			fieldHash := calcHash(fieldName, ctx.caseSensitive())
			_, known := knownHash[fieldHash]
			if known {
				return &generalStructDecoder{typ, fields, false}
			}
			knownHash[fieldHash] = struct{}{}
			return &oneFieldStructDecoder{typ, fieldHash, fieldDecoder}
		}
	case 2:
		var fieldHash1 int64
		var fieldHash2 int64
		var fieldDecoder1 *StructFieldDecoder
		var fieldDecoder2 *StructFieldDecoder
		for fieldName, fieldDecoder := range fields {
			fieldHash := calcHash(fieldName, ctx.caseSensitive())
			_, known := knownHash[fieldHash]
			if known {
				return &generalStructDecoder{typ, fields, false}
			}
			knownHash[fieldHash] = struct{}{}
			if fieldHash1 == 0 {
				fieldHash1 = fieldHash
				fieldDecoder1 = fieldDecoder
			} else {
				fieldHash2 = fieldHash
				fieldDecoder2 = fieldDecoder
			}
		}
		return &twoFieldsStructDecoder{typ, fieldHash1, fieldDecoder1, fieldHash2, fieldDecoder2}
	case 3:
		var fieldName1 int64
		var fieldName2 int64
		var fieldName3 int64
		var fieldDecoder1 *StructFieldDecoder
		var fieldDecoder2 *StructFieldDecoder
		var fieldDecoder3 *StructFieldDecoder
		for fieldName, fieldDecoder := range fields {
			fieldHash := calcHash(fieldName, ctx.caseSensitive())
			_, known := knownHash[fieldHash]
			if known {
				return &generalStructDecoder{typ, fields, false}
			}
			knownHash[fieldHash] = struct{}{}
			if fieldName1 == 0 {
				fieldName1 = fieldHash
				fieldDecoder1 = fieldDecoder
			} else if fieldName2 == 0 {
				fieldName2 = fieldHash
				fieldDecoder2 = fieldDecoder
			} else {
				fieldName3 = fieldHash
				fieldDecoder3 = fieldDecoder
			}
		}
		return &threeFieldsStructDecoder{typ,
			fieldName1, fieldDecoder1,
			fieldName2, fieldDecoder2,
			fieldName3, fieldDecoder3}
	case 4:
		var fieldName1 int64
		var fieldName2 int64
		var fieldName3 int64
		var fieldName4 int64
		var fieldDecoder1 *StructFieldDecoder
		var fieldDecoder2 *StructFieldDecoder
		var fieldDecoder3 *StructFieldDecoder
		var fieldDecoder4 *StructFieldDecoder
		for fieldName, fieldDecoder := range fields {
			fieldHash := calcHash(fieldName, ctx.caseSensitive())
			_, known := knownHash[fieldHash]
			if known {
				return &generalStructDecoder{typ, fields, false}
			}
			knownHash[fieldHash] = struct{}{}
			if fieldName1 == 0 {
				fieldName1 = fieldHash
				fieldDecoder1 = fieldDecoder
			} else if fieldName2 == 0 {
				fieldName2 = fieldHash
				fieldDecoder2 = fieldDecoder
			} else if fieldName3 == 0 {
				fieldName3 = fieldHash
				fieldDecoder3 = fieldDecoder
			} else {
				fieldName4 = fieldHash
				fieldDecoder4 = fieldDecoder
			}
		}
		return &fourFieldsStructDecoder{typ,
			fieldName1, fieldDecoder1,
			fieldName2, fieldDecoder2,
			fieldName3, fieldDecoder3,
			fieldName4, fieldDecoder4}
	case 5:
		var fieldName1 int64
		var fieldName2 int64
		var fieldName3 int64
		var fieldName4 int64
		var fieldName5 int64
		var fieldDecoder1 *StructFieldDecoder
		var fieldDecoder2 *StructFieldDecoder
		var fieldDecoder3 *StructFieldDecoder
		var fieldDecoder4 *StructFieldDecoder
		var fieldDecoder5 *StructFieldDecoder
		for fieldName, fieldDecoder := range fields {
			fieldHash := calcHash(fieldName, ctx.caseSensitive())
			_, known := knownHash[fieldHash]
			if known {
				return &generalStructDecoder{typ, fields, false}
			}
			knownHash[fieldHash] = struct{}{}
			if fieldName1 == 0 {
				fieldName1 = fieldHash
				fieldDecoder1 = fieldDecoder
			} else if fieldName2 == 0 {
				fieldName2 = fieldHash
				fieldDecoder2 = fieldDecoder
			} else if fieldName3 == 0 {
				fieldName3 = fieldHash
				fieldDecoder3 = fieldDecoder
			} else if fieldName4 == 0 {
				fieldName4 = fieldHash
				fieldDecoder4 = fieldDecoder
			} else {
				fieldName5 = fieldHash
				fieldDecoder5 = fieldDecoder
			}
		}
		return &fiveFieldsStructDecoder{typ,
			fieldName1, fieldDecoder1,
			fieldName2, fieldDecoder2,
			fieldName3, fieldDecoder3,
			fieldName4, fieldDecoder4,
			fieldName5, fieldDecoder5}
	case 6:
		var fieldName1 int64
		var fieldName2 int64
		var fieldName3 int64
		var fieldName4 int64
		var fieldName5 int64
		var fieldName6 int64
		var fieldDecoder1 *StructFieldDecoder
		var fieldDecoder2 *StructFieldDecoder
		var fieldDecoder3 *StructFieldDecoder
		var fieldDecoder4 *StructFieldDecoder
		var fieldDecoder5 *StructFieldDecoder
		var fieldDecoder6 *StructFieldDecoder
		for fieldName, fieldDecoder := range fields {
			fieldHash := calcHash(fieldName, ctx.caseSensitive())
			_, known := knownHash[fieldHash]
			if known {
				return &generalStructDecoder{typ, fields, false}
			}
			knownHash[fieldHash] = struct{}{}
			if fieldName1 == 0 {
				fieldName1 = fieldHash
				fieldDecoder1 = fieldDecoder
			} else if fieldName2 == 0 {
				fieldName2 = fieldHash
				fieldDecoder2 = fieldDecoder
			} else if fieldName3 == 0 {
				fieldName3 = fieldHash
				fieldDecoder3 = fieldDecoder
			} else if fieldName4 == 0 {
				fieldName4 = fieldHash
				fieldDecoder4 = fieldDecoder
			} else if fieldName5 == 0 {
				fieldName5 = fieldHash
				fieldDecoder5 = fieldDecoder
			} else {
				fieldName6 = fieldHash
				fieldDecoder6 = fieldDecoder
			}
		}
		return &sixFieldsStructDecoder{typ,
			fieldName1, fieldDecoder1,
			fieldName2, fieldDecoder2,
			fieldName3, fieldDecoder3,
			fieldName4, fieldDecoder4,
			fieldName5, fieldDecoder5,
			fieldName6, fieldDecoder6}
	case 7:
		var fieldName1 int64
		var fieldName2 int64
		var fieldName3 int64
		var fieldName4 int64
		var fieldName5 int64
		var fieldName6 int64
		var fieldName7 int64
		var fieldDecoder1 *StructFieldDecoder
		var fieldDecoder2 *StructFieldDecoder
		var fieldDecoder3 *StructFieldDecoder
		var fieldDecoder4 *StructFieldDecoder
		var fieldDecoder5 *StructFieldDecoder
		var fieldDecoder6 *StructFieldDecoder
		var fieldDecoder7 *StructFieldDecoder
		for fieldName, fieldDecoder := range fields {
			fieldHash := calcHash(fieldName, ctx.caseSensitive())
			_, known := knownHash[fieldHash]
			if known {
				return &generalStructDecoder{typ, fields, false}
			}
			knownHash[fieldHash] = struct{}{}
			if fieldName1 == 0 {
				fieldName1 = fieldHash
				fieldDecoder1 = fieldDecoder
			} else if fieldName2 == 0 {
				fieldName2 = fieldHash
				fieldDecoder2 = fieldDecoder
			} else if fieldName3 == 0 {
				fieldName3 = fieldHash
				fieldDecoder3 = fieldDecoder
			} else if fieldName4 == 0 {
				fieldName4 = fieldHash
				fieldDecoder4 = fieldDecoder
			} else if fieldName5 == 0 {
				fieldName5 = fieldHash
				fieldDecoder5 = fieldDecoder
			} else if fieldName6 == 0 {
				fieldName6 = fieldHash
				fieldDecoder6 = fieldDecoder
			} else {
				fieldName7 = fieldHash
				fieldDecoder7 = fieldDecoder
			}
		}
		return &sevenFieldsStructDecoder{typ,
			fieldName1, fieldDecoder1,
			fieldName2, fieldDecoder2,
			fieldName3, fieldDecoder3,
			fieldName4, fieldDecoder4,
			fieldName5, fieldDecoder5,
			fieldName6, fieldDecoder6,
			fieldName7, fieldDecoder7}
	case 8:
		var fieldName1 int64
		var fieldName2 int64
		var fieldName3 int64
		var fieldName4 int64
		var fieldName5 int64
		var fieldName6 int64
		var fieldName7 int64
		var fieldName8 int64
		var fieldDecoder1 *StructFieldDecoder
		var fieldDecoder2 *StructFieldDecoder
		var fieldDecoder3 *StructFieldDecoder
		var fieldDecoder4 *StructFieldDecoder
		var fieldDecoder5 *StructFieldDecoder
		var fieldDecoder6 *StructFieldDecoder
		var fieldDecoder7 *StructFieldDecoder
		var fieldDecoder8 *StructFieldDecoder
		for fieldName, fieldDecoder := range fields {
			fieldHash := calcHash(fieldName, ctx.caseSensitive())
			_, known := knownHash[fieldHash]
			if known {
				return &generalStructDecoder{typ, fields, false}
			}
			knownHash[fieldHash] = struct{}{}
			if fieldName1 == 0 {
				fieldName1 = fieldHash
				fieldDecoder1 = fieldDecoder
			} else if fieldName2 == 0 {
				fieldName2 = fieldHash
				fieldDecoder2 = fieldDecoder
			} else if fieldName3 == 0 {
				fieldName3 = fieldHash
				fieldDecoder3 = fieldDecoder
			} else if fieldName4 == 0 {
				fieldName4 = fieldHash
				fieldDecoder4 = fieldDecoder
			} else if fieldName5 == 0 {
				fieldName5 = fieldHash
				fieldDecoder5 = fieldDecoder
			} else if fieldName6 == 0 {
				fieldName6 = fieldHash
				fieldDecoder6 = fieldDecoder
			} else if fieldName7 == 0 {
				fieldName7 = fieldHash
				fieldDecoder7 = fieldDecoder
			} else {
				fieldName8 = fieldHash
				fieldDecoder8 = fieldDecoder
			}
		}
		return &eightFieldsStructDecoder{typ,
			fieldName1, fieldDecoder1,
			fieldName2, fieldDecoder2,
			fieldName3, fieldDecoder3,
			fieldName4, fieldDecoder4,
			fieldName5, fieldDecoder5,
			fieldName6, fieldDecoder6,
			fieldName7, fieldDecoder7,
			fieldName8, fieldDecoder8}
	case 9:
		var fieldName1 int64
		var fieldName2 int64
		var fieldName3 int64
		var fieldName4 int64
		var fieldName5 int64
		var fieldName6 int64
		var fieldName7 int64
		var fieldName8 int64
		var fieldName9 int64
		var fieldDecoder1 *StructFieldDecoder
		var fieldDecoder2 *StructFieldDecoder
		var fieldDecoder3 *StructFieldDecoder
		var fieldDecoder4 *StructFieldDecoder
		var fieldDecoder5 *StructFieldDecoder
		var fieldDecoder6 *StructFieldDecoder
		var fieldDecoder7 *StructFieldDecoder
		var fieldDecoder8 *StructFieldDecoder
		var fieldDecoder9 *StructFieldDecoder
		for fieldName, fieldDecoder := range fields {
			fieldHash := calcHash(fieldName, ctx.caseSensitive())
			_, known := knownHash[fieldHash]
			if known {
				return &generalStructDecoder{typ, fields, false}
			}
			knownHash[fieldHash] = struct{}{}
			if fieldName1 == 0 {
				fieldName1 = fieldHash
				fieldDecoder1 = fieldDecoder
			} else if fieldName2 == 0 {
				fieldName2 = fieldHash
				fieldDecoder2 = fieldDecoder
			} else if fieldName3 == 0 {
				fieldName3 = fieldHash
				fieldDecoder3 = fieldDecoder
			} else if fieldName4 == 0 {
				fieldName4 = fieldHash
				fieldDecoder4 = fieldDecoder
			} else if fieldName5 == 0 {
				fieldName5 = fieldHash
				fieldDecoder5 = fieldDecoder
			} else if fieldName6 == 0 {
				fieldName6 = fieldHash
				fieldDecoder6 = fieldDecoder
			} else if fieldName7 == 0 {
				fieldName7 = fieldHash
				fieldDecoder7 = fieldDecoder
			} else if fieldName8 == 0 {
				fieldName8 = fieldHash
				fieldDecoder8 = fieldDecoder
			} else {
				fieldName9 = fieldHash
				fieldDecoder9 = fieldDecoder
			}
		}
		return &nineFieldsStructDecoder{typ,
			fieldName1, fieldDecoder1,
			fieldName2, fieldDecoder2,
			fieldName3, fieldDecoder3,
			fieldName4, fieldDecoder4,
			fieldName5, fieldDecoder5,
			fieldName6, fieldDecoder6,
			fieldName7, fieldDecoder7,
			fieldName8, fieldDecoder8,
			fieldName9, fieldDecoder9}
	case 10:
		var fieldName1 int64
		var fieldName2 int64
		var fieldName3 int64
		var fieldName4 int64
		var fieldName5 int64
		var fieldName6 int64
		var fieldName7 int64
		var fieldName8 int64
		var fieldName9 int64
		var fieldName10 int64
		var fieldDecoder1 *StructFieldDecoder
		var fieldDecoder2 *StructFieldDecoder
		var fieldDecoder3 *StructFieldDecoder
		var fieldDecoder4 *StructFieldDecoder
		var fieldDecoder5 *StructFieldDecoder
		var fieldDecoder6 *StructFieldDecoder
		var fieldDecoder7 *StructFieldDecoder
		var fieldDecoder8 *StructFieldDecoder
		var fieldDecoder9 *StructFieldDecoder
		var fieldDecoder10 *StructFieldDecoder
		for fieldName, fieldDecoder := range fields {
			fieldHash := calcHash(fieldName, ctx.caseSensitive())
			_, known := knownHash[fieldHash]
			if known {
				return &generalStructDecoder{typ, fields, false}
			}
			knownHash[fieldHash] = struct{}{}
			if fieldName1 == 0 {
				fieldName1 = fieldHash
				fieldDecoder1 = fieldDecoder
			} else if fieldName2 == 0 {
				fieldName2 = fieldHash
				fieldDecoder2 = fieldDecoder
			} else if fieldName3 == 0 {
				fieldName3 = fieldHash
				fieldDecoder3 = fieldDecoder
			} else if fieldName4 == 0 {
				fieldName4 = fieldHash
				fieldDecoder4 = fieldDecoder
			} else if fieldName5 == 0 {
				fieldName5 = fieldHash
				fieldDecoder5 = fieldDecoder
			} else if fieldName6 == 0 {
				fieldName6 = fieldHash
				fieldDecoder6 = fieldDecoder
			} else if fieldName7 == 0 {
				fieldName7 = fieldHash
				fieldDecoder7 = fieldDecoder
			} else if fieldName8 == 0 {
				fieldName8 = fieldHash
				fieldDecoder8 = fieldDecoder
			} else if fieldName9 == 0 {
				fieldName9 = fieldHash
				fieldDecoder9 = fieldDecoder
			} else {
				fieldName10 = fieldHash
				fieldDecoder10 = fieldDecoder
			}
		}
		return &tenFieldsStructDecoder{typ,
			fieldName1, fieldDecoder1,
			fieldName2, fieldDecoder2,
			fieldName3, fieldDecoder3,
			fieldName4, fieldDecoder4,
			fieldName5, fieldDecoder5,
			fieldName6, fieldDecoder6,
			fieldName7, fieldDecoder7,
			fieldName8, fieldDecoder8,
			fieldName9, fieldDecoder9,
			fieldName10, fieldDecoder10}
	}
	return &generalStructDecoder{typ, fields, false}
}

type generalStructDecoder struct {
	typ                   reflect2.Type
	fields                map[string]*StructFieldDecoder
	disallowUnknownFields bool
}

func (decoder *generalStructDecoder) Decode(ptr unsafe.Pointer, iter *Iterator) {
	if !iter.readObjectStart() {
		return
	}
	if !iter.incrementDepth() {
		return
	}
	var c byte
	for c = ','; c == ','; c = iter.nextToken() {
		decoder.decodeOneField(ptr, iter)
	}
	if iter.Error != nil && iter.Error != io.EOF && len(decoder.typ.Type1().Name()) != 0 {
		iter.Error = fmt.Errorf("%v.%s", decoder.typ, iter.Error.Error())
	}
	if c != '}' {
		iter.ReportError("struct Decode", `expect }, but found `+string([]byte{c}))
	}
	iter.decrementDepth()
}

func (decoder *generalStructDecoder) decodeOneField(ptr unsafe.Pointer, iter *Iterator) {
	var field string
	var fieldDecoder *StructFieldDecoder
	if iter.cfg.objectFieldMustBeSimpleString {
		fieldBytes := iter.ReadStringAsSlice()
		field = *(*string)(unsafe.Pointer(&fieldBytes))
		fieldDecoder = decoder.fields[field]
		if fieldDecoder == nil && !iter.cfg.caseSensitive {
			fieldDecoder = decoder.fields[strings.ToLower(field)]
		}
	} else {
		field = iter.ReadString()
		fieldDecoder = decoder.fields[field]
		if fieldDecoder == nil && !iter.cfg.caseSensitive {
			fieldDecoder = decoder.fields[strings.ToLower(field)]
		}
	}
	if fieldDecoder == nil {
		if decoder.disallowUnknownFields {
			msg := "found unknown field: " + field
			iter.ReportError("ReadObject", msg)
		}
		c := iter.nextToken()
		if c != ':' {
			iter.ReportError("ReadObject", "expect : after object field, but found "+string([]byte{c}))
		}
		iter.Skip()
		return
	}
	c := iter.nextToken()
	if c != ':' {
		iter.ReportError("ReadObject", "expect : after object field, but found "+string([]byte{c}))
	}
	fieldDecoder.Decode(ptr, iter)
}

type skipObjectDecoder struct {
	typ reflect2.Type
}

func (decoder *skipObjectDecoder) Decode(ptr unsafe.Pointer, iter *Iterator) {
	valueType := iter.WhatIsNext()
	if valueType != ObjectValue && valueType != NilValue {
		iter.ReportError("skipObjectDecoder", "expect object or null")
		return
	}
	iter.Skip()
}

type oneFieldStructDecoder struct {
	typ          reflect2.Type
	fieldHash    int64
	fieldDecoder *StructFieldDecoder
}

func (decoder *oneFieldStructDecoder) Decode(ptr unsafe.Pointer, iter *Iterator) {
	if !iter.readObjectStart() {
		return
	}
	if !iter.incrementDepth() {
		return
	}
	for {
		if iter.readFieldHash() == decoder.fieldHash {
			decoder.fieldDecoder.Decode(ptr, iter)
		} else {
			iter.Skip()
		}
		if iter.isObjectEnd() {
			break
		}
	}
	if iter.Error != nil && iter.Error != io.EOF && len(decoder.typ.Type1().Name()) != 0 {
		iter.Error = fmt.Errorf("%v.%s", decoder.typ, iter.Error.Error())
	}
	iter.decrementDepth()
}

type twoFieldsStructDecoder struct {
	typ           reflect2.Type
	fieldHash1    int64
	fieldDecoder1 *StructFieldDecoder
	fieldHash2    int64
	fieldDecoder2 *StructFieldDecoder
}

func (decoder *twoFieldsStructDecoder) Decode(ptr unsafe.Pointer, iter *Iterator) {
	if !iter.readObjectStart() {
		return
	}
	if !iter.incrementDepth() {
		return
	}
	for {
		switch iter.readFieldHash() {
		case decoder.fieldHash1:
			decoder.fieldDecoder1.Decode(ptr, iter)
		case decoder.fieldHash2:
			decoder.fieldDecoder2.Decode(ptr, iter)
		default:
			iter.Skip()
		}
		if iter.isObjectEnd() {
			break
		}
	}
	if iter.Error != nil && iter.Error != io.EOF && len(decoder.typ.Type1().Name()) != 0 {
		iter.Error = fmt.Errorf("%v.%s", decoder.typ, iter.Error.Error())
	}
	iter.decrementDepth()
}

type threeFieldsStructDecoder struct {
	typ           reflect2.Type
	fieldHash1    int64
	fieldDecoder1 *StructFieldDecoder
	fieldHash2    int64
	fieldDecoder2 *StructFieldDecoder
	fieldHash3    int64
	fieldDecoder3 *StructFieldDecoder
}

func (decoder *threeFieldsStructDecoder) Decode(ptr unsafe.Pointer, iter *Iterator) {
	if !iter.readObjectStart() {
		return
	}
	if !iter.incrementDepth() {
		return
	}
	for {
		switch iter.readFieldHash() {
		case decoder.fieldHash1:
			decoder.fieldDecoder1.Decode(ptr, iter)
		case decoder.fieldHash2:
			decoder.fieldDecoder2.Decode(ptr, iter)
		case decoder.fieldHash3:
			decoder.fieldDecoder3.Decode(ptr, iter)
		default:
			iter.Skip()
		}
		if iter.isObjectEnd() {
			break
		}
	}
	if iter.Error != nil && iter.Error != io.EOF && len(decoder.typ.Type1().Name()) != 0 {
		iter.Error = fmt.Errorf("%v.%s", decoder.typ, iter.Error.Error())
	}
	iter.decrementDepth()
}

type fourFieldsStructDecoder struct {
	typ           reflect2.Type
	fieldHash1    int64
	fieldDecoder1 *StructFieldDecoder
	fieldHash2    int64
	fieldDecoder2 *StructFieldDecoder
	fieldHash3    int64
	fieldDecoder3 *StructFieldDecoder
	fieldHash4    int64
	fieldDecoder4 *StructFieldDecoder
}

func (decoder *fourFieldsStructDecoder) Decode(ptr unsafe.Pointer, iter *Iterator) {
	if !iter.readObjectStart() {
		return
	}
	if !iter.incrementDepth() {
		return
	}
	for {
		switch iter.readFieldHash() {
		case decoder.fieldHash1:
			decoder.fieldDecoder1.Decode(ptr, iter)
		case decoder.fieldHash2:
			decoder.fieldDecoder2.Decode(ptr, iter)
		case decoder.fieldHash3:
			decoder.fieldDecoder3.Decode(ptr, iter)
		case decoder.fieldHash4:
			decoder.fieldDecoder4.Decode(ptr, iter)
		default:
			iter.Skip()
		}
		if iter.isObjectEnd() {
			break
		}
	}
	if iter.Error != nil && iter.Error != io.EOF && len(decoder.typ.Type1().Name()) != 0 {
		iter.Error = fmt.Errorf("%v.%s", decoder.typ, iter.Error.Error())
	}
	iter.decrementDepth()
}

type fiveFieldsStructDecoder struct {
	typ           reflect2.Type
	fieldHash1    int64
	fieldDecoder1 *StructFieldDecoder
	fieldHash2    int64
	fieldDecoder2 *StructFieldDecoder
	fieldHash3    int64
	fieldDecoder3 *StructFieldDecoder
	fieldHash4    int64
	fieldDecoder4 *StructFieldDecoder
	fieldHash5    int64
	fieldDecoder5 *StructFieldDecoder
}

func (decoder *fiveFieldsStructDecoder) Decode(ptr unsafe.Pointer, iter *Iterator) {
	if !iter.readObjectStart() {
		return
	}
	if !iter.incrementDepth() {
		return
	}
	for {
		switch iter.readFieldHash() {
		case decoder.fieldHash1:
			decoder.fieldDecoder1.Decode(ptr, iter)
		case decoder.fieldHash2:
			decoder.fieldDecoder2.Decode(ptr, iter)
		case decoder.fieldHash3:
			decoder.fieldDecoder3.Decode(ptr, iter)
		case decoder.fieldHash4:
			decoder.fieldDecoder4.Decode(ptr, iter)
		case decoder.fieldHash5:
			decoder.fieldDecoder5.Decode(ptr, iter)
		default:
			iter.Skip()
		}
		if iter.isObjectEnd() {
			break
		}
	}
	if iter.Error != nil && iter.Error != io.EOF && len(decoder.typ.Type1().Name()) != 0 {
		iter.Error = fmt.Errorf("%v.%s", decoder.typ, iter.Error.Error())
	}
	iter.decrementDepth()
}

type sixFieldsStructDecoder struct {
	typ           reflect2.Type
	fieldHash1    int64
	fieldDecoder1 *StructFieldDecoder
	fieldHash2    int64
	fieldDecoder2 *StructFieldDecoder
	fieldHash3    int64
	fieldDecoder3 *StructFieldDecoder
	fieldHash4    int64
	fieldDecoder4 *StructFieldDecoder
	fieldHash5    int64
	fieldDecoder5 *StructFieldDecoder
	fieldHash6    int64
	fieldDecoder6 *StructFieldDecoder
}

func (decoder *sixFieldsStructDecoder) Decode(ptr unsafe.Pointer, iter *Iterator) {
	if !iter.readObjectStart() {
		return
	}
	if !iter.incrementDepth() {
		return
	}
	for {
		switch iter.readFieldHash() {
		case decoder.fieldHash1:
			decoder.fieldDecoder1.Decode(ptr, iter)
		case decoder.fieldHash2:
			decoder.fieldDecoder2.Decode(ptr, iter)
		case decoder.fieldHash3:
			decoder.fieldDecoder3.Decode(ptr, iter)
		case decoder.fieldHash4:
			decoder.fieldDecoder4.Decode(ptr, iter)
		case decoder.fieldHash5:
			decoder.fieldDecoder5.Decode(ptr, iter)
		case decoder.fieldHash6:
			decoder.fieldDecoder6.Decode(ptr, iter)
		default:
			iter.Skip()
		}
		if iter.isObjectEnd() {
			break
		}
	}
	if iter.Error != nil && iter.Error != io.EOF && len(decoder.typ.Type1().Name()) != 0 {
		iter.Error = fmt.Errorf("%v.%s", decoder.typ, iter.Error.Error())
	}
	iter.decrementDepth()
}

type sevenFieldsStructDecoder struct {
	typ           reflect2.Type
	fieldHash1    int64
	fieldDecoder1 *StructFieldDecoder
	fieldHash2    int64
	fieldDecoder2 *StructFieldDecoder
	fieldHash3    int64
	fieldDecoder3 *StructFieldDecoder
	fieldHash4    int64
	fieldDecoder4 *StructFieldDecoder
	fieldHash5    int64
	fieldDecoder5 *StructFieldDecoder
	fieldHash6    int64
	fieldDecoder6 *StructFieldDecoder
	fieldHash7    int64
	fieldDecoder7 *StructFieldDecoder
}

func (decoder *sevenFieldsStructDecoder) Decode(ptr unsafe.Pointer, iter *Iterator) {
	if !iter.readObjectStart() {
		return
	}
	if !iter.incrementDepth() {
		return
	}
	for {
		switch iter.readFieldHash() {
		case decoder.fieldHash1:
			decoder.fieldDecoder1.Decode(ptr, iter)
		case decoder.fieldHash2:
			decoder.fieldDecoder2.Decode(ptr, iter)
		case decoder.fieldHash3:
			decoder.fieldDecoder3.Decode(ptr, iter)
		case decoder.fieldHash4:
			decoder.fieldDecoder4.Decode(ptr, iter)
		case decoder.fieldHash5:
			decoder.fieldDecoder5.Decode(ptr, iter)
		case decoder.fieldHash6:
			decoder.fieldDecoder6.Decode(ptr, iter)
		case decoder.fieldHash7:
			decoder.fieldDecoder7.Decode(ptr, iter)
		default:
			iter.Skip()
		}
		if iter.isObjectEnd() {
			break
		}
	}
	if iter.Error != nil && iter.Error != io.EOF && len(decoder.typ.Type1().Name()) != 0 {
		iter.Error = fmt.Errorf("%v.%s", decoder.typ, iter.Error.Error())
	}
	iter.decrementDepth()
}

type eightFieldsStructDecoder struct {
	typ           reflect2.Type
	fieldHash1    int64
	fieldDecoder1 *StructFieldDecoder
	fieldHash2    int64
	fieldDecoder2 *StructFieldDecoder
	fieldHash3    int64
	fieldDecoder3 *StructFieldDecoder
	fieldHash4    int64
	fieldDecoder4 *StructFieldDecoder
	fieldHash5    int64
	fieldDecoder5 *StructFieldDecoder
	fieldHash6    int64
	fieldDecoder6 *StructFieldDecoder
	fieldHash7    int64
	fieldDecoder7 *StructFieldDecoder
	fieldHash8    int64
	fieldDecoder8 *StructFieldDecoder
}

func (decoder *eightFieldsStructDecoder) Decode(ptr unsafe.Pointer, iter *Iterator) {
	if !iter.readObjectStart() {
		return
	}
	if !iter.incrementDepth() {
		return
	}
	for {
		switch iter.readFieldHash() {
		case decoder.fieldHash1:
			decoder.fieldDecoder1.Decode(ptr, iter)
		case decoder.fieldHash2:
			decoder.fieldDecoder2.Decode(ptr, iter)
		case decoder.fieldHash3:
			decoder.fieldDecoder3.Decode(ptr, iter)
		case decoder.fieldHash4:
			decoder.fieldDecoder4.Decode(ptr, iter)
		case decoder.fieldHash5:
			decoder.fieldDecoder5.Decode(ptr, iter)
		case decoder.fieldHash6:
			decoder.fieldDecoder6.Decode(ptr, iter)
		case decoder.fieldHash7:
			decoder.fieldDecoder7.Decode(ptr, iter)
		case decoder.fieldHash8:
			decoder.fieldDecoder8.Decode(ptr, iter)
		default:
			iter.Skip()
		}
		if iter.isObjectEnd() {
			break
		}
	}
	if iter.Error != nil && iter.Error != io.EOF && len(decoder.typ.Type1().Name()) != 0 {
		iter.Error = fmt.Errorf("%v.%s", decoder.typ, iter.Error.Error())
	}
	iter.decrementDepth()
}

type nineFieldsStructDecoder struct {
	typ           reflect2.Type
	fieldHash1    int64
	fieldDecoder1 *StructFieldDecoder
	fieldHash2    int64
	fieldDecoder2 *StructFieldDecoder
	fieldHash3    int64
	fieldDecoder3 *StructFieldDecoder
	fieldHash4    int64
	fieldDecoder4 *StructFieldDecoder
	fieldHash5    int64
	fieldDecoder5 *StructFieldDecoder
	fieldHash6    int64
	fieldDecoder6 *StructFieldDecoder
	fieldHash7    int64
	fieldDecoder7 *StructFieldDecoder
	fieldHash8    int64
	fieldDecoder8 *StructFieldDecoder
	fieldHash9    int64
	fieldDecoder9 *StructFieldDecoder
}

func (decoder *nineFieldsStructDecoder) Decode(ptr unsafe.Pointer, iter *Iterator) {
	if !iter.readObjectStart() {
		return
	}
	if !iter.incrementDepth() {
		return
	}
	for {
		switch iter.readFieldHash() {
		case decoder.fieldHash1:
			decoder.fieldDecoder1.Decode(ptr, iter)
		case decoder.fieldHash2:
			decoder.fieldDecoder2.Decode(ptr, iter)
		case decoder.fieldHash3:
			decoder.fieldDecoder3.Decode(ptr, iter)
		case decoder.fieldHash4:
			decoder.fieldDecoder4.Decode(ptr, iter)
		case decoder.fieldHash5:
			decoder.fieldDecoder5.Decode(ptr, iter)
		case decoder.fieldHash6:
			decoder.fieldDecoder6.Decode(ptr, iter)
		case decoder.fieldHash7:
			decoder.fieldDecoder7.Decode(ptr, iter)
		case decoder.fieldHash8:
			decoder.fieldDecoder8.Decode(ptr, iter)
		case decoder.fieldHash9:
			decoder.fieldDecoder9.Decode(ptr, iter)
		default:
			iter.Skip()
		}
		if iter.isObjectEnd() {
			break
		}
	}
	if iter.Error != nil && iter.Error != io.EOF && len(decoder.typ.Type1().Name()) != 0 {
		iter.Error = fmt.Errorf("%v.%s", decoder.typ, iter.Error.Error())
	}
	iter.decrementDepth()
}

type tenFieldsStructDecoder struct {
	typ            reflect2.Type
	fieldHash1     int64
	fieldDecoder1  *StructFieldDecoder
	fieldHash2     int64
	fieldDecoder2  *StructFieldDecoder
	fieldHash3     int64
	fieldDecoder3  *StructFieldDecoder
	fieldHash4     int64
	fieldDecoder4  *StructFieldDecoder
	fieldHash5     int64
	fieldDecoder5  *StructFieldDecoder
	fieldHash6     int64
	fieldDecoder6  *StructFieldDecoder
	fieldHash7     int64
	fieldDecoder7  *StructFieldDecoder
	fieldHash8     int64
	fieldDecoder8  *StructFieldDecoder
	fieldHash9     int64
	fieldDecoder9  *StructFieldDecoder
	fieldHash10    int64
	fieldDecoder10 *StructFieldDecoder
}

func (decoder *tenFieldsStructDecoder) Decode(ptr unsafe.Pointer, iter *Iterator) {
	if !iter.readObjectStart() {
		return
	}
	if !iter.incrementDepth() {
		return
	}
	for {
		switch iter.readFieldHash() {
		case decoder.fieldHash1:
			decoder.fieldDecoder1.Decode(ptr, iter)
		case decoder.fieldHash2:
			decoder.fieldDecoder2.Decode(ptr, iter)
		case decoder.fieldHash3:
			decoder.fieldDecoder3.Decode(ptr, iter)
		case decoder.fieldHash4:
			decoder.fieldDecoder4.Decode(ptr, iter)
		case decoder.fieldHash5:
			decoder.fieldDecoder5.Decode(ptr, iter)
		case decoder.fieldHash6:
			decoder.fieldDecoder6.Decode(ptr, iter)
		case decoder.fieldHash7:
			decoder.fieldDecoder7.Decode(ptr, iter)
		case decoder.fieldHash8:
			decoder.fieldDecoder8.Decode(ptr, iter)
		case decoder.fieldHash9:
			decoder.fieldDecoder9.Decode(ptr, iter)
		case decoder.fieldHash10:
			decoder.fieldDecoder10.Decode(ptr, iter)
		default:
			iter.Skip()
		}
		if iter.isObjectEnd() {
			break
		}
	}
	if iter.Error != nil && iter.Error != io.EOF && len(decoder.typ.Type1().Name()) != 0 {
		iter.Error = fmt.Errorf("%v.%s", decoder.typ, iter.Error.Error())
	}
	iter.decrementDepth()
}

type StructFieldDecoder struct {
	Field        reflect2.StructField
	FieldDecoder ValDecoder
}

func (decoder *StructFieldDecoder) Decode(ptr unsafe.Pointer, iter *Iterator) {
	fieldPtr := decoder.Field.UnsafeGet(ptr)
	decoder.FieldDecoder.Decode(fieldPtr, iter)
	if iter.Error != nil && iter.Error != io.EOF {
		iter.Error = fmt.Errorf("%s: %s", decoder.Field.Name(), iter.Error.Error())
	}
}

type stringModeStringDecoder struct {
	elemDecoder ValDecoder
	cfg         *frozenConfig
}

func (decoder *stringModeStringDecoder) Decode(ptr unsafe.Pointer, iter *Iterator) {
	decoder.elemDecoder.Decode(ptr, iter)
	str := *((*string)(ptr))
	tempIter := decoder.cfg.BorrowIterator([]byte(str))
	defer decoder.cfg.ReturnIterator(tempIter)
	*((*string)(ptr)) = tempIter.ReadString()
}

type stringModeNumberDecoder struct {
	elemDecoder ValDecoder
}

func (decoder *stringModeNumberDecoder) Decode(ptr unsafe.Pointer, iter *Iterator) {
	if iter.WhatIsNext() == NilValue {
		decoder.elemDecoder.Decode(ptr, iter)
		return
	}

	c := iter.nextToken()
	if c != '"' {
		iter.ReportError("stringModeNumberDecoder", `expect ", but found `+string([]byte{c}))
		return
	}
	decoder.elemDecoder.Decode(ptr, iter)
	if iter.Error != nil {
		return
	}
	c = iter.readByte()
	if c != '"' {
		iter.ReportError("stringModeNumberDecoder", `expect ", but found `+string([]byte{c}))
		return
	}
}
