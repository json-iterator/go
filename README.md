# json iterator (jsoniter)

faster than DOM, more usable than SAX/StAX

for performance numbers, see https://github.com/json-iterator/go-benchmark

# DOM style api

Jsoniter can work as drop in replacement for json.Unmarshal

```
type StructOfTag struct {
	field1 string `json:"field-1"`
	field2 string `json:"-"`
	field3 int `json:",string"`
}

func Test_reflect_struct_tag_field(t *testing.T) {
	err := jsoniter.Unmarshal(`{"field-1": "hello", "field2": "", "field3": "100"}`, &struct_)
	if struct_.field1 != "hello" {
		fmt.Println(err)
		t.Fatal(struct_.field1)
	}
	if struct_.field2 != "world" {
		fmt.Println(err)
		t.Fatal(struct_.field2)
	}
	if struct_.field3 != 100 {
		fmt.Println(err)
		t.Fatal(struct_.field3)
	}
}
```

# StAX style api

When you need the maximum performance, the pull style api allows you to control every bit of parsing process. You
can bind value to object without reflection, or you can calculate the sum of array on the fly without intermediate objects.

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

# Customization

Of course, you can use the low level pull api to do anything you like. But most of the time,
reflection based api is fast enough. How to control the parsing process when we are using the reflection api?
json.Unmarshaller is not flexible enough. Jsoniter provides much better customizability.

```
func Test_customize_type_decoder(t *testing.T) {
	RegisterTypeDecoder("time.Time", func(ptr unsafe.Pointer, iter *Iterator) {
		t, err := time.ParseInLocation("2006-01-02 15:04:05", iter.ReadString(), time.UTC)
		if err != nil {
			iter.Error = err
			return
		}
		*((*time.Time)(ptr)) = t
	})
	defer ClearDecoders()
	val := time.Time{}
	err := Unmarshal([]byte(`"2016-12-05 08:43:28"`), &val)
	if err != nil {
		t.Fatal(err)
	}
	year, month, day := val.Date()
	if year != 2016 || month != 12 || day != 5 {
		t.Fatal(val)
	}
}
```

there is no way to add json.Unmarshaller to time.Time as the type is not defined by you. Using jsoniter, we can.

```
type Tom struct {
	field1 string
}

func Test_customize_field_decoder(t *testing.T) {
	RegisterFieldDecoder("jsoniter.Tom", "field1", func(ptr unsafe.Pointer, iter *Iterator) {
		*((*string)(ptr)) = strconv.Itoa(iter.ReadInt())
	})
	defer ClearDecoders()
	tom := Tom{}
	err := Unmarshal([]byte(`{"field1": 100}`), &tom)
	if err != nil {
		t.Fatal(err)
	}
}
```

It is very common the input json has certain fields massed up. We want string, but it is int, etc. The old way is to
define a struct of exact type like the json. Then we convert from one struct to a new struct. It is just too much work.
Using jsoniter you can tweak the field conversion.

