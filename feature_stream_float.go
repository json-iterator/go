package jsoniter

import (
	"strconv"
)

var POW10 []uint64

func init() {
	POW10 = []uint64{1, 10, 100, 1000, 10000, 100000, 1000000}
}

func (stream *Stream) WriteFloat32(val float32) {
	if val < 0 {
		stream.writeByte('-')
		val = -val
	}
	if val > 0x4ffffff {
		stream.WriteRaw(strconv.FormatFloat(float64(val), 'f', -1, 32));
		return
	}
	precision := 6
	exp := uint64(1000000) // 6
	lval := uint64(float64(val) * float64(exp) + 0.5)
	stream.WriteUint64(lval / exp)
	fval := lval % exp
	if fval == 0 {
		return
	}
	stream.writeByte('.')
	if stream.Available() < 10 {
		stream.Flush()
	}
	for p := precision - 1; p > 0 && fval < POW10[p]; p-- {
		stream.writeByte('0')
	}
	stream.WriteUint64(fval)
	for stream.buf[stream.n - 1] == '0' {
		stream.n--
	}
}

func (stream *Stream) WriteFloat64(val float64) {
	if val < 0 {
		stream.writeByte('-')
		val = -val
	}
	if val > 0x4ffffff {
		stream.WriteRaw(strconv.FormatFloat(val, 'f', -1, 64));
		return
	}
	precision := 6
	exp := uint64(1000000) // 6
	lval := uint64(val * float64(exp) + 0.5)
	stream.WriteUint64(lval / exp)
	fval := lval % exp
	if fval == 0 {
		return
	}
	stream.writeByte('.')
	if stream.Available() < 10 {
		stream.Flush()
	}
	for p := precision - 1; p > 0 && fval < POW10[p]; p-- {
		stream.writeByte('0')
	}
	stream.WriteUint64(fval)
	for stream.buf[stream.n - 1] == '0' {
		stream.n--
	}
}