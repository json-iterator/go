# json iterator (jsoniter)

faster than DOM, more usable than SAX/StAX

for performance numbers, see https://github.com/json-iterator/go-benchmark

# DOM style api

```
type StructOfTag struct {
	field1 string `json:"field-1"`
	field2 string `json:"-"`
	field3 int `json:",string"`
}

func Test_reflect_struct_tag_field(t *testing.T) {
	jsoniter.Unmarshal(`{"field-1": "hello", "field2": "", "field3": "100"}`, &struct_)
	if struct_.field1 != "hello" {
		fmt.Println(iter.Error)
		t.Fatal(struct_.field1)
	}
	if struct_.field2 != "world" {
		fmt.Println(iter.Error)
		t.Fatal(struct_.field2)
	}
	if struct_.field3 != 100 {
		fmt.Println(iter.Error)
		t.Fatal(struct_.field3)
	}
}
```

# StAX style api

Array

```
iter := jsoniter.ParseString(`[1,2,3]`)
for iter.ReadArray() {
  iter.ReadUint64()
}
```

Object

```
type TestObj struct {
    Field1 string
    Field2 uint64
}
```

```
iter := jsoniter.ParseString(`{"field1": "1", "field2": 2}`)
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
```

Skip

```
iter := jsoniter.ParseString(`[ {"a" : [{"b": "c"}], "d": 102 }, "b"]`)
iter.ReadArray()
iter.Skip()
iter.ReadArray()
if iter.ReadString() != "b" {
    t.FailNow()
}
```