jsoniter (json-iterator) is fast and flexible JSON parser available in [Java](https://github.com/json-iterator/java) and [Go](https://github.com/json-iterator/go)

# Why jsoniter?

* Jsoniter is the fastest JSON parser, it could be up to 10x faster than normal parser, data binding included (shameless self [benchmark](/benchmark.html))
* Having a developer friendly api is our #1 prioprity, you can choose from bind-api, any-api or iterator-api or all of them (checkout your [api choices](/api.html))
* Unique iterator api can iterate through JSON directly, zero memory allocation! (see how [iterator](/api.html#iterator-api) works)

# 1 Minute Tutorial

Given this JSON document `[0,1,2,3]`

Parse with Go bind-api

```go
import "github.com/json-iterator/go"
iter := jsoniter.ParseString(`[0,1,2,3]`)
val := []int{}
iter.Read(&val)
fmt.Println(val[3])
```

Parse with Go any-api

```go
import "github.com/json-iterator/go"
iter := jsoniter.ParseString(`[0,1,2,3]`)
val := iter.ReadAny()
fmt.Println(val.Get(3))
```

Parse with Go iterator-api

```go
import "github.com/json-iterator/go"
iter := ParseString(`[0,1,2,3]`)
total := 0
for iter.ReadArray() {
    total += iter.ReadInt()
}
fmt.Println(total)
```

# How to get

```
go get github.com/json-iterator/go
```

# Contribution Welcomed !

Report issue or pull request, or email taowen@gmail.com, or [![Gitter chat](https://badges.gitter.im/gitterHQ/gitter.png)](https://gitter.im/json-iterator/Lobby)