[![rcard](https://goreportcard.com/badge/github.com/json-iterator/go)](https://goreportcard.com/report/github.com/json-iterator/go)

jsoniter (json-iterator) is fast and flexible JSON parser available in [Java](https://github.com/json-iterator/java) and [Go](https://github.com/json-iterator/go)

# Usage

100% compatibility with standard lib

Replace

```
import "encoding/json"
json.Marshal(&data)
```

with 

```
import "github.com/json-iterator/go"
jsoniter.Marshal(&data)
```

Replace

```
import "encoding/json"
json.Unmarshal(input, &data)
```

with

```
import "github.com/json-iterator/go"
jsoniter.Unmarshal(input, &data)
```

# How to get

```
go get github.com/json-iterator/go
```

# Contribution Welcomed !

Report issue or pull request, or email taowen@gmail.com, or [![Gitter chat](https://badges.gitter.im/gitterHQ/gitter.png)](https://gitter.im/json-iterator/Lobby)
