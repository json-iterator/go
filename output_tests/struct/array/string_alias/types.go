package test

type A1 string
type A2 [4]A1

type typeForTest struct {
	F1 [4]A1
	F2 A2
}
