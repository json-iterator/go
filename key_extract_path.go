package jsoniter


const (
	// TypeObject value will be a object
	TypeObject int = 0

	// TypeArray value will be a array
	TypeArray int = 1

	// TypeArrayIndex array[index]
	TypeArrayIndex int = 2
)

// Node represent a Node in a Path
type Node interface {
	String() string
	Type() int
	Index() int
}

var (
	_ Node = ObjectNode("")
	_ Node = ArrayNode("")
	_ Node = ArrayIndex(1)
)

// ObjectNode represent a node in the paths
type ObjectNode string

func (n ObjectNode) String() string {
	return string(n)
}

// Type implement Node
func (n ObjectNode) Type() int {
	return TypeObject
}

// Index implement Node
func (n ObjectNode) Index() int {
	panic(n + "can not index")
}

// ArrayNode represent followed by a array
type ArrayNode string

func (n ArrayNode) String() string {
	return string(n)
}

// Type implement Node
func (n ArrayNode) Type() int {
	return TypeArray
}

// Index implement Node
func (n ArrayNode) Index() int {
	panic(n + "can not index")
}

// ArrayIndex represent a index of array
type ArrayIndex int

func (i ArrayIndex) String() string {
	panic("ArrayIndex can not String")
}

// Type implement Node
func (i ArrayIndex) Type() int {
	return TypeArrayIndex
}

// Index implement Node
func (i ArrayIndex) Index() int {
	return int(i)
}

// Path represent a path to extract
// var input2 = `{"a":["b",{"c","d"}]}`
// c:
// NewPath(ArrayNode("a"), ArrayIndex(1))
type Path []Node

// NewPath contruct a path
func NewPath(nodes ...Node) Path {
	return nodes
}
