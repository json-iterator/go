package test

type Embedded string

type typeForTest struct {
	Embedded `json:"othername"`
}
