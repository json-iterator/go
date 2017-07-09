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
	F1 S1 `json:"F1"`
	F2 S2 `json:"f2"`
	F3 S3 `json:"-"`
	F4 S4 `json:"-,"`
	F5 S5 `json:","`
	F6 S6 `json:""`
}
