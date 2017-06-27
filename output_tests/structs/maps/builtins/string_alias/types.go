package test

type A1 string
type A2 map[string]A1

type T struct {
	F1 map[string]A1
	F2 A2
}
