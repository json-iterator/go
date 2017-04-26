package jsoniter

const (
	// TypeObject value will be a object
	TypeObject int = 0

	// TypeArray value will be a array
	TypeArray int = 1
)

// Path represent a path to extract
// var input2 = `{"a":["b",{"c","d"}]}`
// c:
// NewPath(Field("a"), Index(1))
type Path []IKey

// NewPath contruct a path
func NewPath(keys ...IKey) Path {
	return keys
}
