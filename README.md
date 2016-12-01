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

# array

```
func Benchmark_jsoniter_array(b *testing.B) {
	for n := 0; n < b.N; n++ {
		iter := ParseString(`[1,2,3]`)
		for iter.ReadArray() {
			iter.ReadUint64()
		}
	}
}
```

10000000	       189 ns/op

```
func Benchmark_json_array(b *testing.B) {
	for n := 0; n < b.N; n++ {
		result := []interface{}{}
		json.Unmarshal([]byte(`[1,2,3]`), &result)
	}
}
```
1000000	      1327 ns/op

# object

```
func Benchmark_jsoniter_object(b *testing.B) {
	for n := 0; n < b.N; n++ {
		iter := ParseString(`{"field1": "1", "field2": 2}`)
		obj := TestObj{}
		for field := iter.ReadObject(); field != ""; field = iter.ReadObject() {
			switch field {
			case "field1":
				obj.Field1 = iter.ReadString()
			case "field2":
				obj.Field2 = iter.ReadUint64()
			default:
				iter.ReportError("bind object", "unexpected field")
			}
		}
	}
}
```

5000000	       401 ns/op

```
func Benchmark_json_object(b *testing.B) {
	for n := 0; n < b.N; n++ {
		result := TestObj{}
		json.Unmarshal([]byte(`{"field1": "1", "field2": 2}`), &result)
	}
}
```

1000000	      1318 ns/op