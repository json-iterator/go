package jsoniter

var (
	_ IKey = Index(0)
	_ IKey = Field("")
)

// IKey represent a Key
type IKey interface {
	Type() int
	Equal(k IKey) bool
	Int() int
	String() string
}

// Index be the key
type Index int

// Type implement IKey
func (i Index) Type() int {
	return 0
}

// Equal implement IKey
func (i Index) Equal(k IKey) bool {
	if k.Type() == 0 && k.Int() == int(i) {
		return true
	}
	return false
}

// Int implement IKey
func (i Index) Int() int {
	return int(i)
}

func (i Index) String() string {
	panic("KeyInt cannot String()")
}

// Field represent a object field
type Field string

// Type implement IKey
func (k Field) Type() int {
	return 1
}

// Equal implement IKey
func (k Field) Equal(ik IKey) bool {
	if ik.Type() == 1 && ik.String() == string(k) {
		return true
	}
	return false
}

// Int implement IKey
func (k Field) Int() int {
	panic("STRING cannot Int()")
}

func (k Field) String() string {
	return string(k)
}
