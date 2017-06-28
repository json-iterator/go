package test

type DoubleEmbedded1 struct {
	F1 int32
}

type Embedded1 struct {
	DoubleEmbedded1
}

type DoubleEmbedded2 struct {
	F1 int32 `json:"F1"`
}

type Embedded2 struct {
	DoubleEmbedded2
}

type T struct {
	Embedded1
	Embedded2
}
