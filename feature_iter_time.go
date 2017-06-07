package jsoniter

import (
	"bytes"
	"fmt"
	"regexp"
	"time"
)

func (iter *Iterator) ReadTime() *time.Time {
	//  "2006-01-02T15:04:05"
	var (
		buf   bytes.Buffer
		bts   []byte
		t         = new(time.Time)
		times int = 0
	)
	buf.WriteString(`(\d{2}|\d{4})(?:\-)?([0]{1}\d{1}|[1]{1}[0-2]{1})(?:\-)?([0-2]{1}\d{1}|[3]{1}[0-1]{1})(T)?([0-1]{1}\d{1}|[2]{1}[0-3]{1})(?::)?([0-5]{1}\d{1})(?::)?([0-5]{1}\d{1})`)
	for c := iter.nextToken(); c != 0; c = iter.nextToken() {
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
	matched, err := regexp.MatchString(buf.String(), string(bts))
	if err != nil {
		fmt.Printf("time.Time parse format [%v], only support time format `2006-01-02T15:04:05`\n", string(bts))
		return nil
	} else if matched {
		*t, err = time.Parse("2006-01-02T15:04:05", string(bts))
		if err != nil {
			fmt.Printf("time parse failed, err=%s\n", err.Error())
			return nil
		}
		return t
	}
	return nil
}
