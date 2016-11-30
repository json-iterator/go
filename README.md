# go
faster than DOM, more usable than SAX/StAX

# string

```
func Benchmark_jsoniter_string(b *testing.B) {
	for n := 0; n < b.N; n++ {
		lexer := NewLexerWithArray([]byte(`"\ud83d\udc4a"`))
		lexer.LexString()
	}
}
```

10000000	       140 ns/op

```
func Benchmark_json_string(b *testing.B) {
	for n := 0; n < b.N; n++ {
		result := ""
		json.Unmarshal([]byte(`"\ud83d\udc4a"`), &result)
	}
}
````

2000000	       710 ns/op (5x slower)

# int

```
func Benchmark_jsoniter_int(b *testing.B) {
	for n := 0; n < b.N; n++ {
		lexer := NewLexerWithArray([]byte(`-100`))
		lexer.LexInt64()
	}
}
```

30000000	        60.1 ns/op

```
func Benchmark_json_int(b *testing.B) {
	for n := 0; n < b.N; n++ {
		result := int64(0)
		json.Unmarshal([]byte(`-100`), &result)
	}
}
```

3000000	       505 ns/op (8x slower)