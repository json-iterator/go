package test

type DoubleEmbedded1 struct {
	F1 int32
}

type embedded1 struct {
	DoubleEmbedded1
}

type DoubleEmbedded2 struct {
	F1 int32
}

type embedded2 struct {
	DoubleEmbedded2
}

type typeForTest struct {
	embedded1
	embedded2
}
