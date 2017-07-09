package test

type Embedded1 struct {
	F1 int32
}

type Embedded2 struct {
	F1 int32
}

type typeForTest struct {
	Embedded1
	Embedded2
}
