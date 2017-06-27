package test

type DoubleEmbedded struct {
	F1 int32 `json:"F1"`
}

type Embedded1 struct {
	DoubleEmbedded
	F1 int32
}

type Embedded2 struct {
	F1 int32 `json:"F1"`
	DoubleEmbedded
}

type T struct {
	Embedded1
	Embedded2
}
