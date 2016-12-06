# json iterator (jsoniter)

faster than DOM, more usable than SAX/StAX

# Why json iterator?

## 1. It is faster

jsoniter can work as drop in replacement for json.Unmarshal, reflection-api is not only supported, but recommended.

for performance numbers, see https://github.com/json-iterator/go-benchmark

The reflection-api is very fast, on the same scale of hand written ones.

## 2. io.Reader as input

jsoniter does not read the whole json into memory, it parse the document in a streaming way.
There are too many json parser only take []byte as input, this one does not require so.

## 3. Pull style api

jsoniter can be used like drop-in replacement of json.Unmarshal, for example

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

## 5. Minimum work to parse, use whatever fits the job

I invented this wheel because I find it is tedious to parse json which does not match the object model you want to use.
Parse to `map[string]interface{}` is not only ugly but also slow. Parse to struct is not flexible enough to fix
some field type mismatch or structure mismatch.

If use low level tokenizer/lexer to work at the token level, it is too much work, not to mention there is very few parser
out there allow you to work on this level.

jsoniter pull-api is designed to be easy to use, so that you can map your data structure directly to parsing code.
It is still tedious I am not going to lie to you, but easier than pure tokenizer.
The real power is, you can mix the pull-api with reflection-api.
For example:

```
\\ given [1, {"a": "b"}]
iter.ReadArray()
iter.ReadInt()
iter.ReadArray()
iter.Read(&struct_) // reflection-api
```

Also by using type or field callback, we can switch from reflection-api back to pull-api. The seamless mix of both styles
enabled a unique new way to parse our data.

My advice is always use the reflection-api first. Unless you find pull-api can do a better job in certain area.

# Why not json iterator?

jsoniter does not plan to support `map[string]interface{}`, period.
