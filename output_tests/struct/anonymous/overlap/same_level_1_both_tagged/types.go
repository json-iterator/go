package test

type embedded1 struct {
	F1 int32 `json:"F1"`
}

type embedded2 struct {
	F1 int32 `json:"F1"`
}

type typeForTest struct {
	embedded1
	embedded2
}
