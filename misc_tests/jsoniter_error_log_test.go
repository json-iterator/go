package misc_tests

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"testing"
)

// 7793382               152 ns/op              40 B/op          3 allocs/op
// 14880938               78 ns/op              23 B/op          1 allocs/op
// 10051323              119 ns/op              24 B/op          2 allocs/op
func Benchmark_fmt_errorf(b *testing.B) {
	b.ReportAllocs()
	for n := 0; n < b.N; n++ {
		fmt.Errorf("Error #%d", n)
	}
}

func Benchmark_errors_new_join(b *testing.B) {
	for n := 0; n < b.N; n++ {
		errors.New(strings.Join([]string{"Error #", strconv.Itoa(n)}, ""))
	}
}

func Benchmark_errors_new_sprintf(b *testing.B) {
	for n := 0; n < b.N; n++ {
		errors.New(fmt.Sprintf("Error #%d", n))
	}
}

func TestErrorsAreEquivalent(t *testing.T) {
	errorf := fmt.Errorf("Error #%d", 1).Error()
	join := errors.New(strings.Join([]string{"Error #", strconv.Itoa(1)}, "")).Error()
	sprintf := errors.New(fmt.Sprintf("Error #%d", 1)).Error()
	if errorf != join || errorf != sprintf {
		t.Fatalf("Errors are not equal. [errorf: %s] [join: %s] [sprintf: %s]", errorf, join, sprintf)
	}
}
