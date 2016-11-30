package jsoniter

import (
	fflib "github.com/pquerna/ffjson/fflib/v1"
	"strconv"
	"fmt"
)

type Token fflib.FFTok

const TokenInteger = Token(fflib.FFTok_integer)
const TokenDouble = Token(fflib.FFTok_double)
const TokenBool = Token(fflib.FFTok_bool)
const TokenError = Token(fflib.FFTok_error)
const TokenLeftBrace = Token(fflib.FFTok_left_brace)
const TokenLeftBracket = Token(fflib.FFTok_left_bracket)
const TokenRightBrace = Token(fflib.FFTok_right_brace)
const TokenRightBracket = Token(fflib.FFTok_right_bracket)
const TokenComma = Token(fflib.FFTok_comma)
const TokenString = Token(fflib.FFTok_string)
const TokenColon = Token(fflib.FFTok_colon)

func (tok Token) ToString() string {
	return fmt.Sprintf("%v", fflib.FFTok(tok))
}

type Iterator struct {
	ErrorHandler func(error)
	lexer        *fflib.FFLexer
}

type UnexpectedToken struct {
	Expected Token
	Actual   Token
}

func (err *UnexpectedToken) Error() string {
	return fmt.Sprintf("unexpected token, expected %v, actual %v", fflib.FFTok(err.Expected), fflib.FFTok(err.Actual))
}

func NewIterator(input []byte) Iterator {
	lexer := fflib.NewFFLexer(input)
	return Iterator{
		lexer: lexer,
	}
}

func (iter Iterator) ReadArray(callback func(Iterator, int)) {
	lexer := iter.lexer
	if !iter.AssertToken(TokenLeftBrace) {
		return
	}
	index := 0
	for {
		lexer.Scan()
		callback(iter, index)
		index += 1
		if iter.Token() != TokenComma {
			break
		}
	}
	iter.AssertToken(TokenRightBrace)
	lexer.Scan()
}

func (iter Iterator) ReadObject(callback func(Iterator, string)) {
	lexer := iter.lexer
	if !iter.AssertToken(TokenLeftBracket) {
		return
	}
	for {
		lexer.Scan()
		field := iter.ReadString()
		iter.AssertToken(TokenColon)
		lexer.Scan()
		callback(iter, field)
		if iter.Token() != TokenComma {
			break
		}
	}
	iter.AssertToken(TokenRightBracket)
	lexer.Scan()
}

func (iter Iterator) Skip() {
	switch iter.Token() {
	case TokenLeftBracket:
		iter.ReadObject(func(iter Iterator, field string) {
			iter.Skip()
		})
	case TokenLeftBrace:
		iter.ReadArray(func(iter Iterator, index int) {
			iter.Skip()
		})
	default:
		iter.lexer.Scan()
	}
}

func (iter Iterator) ReadInt8() (rval int8) {
	lexer := iter.lexer
	if !iter.AssertToken(TokenInteger) {
		return
	}
	field, err := lexer.CaptureField(lexer.Token)
	if err != nil {
		iter.OnError(err)
		return
	}
	lexer.Scan()
	number, err := strconv.ParseInt(string(field), 10, 8)
	if err != nil {
		iter.OnError(fmt.Errorf("failed to convert %v: %v", string(field), err.Error()))
		return
	}
	return int8(number)
}

func (iter Iterator) ReadInt16() (rval int16) {
	lexer := iter.lexer
	if !iter.AssertToken(TokenInteger) {
		return
	}
	field, err := lexer.CaptureField(lexer.Token)
	if err != nil {
		iter.OnError(err)
		return
	}
	lexer.Scan()
	number, err := strconv.ParseInt(string(field), 10, 16)
	if err != nil {
		iter.OnError(fmt.Errorf("failed to convert %v: %v", string(field), err.Error()))
		return
	}
	return int16(number)
}

func (iter Iterator) ReadInt32() (rval int32) {
	lexer := iter.lexer
	if !iter.AssertToken(TokenInteger) {
		return
	}
	field, err := lexer.CaptureField(lexer.Token)
	if err != nil {
		iter.OnError(err)
		return
	}
	lexer.Scan()
	number, err := strconv.ParseInt(string(field), 10, 32)
	if err != nil {
		iter.OnError(fmt.Errorf("failed to convert %v: %v", string(field), err.Error()))
		return
	}
	return int32(number)
}

func (iter Iterator) ReadInt64() (rval int64) {
	lexer := iter.lexer
	if !iter.AssertToken(TokenInteger) {
		return
	}
	field, err := lexer.CaptureField(lexer.Token)
	if err != nil {
		iter.OnError(err)
		return
	}
	lexer.Scan()
	number, err := strconv.ParseInt(string(field), 10, 64)
	if err != nil {
		iter.OnError(fmt.Errorf("failed to convert %v: %v", string(field), err.Error()))
		return
	}
	return number
}

func (iter Iterator) ReadInt() (rval int) {
	lexer := iter.lexer
	n, err := lexer.LexInt64()
	if err != nil {
		iter.OnError(err)
		return
	}
	return int(n)
}

func (iter Iterator) ReadUint() (rval uint) {
	lexer := iter.lexer
	if !iter.AssertToken(TokenInteger) {
		return
	}
	field, err := lexer.CaptureField(lexer.Token)
	if err != nil {
		iter.OnError(err)
		return
	}
	lexer.Scan()
	number, err := strconv.ParseUint(string(field), 10, 64)
	if err != nil {
		iter.OnError(fmt.Errorf("failed to convert %v: %v", string(field), err.Error()))
		return
	}
	return uint(number)
}

func (iter Iterator) ReadUint8() (rval uint8) {
	lexer := iter.lexer
	if !iter.AssertToken(TokenInteger) {
		return
	}
	field, err := lexer.CaptureField(lexer.Token)
	if err != nil {
		iter.OnError(err)
		return
	}
	lexer.Scan()
	number, err := strconv.ParseUint(string(field), 10, 8)
	if err != nil {
		iter.OnError(fmt.Errorf("failed to convert %v: %v", string(field), err.Error()))
		return
	}
	return uint8(number)
}

func (iter Iterator) ReadUint16() (rval uint16) {
	lexer := iter.lexer
	if !iter.AssertToken(TokenInteger) {
		return
	}
	field, err := lexer.CaptureField(lexer.Token)
	if err != nil {
		iter.OnError(err)
		return
	}
	lexer.Scan()
	number, err := strconv.ParseUint(string(field), 10, 16)
	if err != nil {
		iter.OnError(fmt.Errorf("failed to convert %v: %v", string(field), err.Error()))
		return
	}
	return uint16(number)
}

func (iter Iterator) ReadUint32() (rval uint32) {
	lexer := iter.lexer
	if !iter.AssertToken(TokenInteger) {
		return
	}
	field, err := lexer.CaptureField(lexer.Token)
	if err != nil {
		iter.OnError(err)
		return
	}
	lexer.Scan()
	number, err := strconv.ParseUint(string(field), 10, 32)
	if err != nil {
		iter.OnError(fmt.Errorf("failed to convert %v: %v", string(field), err.Error()))
		return
	}
	return uint32(number)
}

func (iter Iterator) ReadUint64() (rval uint64) {
	lexer := iter.lexer
	if !iter.AssertToken(TokenInteger) {
		return
	}
	field, err := lexer.CaptureField(lexer.Token)
	if err != nil {
		iter.OnError(err)
		return
	}
	lexer.Scan()
	number, err := strconv.ParseUint(string(field), 10, 64)
	if err != nil {
		iter.OnError(fmt.Errorf("failed to convert %v: %v", string(field), err.Error()))
		return
	}
	return uint64(number)
}

func (iter Iterator) ReadFloat32() (rval float32) {
	lexer := iter.lexer
	if !iter.AssertToken(TokenDouble) {
		return
	}
	field, err := lexer.CaptureField(lexer.Token)
	if err != nil {
		iter.OnError(err)
		return
	}
	lexer.Scan()
	number, err := strconv.ParseFloat(string(field), 32)
	if err != nil {
		iter.OnError(fmt.Errorf("failed to convert %v: %v", string(field), err.Error()))
		return
	}
	return float32(number)
}

func (iter Iterator) ReadFloat64() (rval float64) {
	lexer := iter.lexer
	if !iter.AssertToken(TokenDouble) {
		return
	}
	field, err := lexer.CaptureField(lexer.Token)
	if err != nil {
		iter.OnError(err)
		return
	}
	lexer.Scan()
	number, err := strconv.ParseFloat(string(field), 64)
	if err != nil {
		iter.OnError(fmt.Errorf("failed to convert %v: %v", string(field), err.Error()))
		return
	}
	return float64(number)
}

func (iter Iterator) ReadString() (rval string) {
	lexer := iter.lexer
	if !iter.AssertToken(TokenString) {
		return
	}
	field, err := lexer.CaptureField(lexer.Token)
	if err != nil {
		iter.OnError(err)
		return
	}
	lexer.Scan()
	return string(field[1:len(field) - 1])
}

func (iter Iterator) AssertToken(expected Token) bool {
	actual := iter.Token()
	if expected != actual {
		if actual == TokenError {
			fmt.Println(iter.lexer.BigError.Error())
		}
		iter.OnError(&UnexpectedToken{
			Expected: expected,
			Actual: actual,
		})
		return false
	}
	return true
}

func (iter Iterator) Token() Token {
	return Token(iter.lexer.Token)
}

func (iter *Iterator) OnError(err error) {
	if iter.ErrorHandler == nil {
		panic(err.Error())
	} else {
		iter.ErrorHandler(err)
	}
}