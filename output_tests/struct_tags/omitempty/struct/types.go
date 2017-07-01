package test

type T struct {
	F struct{} `json:"f,omitempty"` // omitempty is meaningless here
}
