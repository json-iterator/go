package jsoniter

var (
	_ IKey = INT(0)
	_ IKey = STRING("")
)

// IKey represent a Key
type IKey interface {
	Type() int
	Equal(k IKey) bool
	Int() int
	String() string
}

// INT be the key
type INT int

// Type implement IKey
func (i INT) Type() int {
	return 0
}

// Equal implement IKey
func (i INT) Equal(k IKey) bool {
	if k.Type() == 0 && k.Int() == int(i) {
		return true
	}
	return false
}

// Int implement IKey
func (i INT) Int() int {
	return int(i)
}

func (i INT) String() string {
	panic("KeyInt cannot String()")
}

// STRING represent a object field
type STRING string

// Type implement IKey
func (k STRING) Type() int {
	return 1
}

// Equal implement IKey
func (k STRING) Equal(ik IKey) bool {
	if ik.Type() == 1 && ik.String() == string(k) {
		return true
	}
	return false
}

// Int implement IKey
func (k STRING) Int() int {
	panic("STRING cannot Int()")
}

func (k STRING) String() string {
	return string(k)
}
