[![rcard](https://goreportcard.com/badge/github.com/json-iterator/go)](https://goreportcard.com/report/github.com/json-iterator/go)

jsoniter (json-iterator) is fast and flexible JSON parser available in [Java](https://github.com/json-iterator/java) and [Go](https://github.com/json-iterator/go)

# Why jsoniter?

* Jsoniter is the fastest JSON parser. It could be up to 10x faster than normal parser, data binding included. Shameless self [benchmark](http://jsoniter.com/benchmark.html)
* Extremely flexible api. You can mix and match three different styles: bind-api, any-api or iterator-api. Checkout your [api choices](http://jsoniter.com/api.html)
* Unique iterator api can iterate through JSON directly, zero memory allocation! See how [iterator](http://jsoniter.com/api.html#iterator-api) works

# Show off

Here is a quick show off, for more complete report you can checkout the full [benchmark](http://jsoniter.com/benchmark.html) with [in-depth optimization](http://jsoniter.com/benchmark.html#optimization-used) to back the numbers up

![go-medium](http://jsoniter.com/benchmarks/go-medium.png)

# Bind-API is the best

Bind-api should always be the first choice. Given this JSON document `[0,1,2,3]`

Parse with Go bind-api

```go
import "github.com/json-iterator/go"
iter := jsoniter.ParseString(`[0,1,2,3]`)
var := iter.Read()
fmt.Println(val)
```

# Iterator-API for quick extraction

When you do not need to get all the data back, just extract some.

Parse with Go iterator-api

```go
import "github.com/json-iterator/go"
iter := ParseString(`[0, [1, 2], [3, 4], 5]`)
count := 0
for iter.ReadArray() {
    iter.Skip()
    count++
}
fmt.Println(count) // 4
```

# Any-API for maximum flexibility

Parse with Go any-api

```go
import "github.com/json-iterator/go"
iter := jsoniter.ParseString(`[{"field1":"11","field2":"12"},{"field1":"21","field2":"22"}]`)
val := iter.ReadAny()
fmt.Println(val.ToInt(1, "field2")) // 22
```

Notice you can extract from nested data structure, and convert any type to the type to you want.

# How to get

```
go get github.com/json-iterator/go
```

# Contribution Welcomed !

Report issue or pull request, or email taowen@gmail.com, or [![Gitter chat](https://badges.gitter.im/gitterHQ/gitter.png)](https://gitter.im/json-iterator/Lobby)
