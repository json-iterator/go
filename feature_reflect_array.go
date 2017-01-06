package jsoniter

import (
	"unsafe"
	"reflect"
	"io"
	"fmt"
)

type sliceDecoder struct {
	sliceType   reflect.Type
	elemType    reflect.Type
	elemDecoder Decoder
}

// sliceHeader is a safe version of SliceHeader used within this package.
type sliceHeader struct {
	Data unsafe.Pointer
	Len  int
	Cap  int
}

func (decoder *sliceDecoder) decode(ptr unsafe.Pointer, iter *Iterator) {
	decoder.doDecode(ptr, iter)
	if iter.Error != nil && iter.Error != io.EOF {
		iter.Error = fmt.Errorf("%v: %s", decoder.sliceType, iter.Error.Error())
	}
}

func (decoder *sliceDecoder) doDecode(ptr unsafe.Pointer, iter *Iterator) {
	slice := (*sliceHeader)(ptr)
	reuseSlice(slice, decoder.sliceType, 4)
	if !iter.ReadArray() {
		return
	}
	offset := uintptr(0)
	decoder.elemDecoder.decode(unsafe.Pointer(uintptr(slice.Data) + offset), iter)
	if !iter.ReadArray() {
		slice.Len = 1
		return
	}
	offset += decoder.elemType.Size()
	decoder.elemDecoder.decode(unsafe.Pointer(uintptr(slice.Data) + offset), iter)
	if !iter.ReadArray() {
		slice.Len = 2
		return
	}
	offset += decoder.elemType.Size()
	decoder.elemDecoder.decode(unsafe.Pointer(uintptr(slice.Data) + offset), iter)
	if !iter.ReadArray() {
		slice.Len = 3
		return
	}
	offset += decoder.elemType.Size()
	decoder.elemDecoder.decode(unsafe.Pointer(uintptr(slice.Data) + offset), iter)
	slice.Len = 4
	for iter.ReadArray() {
		growOne(slice, decoder.sliceType, decoder.elemType)
		offset += decoder.elemType.Size()
		decoder.elemDecoder.decode(unsafe.Pointer(uintptr(slice.Data) + offset), iter)
	}
}

// grow grows the slice s so that it can hold extra more values, allocating
// more capacity if needed. It also returns the old and new slice lengths.
func growOne(slice *sliceHeader, sliceType reflect.Type, elementType reflect.Type) {
	newLen := slice.Len + 1
	if newLen <= slice.Cap {
		slice.Len = newLen
		return
	}
	newCap := slice.Cap
	if newCap == 0 {
		newCap = 1
	} else {
		for newCap < newLen {
			if slice.Len < 1024 {
				newCap += newCap
			} else {
				newCap += newCap / 4
			}
		}
	}
	dst := unsafe.Pointer(reflect.MakeSlice(sliceType, newLen, newCap).Pointer())
	// copy old array into new array
	originalBytesCount := uintptr(slice.Len) * elementType.Size()
	srcPtr := (*[1 << 30]byte)(slice.Data)
	dstPtr := (*[1 << 30]byte)(dst)
	for i := uintptr(0); i < originalBytesCount; i++ {
		dstPtr[i] = srcPtr[i]
	}
	slice.Len = newLen
	slice.Cap = newCap
	slice.Data = dst
}

func reuseSlice(slice *sliceHeader, sliceType reflect.Type, expectedCap int) {
	if expectedCap <= slice.Cap {
		return
	}
	dst := unsafe.Pointer(reflect.MakeSlice(sliceType, 0, expectedCap).Pointer())
	slice.Cap = expectedCap
	slice.Data = dst
}