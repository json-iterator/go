package test

type Embedded map[string]string

type typeForTest struct {
	*Embedded
}
