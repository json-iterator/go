package jsoniter

var digits []uint8
var digitTens []uint8
var digitOnes []uint8
var DIGITS []uint32

func init() {
	digits = []uint8{
		'0', '1', '2', '3', '4', '5',
		'6', '7', '8', '9', 'a', 'b',
		'c', 'd', 'e', 'f', 'g', 'h',
		'i', 'j', 'k', 'l', 'm', 'n',
		'o', 'p', 'q', 'r', 's', 't',
		'u', 'v', 'w', 'x', 'y', 'z',
	}
	digitTens = []uint8{
		'0', '0', '0', '0', '0', '0', '0', '0', '0', '0',
		'1', '1', '1', '1', '1', '1', '1', '1', '1', '1',
		'2', '2', '2', '2', '2', '2', '2', '2', '2', '2',
		'3', '3', '3', '3', '3', '3', '3', '3', '3', '3',
		'4', '4', '4', '4', '4', '4', '4', '4', '4', '4',
		'5', '5', '5', '5', '5', '5', '5', '5', '5', '5',
		'6', '6', '6', '6', '6', '6', '6', '6', '6', '6',
		'7', '7', '7', '7', '7', '7', '7', '7', '7', '7',
		'8', '8', '8', '8', '8', '8', '8', '8', '8', '8',
		'9', '9', '9', '9', '9', '9', '9', '9', '9', '9',
	}
	digitOnes = []uint8{
		'0', '1', '2', '3', '4', '5', '6', '7', '8', '9',
		'0', '1', '2', '3', '4', '5', '6', '7', '8', '9',
		'0', '1', '2', '3', '4', '5', '6', '7', '8', '9',
		'0', '1', '2', '3', '4', '5', '6', '7', '8', '9',
		'0', '1', '2', '3', '4', '5', '6', '7', '8', '9',
		'0', '1', '2', '3', '4', '5', '6', '7', '8', '9',
		'0', '1', '2', '3', '4', '5', '6', '7', '8', '9',
		'0', '1', '2', '3', '4', '5', '6', '7', '8', '9',
		'0', '1', '2', '3', '4', '5', '6', '7', '8', '9',
		'0', '1', '2', '3', '4', '5', '6', '7', '8', '9',
	}
	DIGITS = make([]uint32, 1000)
	for i := uint32(0); i < 1000; i++ {
		DIGITS[i] = (((i / 100) + '0') << 16) + ((((i / 10) % 10) + '0') << 8) + i % 10 + '0';
		if i < 10 {
			DIGITS[i] += 2 << 24
		} else if i < 100 {
			DIGITS[i] += 1 << 24
		}
	}
}

func writeFirstBuf(buf []byte, v uint32, n int) int {
	start := v >> 24
	if start == 0 {
		buf[n] = byte(v >> 16)
		n++
		buf[n] = byte(v >> 8)
		n++
	} else if start == 1 {
		buf[n] = byte(v >> 8)
		n++
	}
	buf[n] = byte(v)
	n++
	return n
}

func writeBuf(buf []byte, v uint32, n int) {
	buf[n] = byte(v >> 16)
	buf[n + 1] = byte(v >> 8)
	buf[n + 2] = byte(v)
}

func (stream *Stream) WriteUint8(val uint8) {
	if stream.Available() < 3 {
		stream.Flush()
	}
	stream.n = writeFirstBuf(stream.buf, DIGITS[val], stream.n)
}

func (stream *Stream) WriteInt8(nval int8) {
	if stream.Available() < 4 {
		stream.Flush()
	}
	n := stream.n
	var val uint8
	if (nval < 0) {
		val = uint8(-nval)
		stream.buf[n] = '-'
		n++
	} else {
		val = uint8(nval)
	}
	stream.n = writeFirstBuf(stream.buf, DIGITS[val], n)
}

func (stream *Stream) WriteUint16(val uint16) {
	if stream.Available() < 5 {
		stream.Flush()
	}
	q1 := val / 1000
	if q1 == 0 {
		stream.n = writeFirstBuf(stream.buf, DIGITS[val], stream.n)
		return
	}
	r1 := val - q1 * 1000;
	n := writeFirstBuf(stream.buf, DIGITS[q1], stream.n)
	writeBuf(stream.buf, DIGITS[r1], n)
	stream.n = n + 3
	return
}

func (stream *Stream) WriteInt16(nval int16) {
	if stream.Available() < 6 {
		stream.Flush()
	}
	n := stream.n
	var val uint16
	if (nval < 0) {
		val = uint16(-nval)
		stream.buf[n] = '-'
		n++
	} else {
		val = uint16(nval)
	}
	q1 := val / 1000
	if q1 == 0 {
		stream.n = writeFirstBuf(stream.buf, DIGITS[val], n)
		return
	}
	r1 := val - q1 * 1000;
	n = writeFirstBuf(stream.buf, DIGITS[q1], n)
	writeBuf(stream.buf, DIGITS[r1], n)
	stream.n = n + 3
	return
}

func (stream *Stream) WriteUint32(val uint32) {
	if stream.Available() < 10 {
		stream.Flush()
	}
	n := stream.n
	q1 := val / 1000
	if q1 == 0 {
		stream.n = writeFirstBuf(stream.buf, DIGITS[val], n)
		return
	}
	r1 := val - q1 * 1000;
	q2 := q1 / 1000
	if q2 == 0 {
		n := writeFirstBuf(stream.buf, DIGITS[q1], n)
		writeBuf(stream.buf, DIGITS[r1], n)
		stream.n = n + 3
		return
	}
	r2 := q1 - q2 * 1000
	q3 := q2 / 1000
	if q3 == 0 {
		n = writeFirstBuf(stream.buf, DIGITS[q2], n)
	} else {
		r3 := q2 - q3 * 1000
		stream.buf[n] = byte(q3 + '0')
		n++
		writeBuf(stream.buf, DIGITS[r3], n)
		n += 3
	}
	writeBuf(stream.buf, DIGITS[r2], n)
	writeBuf(stream.buf, DIGITS[r1], n + 3)
	stream.n = n + 6
}

func (stream *Stream) WriteInt32(nval int32) {
	if stream.Available() < 11 {
		stream.Flush()
	}
	n := stream.n
	var val uint32
	if (nval < 0) {
		val = uint32(-nval)
		stream.buf[n] = '-'
		n++
	} else {
		val = uint32(nval)
	}
	q1 := val / 1000
	if q1 == 0 {
		stream.n = writeFirstBuf(stream.buf, DIGITS[val], n)
		return
	}
	r1 := val - q1 * 1000;
	q2 := q1 / 1000
	if q2 == 0 {
		n := writeFirstBuf(stream.buf, DIGITS[q1], n)
		writeBuf(stream.buf, DIGITS[r1], n)
		stream.n = n + 3
		return
	}
	r2 := q1 - q2 * 1000
	q3 := q2 / 1000
	if q3 == 0 {
		n = writeFirstBuf(stream.buf, DIGITS[q2], n)
	} else {
		r3 := q2 - q3 * 1000
		stream.buf[n] = byte(q3 + '0')
		n++
		writeBuf(stream.buf, DIGITS[r3], n)
		n += 3
	}
	writeBuf(stream.buf, DIGITS[r2], n)
	writeBuf(stream.buf, DIGITS[r1], n + 3)
	stream.n = n + 6
}

func (stream *Stream) WriteUint64(val uint64) {
	if stream.Available() < 20 {
		stream.Flush()
	}
	n := stream.n
	q1 := val / 1000
	if q1 == 0 {
		stream.n = writeFirstBuf(stream.buf, DIGITS[val], n)
		return
	}
	r1 := val - q1 * 1000;
	q2 := q1 / 1000
	if q2 == 0 {
		n := writeFirstBuf(stream.buf, DIGITS[q1], n)
		writeBuf(stream.buf, DIGITS[r1], n)
		stream.n = n + 3
		return
	}
	r2 := q1 - q2 * 1000
	q3 := q2 / 1000
	if q3 == 0 {
		n = writeFirstBuf(stream.buf, DIGITS[q2], n)
		writeBuf(stream.buf, DIGITS[r2], n)
		writeBuf(stream.buf, DIGITS[r1], n + 3)
		stream.n = n + 6
		return
	}
	r3 := q2 - q3 * 1000
	q4 := q3 / 1000
	if q4 == 0 {
		n = writeFirstBuf(stream.buf, DIGITS[q3], n)
		writeBuf(stream.buf, DIGITS[r3], n)
		writeBuf(stream.buf, DIGITS[r2], n + 3)
		writeBuf(stream.buf, DIGITS[r1], n + 6)
		stream.n = n + 9
		return
	}
	r4 := q3 - q4 * 1000
	q5 := q4 / 1000
	if q5 == 0 {
		n = writeFirstBuf(stream.buf, DIGITS[q4], n)
		writeBuf(stream.buf, DIGITS[r4], n)
		writeBuf(stream.buf, DIGITS[r3], n + 3)
		writeBuf(stream.buf, DIGITS[r2], n + 6)
		writeBuf(stream.buf, DIGITS[r1], n + 9)
		stream.n = n + 12
		return
	}
	r5 := q4 - q5 * 1000
	q6 := q5 / 1000
	if q6 == 0 {
		n = writeFirstBuf(stream.buf, DIGITS[q5], n)
	} else {
		n = writeFirstBuf(stream.buf, DIGITS[q6], n)
		r6 := q5 - q6 * 1000
		writeBuf(stream.buf, DIGITS[r6], n)
		n += 3
	}
	writeBuf(stream.buf, DIGITS[r5], n)
	writeBuf(stream.buf, DIGITS[r4], n + 3)
	writeBuf(stream.buf, DIGITS[r3], n + 6)
	writeBuf(stream.buf, DIGITS[r2], n + 9)
	writeBuf(stream.buf, DIGITS[r1], n + 12)
	stream.n = n + 15
}

func (stream *Stream) WriteInt64(nval int64) {
	if stream.Available() < 20 {
		stream.Flush()
	}
	n := stream.n
	var val uint64
	if (nval < 0) {
		val = uint64(-nval)
		stream.buf[n] = '-'
		n++
	} else {
		val = uint64(nval)
	}
	q1 := val / 1000
	if q1 == 0 {
		stream.n = writeFirstBuf(stream.buf, DIGITS[val], n)
		return
	}
	r1 := val - q1 * 1000;
	q2 := q1 / 1000
	if q2 == 0 {
		n := writeFirstBuf(stream.buf, DIGITS[q1], n)
		writeBuf(stream.buf, DIGITS[r1], n)
		stream.n = n + 3
		return
	}
	r2 := q1 - q2 * 1000
	q3 := q2 / 1000
	if q3 == 0 {
		n = writeFirstBuf(stream.buf, DIGITS[q2], n)
		writeBuf(stream.buf, DIGITS[r2], n)
		writeBuf(stream.buf, DIGITS[r1], n + 3)
		stream.n = n + 6
		return
	}
	r3 := q2 - q3 * 1000
	q4 := q3 / 1000
	if q4 == 0 {
		n = writeFirstBuf(stream.buf, DIGITS[q3], n)
		writeBuf(stream.buf, DIGITS[r3], n)
		writeBuf(stream.buf, DIGITS[r2], n + 3)
		writeBuf(stream.buf, DIGITS[r1], n + 6)
		stream.n = n + 9
		return
	}
	r4 := q3 - q4 * 1000
	q5 := q4 / 1000
	if q5 == 0 {
		n = writeFirstBuf(stream.buf, DIGITS[q4], n)
		writeBuf(stream.buf, DIGITS[r4], n)
		writeBuf(stream.buf, DIGITS[r3], n + 3)
		writeBuf(stream.buf, DIGITS[r2], n + 6)
		writeBuf(stream.buf, DIGITS[r1], n + 9)
		stream.n = n + 12
		return
	}
	r5 := q4 - q5 * 1000
	q6 := q5 / 1000
	if q6 == 0 {
		n = writeFirstBuf(stream.buf, DIGITS[q5], n)
	} else {
		stream.buf[n] = byte(q6 + '0')
		n++
		r6 := q5 - q6 * 1000
		writeBuf(stream.buf, DIGITS[r6], n)
		n += 3
	}
	writeBuf(stream.buf, DIGITS[r5], n)
	writeBuf(stream.buf, DIGITS[r4], n + 3)
	writeBuf(stream.buf, DIGITS[r3], n + 6)
	writeBuf(stream.buf, DIGITS[r2], n + 9)
	writeBuf(stream.buf, DIGITS[r1], n + 12)
	stream.n = n + 15
}

func (stream *Stream) WriteInt(val int) {
	stream.WriteInt64(int64(val))
}

func (stream *Stream) WriteUint(val uint) {
	stream.WriteUint64(uint64(val))
}