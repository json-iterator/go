package test

type doubleEmbedded struct {
	F1 int32 `json:"F1"`
}

type embedded1 struct {
	doubleEmbedded
	F1 int32
}

type embedded2 struct {
	F1 int32 `json:"F1"`
	doubleEmbedded
}

type typeForTest struct {
	embedded1
	embedded2
}
