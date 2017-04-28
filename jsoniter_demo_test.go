package jsoniter

import (
	"fmt"
	"testing"
)

func Test_bind_api_demo(t *testing.T) {
	iter := ParseString(`[0,1,2,3]`)
	val := []int{}
	iter.ReadVal(&val)
	fmt.Println(val[3])
}

func Test_iterator_api_demo(t *testing.T) {
	iter := ParseString(`[0,1,2,3]`)
	total := 0
	for iter.ReadArray() {
		total += iter.ReadInt()
	}
	fmt.Println(total)
}

type User struct {
	userID int
	name   string
	tags   []string
}

func Test_iterator_and_bind_api(t *testing.T) {
	iter := ParseString(`[123, {"name": "taowen", "tags": ["crazy", "hacker"]}]`)
	user := User{}
	iter.ReadArray()
	user.userID = iter.ReadInt()
	iter.ReadArray()
	iter.ReadVal(&user)
	iter.ReadArray() // array end
	fmt.Println(user)
}

type TaskBidLog struct {
	age int
}

func Test2(t *testing.T) {
	rawString :=`
	{"id":0,"bidId":"bid01492692440885","impId":"imp0","taskId":"1024","bidPrice":80,"winPrice":0,"isWon":0,"createTime":1492692440885,"updateTime":null,"device":"","age":30,"gender":"","location":"[中国, 山西, , ]","conType":"0","os":"iOS","osv":"","brand":"","geo":"","ip":"1.68.4.193","idfa":"","waxUserid":""}`
	var log TaskBidLog
	err := UnmarshalFromString(rawString, &log)
	fmt.Println(err)
	fmt.Println(log.age)
}
