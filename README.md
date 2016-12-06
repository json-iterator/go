# json iterator (jsoniter)

faster than DOM, more usable than SAX/StAX

# Why json iterator?

## 1. It is faster

jsoniter can work as drop in replacement for json.Unmarshal, with or without reflection. Unlike https://github.com/pquerna/ffjson
jsoniter does not require `go generate`

for performance numbers, see https://github.com/json-iterator/go-benchmark

## 2. io.Reader as input

jsoniter does not read the whole json into memory, it parse the document in a streaming way. Unlike https://github.com/pquerna/ffjson
it requires []byte as input.

## 3. Pull style api

jsoniter can be used just like json.Unmarshal, for example

```
type StructOfTag struct {
    field1 string `json:"field-1"`
    field2 string `json:"-"`
    field3 int `json:",string"`
}

struct_ := StructOfTag{}
jsoniter.Unmarshal(`{"field-1": "hello", "field2": "", "field3": "100"}`, &struct_)
```

But it allows you to go down one level lower, to control the parsing process using pull style api (like StAX, if you
know what I mean). Here is just a demo of what you can do

```
iter := jsoniter.ParseString(`[1,2,3]`)
for iter.ReadArray() {
  iter.ReadUint64()
}
```

## 4. Customization

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

there is no way to add json.Unmarshaller to time.Time as the type is not defined by you (type alias time.Time is not fun to use).
Using jsoniter, we can.

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

# Why not json iterator?

jsoniter does not plan to support `map[string]interface{}`, period.
