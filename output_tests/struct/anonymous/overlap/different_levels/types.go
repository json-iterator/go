package test

type Embedded struct {
	F1 int32
}

type T struct {
	F1 string
	Embedded
}
