package jsoniter

import "fmt"

// ReadNil reads a json object as nil and
// returns whether it's a nil or not
func (iter *Iterator) ReadNil() (ret bool) {
	c := iter.nextToken()
	if c == 'n' {
		iter.skipThreeBytes('u', 'l', 'l') // null
		return true
	}
	iter.unreadByte()
	return false
}

// ReadBool reads a json object as Bool
func (iter *Iterator) ReadBool() (ret bool) {
	c := iter.nextToken()
	if c == 't' {
		iter.skipThreeBytes('r', 'u', 'e')
		return true
	}
	if c == 'f' {
		iter.skipFourBytes('a', 'l', 's', 'e')
		return false
	}
	iter.ReportError("ReadBool", "expect t or f")
	return
}

// SkipAndReturnBytes skip next JSON element, and return its content as []byte.
// The []byte can be kept, it is a copy of data.
func (iter *Iterator) SkipAndReturnBytes() []byte {
	iter.startCapture(iter.head)
	iter.Skip()
	return iter.stopCapture()
}

type captureBuffer struct {
	startedAt int
	captured  []byte
}

func (iter *Iterator) startCapture(captureStartedAt int) {
	if iter.captured != nil {
		panic("already in capture mode")
	}
	iter.captureStartedAt = captureStartedAt
	iter.captured = make([]byte, 0, 32)
}

func (iter *Iterator) stopCapture() []byte {
	if iter.captured == nil {
		panic("not in capture mode")
	}
	captured := iter.captured
	remaining := iter.buf[iter.captureStartedAt:iter.head]
	iter.captureStartedAt = -1
	iter.captured = nil
	if len(captured) == 0 {
		return remaining
	}
	captured = append(captured, remaining...)
	return captured
}

// Skip skips a json object and positions to relatively the next json object
func (iter *Iterator) Skip() {
	c := iter.nextToken()
	switch c {
	case '"':
		iter.skipString()
	case 'n':
		iter.skipThreeBytes('u', 'l', 'l') // null
	case 't':
		iter.skipThreeBytes('r', 'u', 'e') // true
	case 'f':
		iter.skipFourBytes('a', 'l', 's', 'e') // false
	case '-', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		iter.skipNumber()
	case '[':
		iter.skipArray()
	case '{':
		iter.skipObject()
	default:
		iter.ReportError("Skip", fmt.Sprintf("do not know how to skip: %v", c))
		return
	}
}

func (iter *Iterator) skipString() {
	for {
		end, escaped := iter.findStringEnd()
		if end == -1 {
			if !iter.loadMore() {
				iter.ReportError("skipString", "incomplete string")
				return
			}
			if escaped {
				iter.head = 1 // skip the first char as last char read is \
			}
		} else {
			iter.head = end
			return
		}
	}
}

// adapted from: https://github.com/buger/jsonparser/blob/master/parser.go
// Tries to find the end of string
// Support if string contains escaped quote symbols.
func (iter *Iterator) findStringEnd() (int, bool) {
	escaped := false
	for i := iter.head; i < iter.tail; i++ {
		c := iter.buf[i]
		if c == '"' {
			if !escaped {
				return i + 1, false
			}
			j := i - 1
			for {
				if j < iter.head || iter.buf[j] != '\\' {
					// even number of backslashes
					// either end of buffer, or " found
					return i + 1, true
				}
				j--
				if j < iter.head || iter.buf[j] != '\\' {
					// odd number of backslashes
					// it is \" or \\\"
					break
				}
				j--
			}
		} else if c == '\\' {
			escaped = true
		}
	}
	j := iter.tail - 1
	for {
		if j < iter.head || iter.buf[j] != '\\' {
			// even number of backslashes
			// either end of buffer, or " found
			return -1, false // do not end with \
		}
		j--
		if j < iter.head || iter.buf[j] != '\\' {
			// odd number of backslashes
			// it is \" or \\\"
			break
		}
		j--

	}
	return -1, true // end with \
}

func (iter *Iterator) skipArray() {
	level := 1
	for {
		for i := iter.head; i < iter.tail; i++ {
			switch iter.buf[i] {
			case '"': // If inside string, skip it
				iter.head = i + 1
				iter.skipString()
				i = iter.head - 1 // it will be i++ soon
			case '[': // If open symbol, increase level
				level++
			case ']': // If close symbol, increase level
				level--

				// If we have returned to the original level, we're done
				if level == 0 {
					iter.head = i + 1
					return
				}
			}
		}
		if !iter.loadMore() {
			iter.ReportError("skipObject", "incomplete array")
			return
		}
	}
}

func (iter *Iterator) skipObject() {
	level := 1
	for {
		for i := iter.head; i < iter.tail; i++ {
			switch iter.buf[i] {
			case '"': // If inside string, skip it
				iter.head = i + 1
				iter.skipString()
				i = iter.head - 1 // it will be i++ soon
			case '{': // If open symbol, increase level
				level++
			case '}': // If close symbol, increase level
				level--

				// If we have returned to the original level, we're done
				if level == 0 {
					iter.head = i + 1
					return
				}
			}
		}
		if !iter.loadMore() {
			iter.ReportError("skipObject", "incomplete object")
			return
		}
	}
}

func (iter *Iterator) skipNumber() {
	for {
		for i := iter.head; i < iter.tail; i++ {
			c := iter.buf[i]
			switch c {
			case ' ', '\n', '\r', '\t', ',', '}', ']':
				iter.head = i
				return
			}
		}
		if !iter.loadMore() {
			return
		}
	}
}

func (iter *Iterator) skipFourBytes(b1, b2, b3, b4 byte) {
	if iter.readByte() != b1 {
		iter.ReportError("skipFourBytes", fmt.Sprintf("expect %s", string([]byte{b1, b2, b3, b4})))
		return
	}
	if iter.readByte() != b2 {
		iter.ReportError("skipFourBytes", fmt.Sprintf("expect %s", string([]byte{b1, b2, b3, b4})))
		return
	}
	if iter.readByte() != b3 {
		iter.ReportError("skipFourBytes", fmt.Sprintf("expect %s", string([]byte{b1, b2, b3, b4})))
		return
	}
	if iter.readByte() != b4 {
		iter.ReportError("skipFourBytes", fmt.Sprintf("expect %s", string([]byte{b1, b2, b3, b4})))
		return
	}
}

func (iter *Iterator) skipThreeBytes(b1, b2, b3 byte) {
	if iter.readByte() != b1 {
		iter.ReportError("skipThreeBytes", fmt.Sprintf("expect %s", string([]byte{b1, b2, b3})))
		return
	}
	if iter.readByte() != b2 {
		iter.ReportError("skipThreeBytes", fmt.Sprintf("expect %s", string([]byte{b1, b2, b3})))
		return
	}
	if iter.readByte() != b3 {
		iter.ReportError("skipThreeBytes", fmt.Sprintf("expect %s", string([]byte{b1, b2, b3})))
		return
	}
}
