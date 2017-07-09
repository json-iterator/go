package test

type embedded1 struct {
	F1 int32
}

type embedded2 struct {
	F1 int32
}

type typeForTest struct {
	embedded1
	embedded2
}
