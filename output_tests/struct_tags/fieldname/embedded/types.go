package test

type S1 struct {
	S1F string
}
type S2 struct {
	S2F string
}
type S3 struct {
	S3F string
}
type S4 struct {
	S4F string
}
type S5 struct {
	S5F string
}
type S6 struct {
	S6F string
}

type typeForTest struct {
	S1 `json:"F1"`
	S2 `json:"f2"`
	S3 `json:"-"`
	S4 `json:"-,"`
	S5 `json:","`
	S6 `json:""`
}
