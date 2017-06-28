package test

type Embedded map[string]string

type T struct {
	*Embedded
}
