package test

type A1 float32
type A2 *float32

type T struct {
	F1 *A1
	F2 A2
	F3 *A2
}
