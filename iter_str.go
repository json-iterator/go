package jsoniter

import (
	"fmt"
	"strings"
	"unicode/utf16"
)

// ReadString read string from iterator
func (iter *Iterator) ReadString() (ret string) {
	c := iter.nextToken()
	if c == '"' {
		for i := iter.head; i < iter.tail; i++ {
			c := iter.buf[i]
			if c == '"' {
				ret = iter.strf.NewString(iter.buf[iter.head:i])
				iter.head = i + 1
				return ret
			} else if c == '\\' {
				break
			} else if c < ' ' {
				iter.ReportError("ReadString",
					fmt.Sprintf(`invalid control character found: %d`, c))
				return
			}
		}
		return iter.readStringSlowPath()
	} else if c == 'n' {
		iter.skipThreeBytes('u', 'l', 'l')
		return ""
	}
	iter.ReportError("ReadString", `expects " or n, but found `+iter.strf.NewString([]byte{c}))
	return
}

func (iter *Iterator) readStringSlowPath() string {
	strb := iter.cfg.borrowStringBuilder()
	defer iter.cfg.returnStringBuilder(strb)
	var c byte
	for iter.Error == nil {
		c = iter.readByte()
		if c == '"' {
			return strb.String()
		}
		if c == '\\' {
			c = iter.readByte()
			iter.readEscapedChar(c, strb)
		} else {
			strb.WriteByte(c)
		}
	}
	iter.ReportError("readStringSlowPath", "unexpected end of input")
	return ""
}

func (iter *Iterator) readEscapedChar(c byte, strb *strings.Builder) {
	switch c {
	case 'u':
		r := iter.readU4()
		if utf16.IsSurrogate(r) {
			c = iter.readByte()
			if iter.Error != nil {
				strb.Reset()
				return
			}
			if c != '\\' {
				iter.unreadByte()
				appendRune(strb, r)
				return
			}
			c = iter.readByte()
			if iter.Error != nil {
				strb.Reset()
				return
			}
			if c != 'u' {
				appendRune(strb, r)
				iter.readEscapedChar(c, strb)
				return
			}
			r2 := iter.readU4()
			if iter.Error != nil {
				strb.Reset()
				return
			}
			combined := utf16.DecodeRune(r, r2)
			if combined == '\uFFFD' {
				appendRune(strb, r)
				appendRune(strb, r2)
			} else {
				appendRune(strb, combined)
			}
		} else {
			appendRune(strb, r)
		}
	case '"':
		strb.WriteByte('"')
	case '\\':
		strb.WriteByte('\\')
	case '/':
		strb.WriteByte('/')
	case 'b':
		strb.WriteByte('\b')
	case 'f':
		strb.WriteByte('\f')
	case 'n':
		strb.WriteByte('\n')
	case 'r':
		strb.WriteByte('\r')
	case 't':
		strb.WriteByte('\t')
	default:
		iter.ReportError("readEscapedChar",
			`invalid escape char after \`)
		strb.Reset()
		return
	}
	return
}

// ReadStringAsSlice read string from iterator without copying into string form.
// The []byte can not be kept, as it will change after next iterator call.
func (iter *Iterator) ReadStringAsSlice() (ret []byte) {
	c := iter.nextToken()
	if c == '"' {
		for i := iter.head; i < iter.tail; i++ {
			// require ascii string and no escape
			// for: field name, base64, number
			if iter.buf[i] == '"' {
				// fast path: reuse the underlying buffer
				ret = iter.buf[iter.head:i]
				iter.head = i + 1
				return ret
			}
		}
		readLen := iter.tail - iter.head
		copied := make([]byte, readLen, readLen*2)
		copy(copied, iter.buf[iter.head:iter.tail])
		iter.head = iter.tail
		for iter.Error == nil {
			c := iter.readByte()
			if c == '"' {
				return copied
			}
			copied = append(copied, c)
		}
		return copied
	}
	iter.ReportError("ReadStringAsSlice", `expects " or n, but found `+string([]byte{c}))
	return
}

func (iter *Iterator) readU4() (ret rune) {
	for i := 0; i < 4; i++ {
		c := iter.readByte()
		if iter.Error != nil {
			return
		}
		if c >= '0' && c <= '9' {
			ret = ret*16 + rune(c-'0')
		} else if c >= 'a' && c <= 'f' {
			ret = ret*16 + rune(c-'a'+10)
		} else if c >= 'A' && c <= 'F' {
			ret = ret*16 + rune(c-'A'+10)
		} else {
			iter.ReportError("readU4", "expects 0~9 or a~f, but found "+string([]byte{c}))
			return
		}
	}
	return ret
}

const (
	t1 = 0x00 // 0000 0000
	tx = 0x80 // 1000 0000
	t2 = 0xC0 // 1100 0000
	t3 = 0xE0 // 1110 0000
	t4 = 0xF0 // 1111 0000
	t5 = 0xF8 // 1111 1000

	maskx = 0x3F // 0011 1111
	mask2 = 0x1F // 0001 1111
	mask3 = 0x0F // 0000 1111
	mask4 = 0x07 // 0000 0111

	rune1Max = 1<<7 - 1
	rune2Max = 1<<11 - 1
	rune3Max = 1<<16 - 1

	surrogateMin = 0xD800
	surrogateMax = 0xDFFF

	maxRune   = '\U0010FFFF' // Maximum valid Unicode code point.
	runeError = '\uFFFD'     // the "error" Rune or "Unicode replacement character"
)

func appendRune(p *strings.Builder, r rune) {
	// Negative values are erroneous. Making it unsigned addresses the problem.
	switch i := uint32(r); {
	case i <= rune1Max:
		p.WriteByte(byte(r))
	case i <= rune2Max:
		p.WriteByte(t2 | byte(r>>6))
		p.WriteByte(tx | byte(r)&maskx)
	case i > maxRune, surrogateMin <= i && i <= surrogateMax:
		r = runeError
		fallthrough
	case i <= rune3Max:
		p.WriteByte(t3 | byte(r>>12))
		p.WriteByte(tx | byte(r>>6)&maskx)
		p.WriteByte(tx | byte(r)&maskx)
	default:
		p.WriteByte(t4 | byte(r>>18))
		p.WriteByte(tx | byte(r>>12)&maskx)
		p.WriteByte(tx | byte(r>>6)&maskx)
		p.WriteByte(tx | byte(r)&maskx)
	}
}
