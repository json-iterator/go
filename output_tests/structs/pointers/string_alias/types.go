package test

type A1 string
type A2 *string

type T struct {
	F1 *A1
	F2 A2
	F3 *A2
}
