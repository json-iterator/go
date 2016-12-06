package jsoniter

import (
	"testing"
	"encoding/json"
)


func Test_skip_number(t *testing.T) {
	iter := ParseString(`[-0.12, "b"]`)
	iter.ReadArray()
	iter.Skip()
	iter.ReadArray()
	if iter.ReadString() != "b" {
		t.FailNow()
	}
}

func Test_skip_null(t *testing.T) {
	iter := ParseString(`[null , "b"]`)
	iter.ReadArray()
	iter.Skip()
	iter.ReadArray()
	if iter.ReadString() != "b" {
		t.FailNow()
	}
}

func Test_skip_true(t *testing.T) {
	iter := ParseString(`[true , "b"]`)
	iter.ReadArray()
	iter.Skip()
	iter.ReadArray()
	if iter.ReadString() != "b" {
		t.FailNow()
	}
}

func Test_skip_false(t *testing.T) {
	iter := ParseString(`[false , "b"]`)
	iter.ReadArray()
	iter.Skip()
	iter.ReadArray()
	if iter.ReadString() != "b" {
		t.FailNow()
	}
}

func Test_skip_array(t *testing.T) {
	iter := ParseString(`[[1, [2, [3], 4]], "b"]`)
	iter.ReadArray()
	iter.Skip()
	iter.ReadArray()
	if iter.ReadString() != "b" {
		t.FailNow()
	}
}

func Test_skip_empty_array(t *testing.T) {
	iter := ParseString(`[ [ ], "b"]`)
	iter.ReadArray()
	iter.Skip()
	iter.ReadArray()
	if iter.ReadString() != "b" {
		t.FailNow()
	}
}

func Test_skip_nested(t *testing.T) {
	iter := ParseString(`[ {"a" : [{"b": "c"}], "d": 102 }, "b"]`)
	iter.ReadArray()
	iter.Skip()
	iter.ReadArray()
	if iter.ReadString() != "b" {
		t.FailNow()
	}
}

type TestResp struct {
	Code uint64
}

func Benchmark_jsoniter_skip(b *testing.B) {
	input := []byte(`
{
    "_shards":{
        "total" : 5,
        "successful" : 5,
        "failed" : 0
    },
    "hits":{
        "total" : 1,
        "hits" : [
            {
                "_index" : "twitter",
                "_type" : "tweet",
                "_id" : "1",
                "_source" : {
                    "user" : "kimchy",
                    "postDate" : "2009-11-15T14:12:12",
                    "message" : "trying out Elasticsearch"
                }
            }
        ]
    },
    "code": 200
}`)
	for n := 0; n < b.N; n++ {
		result := TestResp{}
		iter := ParseBytes(input)
		for field := iter.ReadObject(); field != ""; field = iter.ReadObject() {
			switch field {
			case "code":
				result.Code = iter.ReadUint64()
			default:
				iter.Skip()
			}
		}
	}
}

func Benchmark_json_skip(b *testing.B) {
	input := []byte(`
{
    "_shards":{
        "total" : 5,
        "successful" : 5,
        "failed" : 0
    },
    "hits":{
        "total" : 1,
        "hits" : [
            {
                "_index" : "twitter",
                "_type" : "tweet",
                "_id" : "1",
                "_source" : {
                    "user" : "kimchy",
                    "postDate" : "2009-11-15T14:12:12",
                    "message" : "trying out Elasticsearch"
                }
            }
        ]
    },
    "code": 200
}`)
	for n := 0; n < b.N; n++ {
		result := TestResp{}
		json.Unmarshal(input, &result)
	}
}