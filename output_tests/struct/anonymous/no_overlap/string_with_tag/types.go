package test

type Embedded string

type T struct {
	Embedded `json:"othername"`
}
