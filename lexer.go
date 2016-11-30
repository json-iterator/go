package jsoniter

import (
	"io"
	"errors"
	"fmt"
	"unicode/utf16"
)

type Lexer struct {
	reader io.Reader
	buf    []byte
	head   int
	tail   int
}

func NewLexer(reader io.Reader, bufSize int) *Lexer {
	return &Lexer{
		reader: reader,
		buf: make([]byte, bufSize),
		head: 0,
		tail: 0,
	}
}

func NewLexerWithArray(input []byte) *Lexer {
	return &Lexer{
		reader: nil,
		buf: input,
		head: 0,
		tail: len(input),
	}
}

func (lexer *Lexer) readByte() (byte, error) {
	if lexer.head == lexer.tail {
		if lexer.reader == nil {
			return 0, io.EOF
		}
		n, err := lexer.reader.Read(lexer.buf)
		if err != nil {
			return 0, err
		}
		if n == 0 {
			return 0, io.EOF
		}
		lexer.head = 0
		lexer.tail = n
	}
	b := lexer.buf[lexer.head]
	lexer.head += 1
	return b, nil
}

func (lexer *Lexer) unreadByte() error {
	if lexer.head == 0 {
		return errors.New("unread too many bytes")
	}
	lexer.head -= 1
	return nil
}

const maxUint64 = (1 << 64 - 1)
const cutoffUint64 = maxUint64 / 10 + 1
const maxUint32 = (1 << 32 - 1)
const cutoffUint32 = maxUint32 / 10 + 1

func (lexer *Lexer) LexUin64() (uint64, error) {
	var n uint64
	c, err := lexer.readByte()
	if err != nil {
		return 0, err
	}

	/* a single zero, or a series of integers */
	if c == '0' {
		c, err = lexer.readByte()
		if err != nil && err != io.EOF {
			return 0, err
		}
	} else if c >= '1' && c <= '9' {
		for c >= '0' && c <= '9' {
			var v byte
			v = c - '0'
			if n >= cutoffUint64 {
				return 0, errors.New("overflow")
			}
			n = n * uint64(10) + uint64(v)
			c, err = lexer.readByte()
			if err != nil && err != io.EOF {
				return 0, err
			}
		}
		lexer.unreadByte()
	} else {
		lexer.unreadByte()
		return 0, errors.New("unexpected")
	}
	return n, nil
}

func (lexer *Lexer) LexInt64() (int64, error) {
	c, err := lexer.readByte()
	if err != nil {
		return 0, err
	}

	/* optional leading minus */
	if c == '-' {
		n, err := lexer.LexUin64()
		if err != nil {
			return 0, err
		}
		return -int64(n), nil
	} else {
		lexer.unreadByte()
		n, err := lexer.LexUin64()
		if err != nil {
			return 0, err
		}
		return int64(n), nil
	}
}

func (lexer *Lexer) LexString() (string, error) {
	str := make([]byte, 0, 10)
	c, err := lexer.readByte()
	if err != nil {
		return "", err
	}
	if c != '"' {
		return "", errors.New("unexpected")
	}
	for {
		c, err = lexer.readByte()
		if err != nil {
			return "", err
		}
		switch c {
		case '\\':
			c, err = lexer.readByte()
			if err != nil {
				return "", err
			}
			switch c {
			case 'u':
				r, err := lexer.readU4()
				if err != nil {
					return "", err
				}
				if utf16.IsSurrogate(r) {
					c, err = lexer.readByte()
					if err != nil {
						return "", err
					}
					if c != '\\' {
						return "", fmt.Errorf("unexpected: %v", c)
					}
					c, err = lexer.readByte()
					if err != nil {
						return "", err
					}
					if c != 'u' {
						return "", fmt.Errorf("unexpected: %v", c)
					}
					r2, err := lexer.readU4()
					if err != nil {
						return "", err
					}
					combined := utf16.DecodeRune(r, r2)
					str = appendRune(str, combined)
				} else {
					str = appendRune(str, r)
				}
			case '"':
				str = append(str, '"')
			case '\\':
				str = append(str, '\\')
			case '/':
				str = append(str, '/')
			case 'b':
				str = append(str, '\b')
			case 'f':
				str = append(str, '\f')
			case 'n':
				str = append(str, '\n')
			case 'r':
				str = append(str, '\r')
			case 't':
				str = append(str, '\t')
			default:
				return "", errors.New("unexpected")
			}
		case '"':
			return string(str), nil
		default:
			str = append(str, c)
		}
	}
}

func (lexer *Lexer) readU4() (rune, error) {
	var u4 rune
	for i := 0; i < 4; i++ {
		c, err := lexer.readByte()
		if err != nil {
			return 0, err
		}
		if (c >= '0' && c <= '9') {
			if u4 >= cutoffUint32 {
				return 0, errors.New("overflow")
			}
			u4 = u4 * 16 + rune(c - '0')
		} else if ((c >= 'a' && c <= 'f') ) {
			if u4 >= cutoffUint32 {
				return 0, errors.New("overflow")
			}
			u4 = u4 * 16 + rune(c - 'a' + 10)
		} else {
			return 0, fmt.Errorf("unexpected: %v", c)
		}
	}
	return u4, nil
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

	rune1Max = 1 << 7 - 1
	rune2Max = 1 << 11 - 1
	rune3Max = 1 << 16 - 1

	surrogateMin = 0xD800
	surrogateMax = 0xDFFF

	MaxRune = '\U0010FFFF' // Maximum valid Unicode code point.
	RuneError = '\uFFFD'     // the "error" Rune or "Unicode replacement character"
)

func appendRune(p []byte, r rune) []byte {
	// Negative values are erroneous. Making it unsigned addresses the problem.
	switch i := uint32(r); {
	case i <= rune1Max:
		p = append(p, byte(r))
		return p
	case i <= rune2Max:
		p = append(p, t2 | byte(r >> 6))
		p = append(p, tx | byte(r) & maskx)
		return p
	case i > MaxRune, surrogateMin <= i && i <= surrogateMax:
		r = RuneError
		fallthrough
	case i <= rune3Max:
		p = append(p, t3 | byte(r >> 12))
		p = append(p, tx | byte(r >> 6) & maskx)
		p = append(p, tx | byte(r) & maskx)
		return p
	default:
		p = append(p, t4 | byte(r >> 18))
		p = append(p, tx | byte(r >> 12) & maskx)
		p = append(p, tx | byte(r >> 6) & maskx)
		p = append(p, tx | byte(r) & maskx)
		return p
	}
}

