package test

type Embedded1 struct {
	F1 int32
}

type Embedded2 struct {
	F1 int32
}

type T struct {
	Embedded1
	Embedded2
}
