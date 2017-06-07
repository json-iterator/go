package jsoniter

import (
	"time"
	"unsafe"
)

func TimeDecoder(ptr unsafe.Pointer, iter *Iterator) {
	t := ReadTime(iter)
	*((*time.Time)(ptr)) = *t
}

// 需要实现全部标准
const (
	RFC3339     = "2006-01-02T15:04:05Z07:00"
	ANSIC       = "Mon Jan _2 15:04:05 2006"
	UnixDate    = "Mon Jan _2 15:04:05 MST 2006"
	RubyDate    = "Mon Jan 02 15:04:05 -0700 2006"
	RFC822      = "02 Jan 06 15:04 MST"
	RFC822Z     = "02 Jan 06 15:04 -0700" // RFC822 with numeric zone
	RFC850      = "Monday, 02-Jan-06 15:04:05 MST"
	RFC1123     = "Mon, 02 Jan 2006 15:04:05 MST"
	RFC1123Z    = "Mon, 02 Jan 2006 15:04:05 -0700" // RFC1123 with numeric zone
	RFC3339Nano = "2006-01-02T15:04:05.999999999Z07:00"
	Kitchen     = "3:04PM"
	// Handy time stamps.
	Stamp      = "Jan _2 15:04:05"
	StampMilli = "Jan _2 15:04:05.000"
	StampMicro = "Jan _2 15:04:05.000000"
	StampNano  = "Jan _2 15:04:05.000000000"
)

func ReadTime(iter *Iterator) *time.Time {
	var (
		bts   []byte
		times int = 0
		t     *time.Time
	)
	for c := iter.nextToken(); c != 0; c = timeNextToken(iter) {
		if c == '"' {
			times++
			if times == 2 {
				break
			} else {
				continue
			}
		}
		bts = append(bts, c)
	}
	timeStr := string(bts)
	t, iter.Error = ConvertStrToTime(timeStr)
	return t
}

func ConvertStrToTime(str string) (t *time.Time, err error) {
	// RFC3339
	t = new(time.Time)
	*t, err = time.Parse(time.RFC3339, str)
	if err == nil {
		return
	}

	// RFC3339Nano
	*t, err = time.Parse(time.RFC3339Nano, str)
	if err == nil {
		return
	}

	*t, err = time.Parse("2006-01-02 15:04:05", str)
	if err == nil {
		return
	}

	*t, err = time.Parse(time.ANSIC, str)
	if err == nil {
		return
	}

	*t, err = time.Parse(time.UnixDate, str)
	if err == nil {
		return
	}

	*t, err = time.Parse(time.RFC850, str)
	if err == nil {
		return
	}

	*t, err = time.Parse(time.RFC1123, str)
	if err == nil {
		return
	}

	return
}

func timeNextToken(iter *Iterator) byte {
	// a variation of skip whitespaces, returning the next non-whitespace token
	for {
		for i := iter.head; i < iter.tail; i++ {
			c := iter.buf[i]
			switch c {
			case '\n', '\t', '\r':
				continue
			}
			iter.head = i + 1
			return c
		}
		if !iter.loadMore() {
			return 0
		}
	}
}
