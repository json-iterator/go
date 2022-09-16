package jsoniter

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func ExampleMarshal() {
	type ColorGroup struct {
		ID     int
		Name   string
		Colors []string
	}
	group := ColorGroup{
		ID:     1,
		Name:   "Reds",
		Colors: []string{"Crimson", "Red", "Ruby", "Maroon"},
	}
	b, err := Marshal(group)
	if err != nil {
		fmt.Println("error:", err)
	}
	os.Stdout.Write(b)
	// Output:
	// {"ID":1,"Name":"Reds","Colors":["Crimson","Red","Ruby","Maroon"]}
}

func ExampleUnmarshal() {
	var jsonBlob = []byte(`[
		{"Name": "Platypus", "Order": "Monotremata"},
		{"Name": "Quoll",    "Order": "Dasyuromorphia"}
	]`)
	type Animal struct {
		Name  string
		Order string
	}
	var animals []Animal
	err := Unmarshal(jsonBlob, &animals)
	if err != nil {
		fmt.Println("error:", err)
	}
	fmt.Printf("%+v", animals)
	// Output:
	// [{Name:Platypus Order:Monotremata} {Name:Quoll Order:Dasyuromorphia}]
}

func ExampleConfigFastest_Marshal() {
	type ColorGroup struct {
		ID     int
		Name   string
		Colors []string
	}
	group := ColorGroup{
		ID:     1,
		Name:   "Reds",
		Colors: []string{"Crimson", "Red", "Ruby", "Maroon"},
	}
	stream := ConfigFastest.BorrowStream(nil)
	defer ConfigFastest.ReturnStream(stream)
	stream.WriteVal(group)
	if stream.Error != nil {
		fmt.Println("error:", stream.Error)
	}
	os.Stdout.Write(stream.Buffer())
	// Output:
	// {"ID":1,"Name":"Reds","Colors":["Crimson","Red","Ruby","Maroon"]}
}

func ExampleConfigFastest_Unmarshal() {
	var jsonBlob = []byte(`[
		{"Name": "Platypus", "Order": "Monotremata"},
		{"Name": "Quoll",    "Order": "Dasyuromorphia"}
	]`)
	type Animal struct {
		Name  string
		Order string
	}
	var animals []Animal
	iter := ConfigFastest.BorrowIterator(jsonBlob)
	defer ConfigFastest.ReturnIterator(iter)
	iter.ReadVal(&animals)
	if iter.Error != nil {
		fmt.Println("error:", iter.Error)
	}
	fmt.Printf("%+v", animals)
	// Output:
	// [{Name:Platypus Order:Monotremata} {Name:Quoll Order:Dasyuromorphia}]
}

func ExampleGet() {
	val := []byte(`{"ID":1,"Name":"Reds","Colors":["Crimson","Red","Ruby","Maroon"]}`)
	fmt.Printf(Get(val, "Colors", 0).ToString())
	// Output:
	// Crimson
}

func ExampleMyKey() {
	hello := MyKey("hello")
	output, _ := Marshal(map[*MyKey]string{&hello: "world"})
	fmt.Println(string(output))
	obj := map[*MyKey]string{}
	Unmarshal(output, &obj)
	for k, v := range obj {
		fmt.Println(*k, v)
	}
	// Output:
	// {"Hello":"world"}
	// Hel world
}

type MyKey string

func (m *MyKey) MarshalText() ([]byte, error) {
	return []byte(strings.Replace(string(*m), "h", "H", -1)), nil
}

func (m *MyKey) UnmarshalText(text []byte) error {
	*m = MyKey(text[:3])
	return nil
}

type Target struct {
	FieldA string `json:"fieldA"`
}

func Example_duplicateFieldsCaseSensitive() {
	api := Config{
		CaseSensitive:           true,
		DisallowDuplicateFields: true,
	}.Froze()

	t := &Target{}
	err := api.Unmarshal([]byte(`{"fieldA": "value", "fielda": "val2"}`), t)
	fmt.Printf("Case-sensitiveness means no duplicates: 'fieldA' = %q, err = %v\n", t.FieldA, err)

	t = &Target{}
	err = api.Unmarshal([]byte(`{"fieldA": "value", "fieldA": "val2"}`), t)
	fmt.Printf("Got duplicates in struct field: 'fieldA' = %q, err = %v\n", t.FieldA, err)

	t = &Target{}
	err = api.Unmarshal([]byte(`{"fielda": "value", "fielda": "val2"}`), t)
	fmt.Printf("Got duplicates not in struct field: 'fieldA' = %q, err = %v\n", t.FieldA, err)

	// Output:
	// Case-sensitiveness means no duplicates: 'fieldA' = "value", err = <nil>
	// Got duplicates in struct field: 'fieldA' = "value", err = jsoniter.Target.ReadObject: found duplicate field: fieldA, error found in #10 byte of ...|, "fieldA": "val2"}|..., bigger context ...|{"fieldA": "value", "fieldA": "val2"}|...
	// Got duplicates not in struct field: 'fieldA' = "", err = jsoniter.Target.ReadObject: found duplicate field: fielda, error found in #10 byte of ...|, "fielda": "val2"}|..., bigger context ...|{"fielda": "value", "fielda": "val2"}|...
}

func Example_noDuplicateFieldsCaseSensitive() {
	api := Config{
		CaseSensitive:           true,
		DisallowDuplicateFields: false,
	}.Froze()

	t := &Target{}
	err := api.Unmarshal([]byte(`{"fieldA": "value", "fielda": "val2"}`), t)
	fmt.Printf("Case-sensitiveness means no duplicates: 'fieldA' = %q, err = %v\n", t.FieldA, err)

	t = &Target{}
	err = api.Unmarshal([]byte(`{"fieldA": "value", "fieldA": "val2"}`), t)
	fmt.Printf("Got duplicates in struct field: 'fieldA' = %q, err = %v\n", t.FieldA, err)

	t = &Target{}
	err = api.Unmarshal([]byte(`{"fielda": "value", "fielda": "val2"}`), t)
	fmt.Printf("Got duplicates not in struct field: 'fieldA' = %q, err = %v\n", t.FieldA, err)

	// Output:
	// Case-sensitiveness means no duplicates: 'fieldA' = "value", err = <nil>
	// Got duplicates in struct field: 'fieldA' = "val2", err = <nil>
	// Got duplicates not in struct field: 'fieldA' = "", err = <nil>
}

func Example_duplicateFieldsInCaseSensitive() {
	api := Config{
		CaseSensitive:           false,
		DisallowDuplicateFields: true,
	}.Froze()

	t := &Target{}
	err := api.Unmarshal([]byte(`{"fieldA": "value", "fielda": "val2"}`), t)
	fmt.Printf("In-case-sensitive duplicates: 'fieldA' = %q, err = %v\n", t.FieldA, err)

	t = &Target{}
	err = api.Unmarshal([]byte(`{"fieldA": "value", "fieldA": "val2"}`), t)
	fmt.Printf("Got duplicates in exact struct field match: 'fieldA' = %q, err = %v\n", t.FieldA, err)

	t = &Target{}
	err = api.Unmarshal([]byte(`{"fielda": "value", "fielda": "val2"}`), t)
	fmt.Printf("Got duplicates not in notexact struct field match: 'fieldA' = %q, err = %v\n", t.FieldA, err)

	// Output:
	// In-case-sensitive duplicates: 'fieldA' = "value", err = jsoniter.Target.ReadObject: found duplicate field: fielda, error found in #10 byte of ...|, "fielda": "val2"}|..., bigger context ...|{"fieldA": "value", "fielda": "val2"}|...
	// Got duplicates in exact struct field match: 'fieldA' = "value", err = jsoniter.Target.ReadObject: found duplicate field: fieldA, error found in #10 byte of ...|, "fieldA": "val2"}|..., bigger context ...|{"fieldA": "value", "fieldA": "val2"}|...
	// Got duplicates not in notexact struct field match: 'fieldA' = "value", err = jsoniter.Target.ReadObject: found duplicate field: fielda, error found in #10 byte of ...|, "fielda": "val2"}|..., bigger context ...|{"fielda": "value", "fielda": "val2"}|...
}

func Example_noDuplicateFieldsInCaseSensitive() {
	api := Config{
		CaseSensitive:           false,
		DisallowDuplicateFields: false,
	}.Froze()

	t := &Target{}
	err := api.Unmarshal([]byte(`{"fieldA": "value", "fielda": "val2"}`), t)
	fmt.Printf("Case-sensitiveness means no duplicates: 'fieldA' = %q, err = %v\n", t.FieldA, err)

	t = &Target{}
	err = api.Unmarshal([]byte(`{"fieldA": "value", "fieldA": "val2"}`), t)
	fmt.Printf("Got duplicates in struct field: 'fieldA' = %q, err = %v\n", t.FieldA, err)

	t = &Target{}
	err = api.Unmarshal([]byte(`{"fielda": "value", "fielda": "val2"}`), t)
	fmt.Printf("Got duplicates not in struct field: 'fieldA' = %q, err = %v\n", t.FieldA, err)

	// Output:
	// Case-sensitiveness means no duplicates: 'fieldA' = "val2", err = <nil>
	// Got duplicates in struct field: 'fieldA' = "val2", err = <nil>
	// Got duplicates not in struct field: 'fieldA' = "val2", err = <nil>
}

func TestEncoder(t *testing.T) {
	api := Config{
		CaseSensitive:           true,
		DisallowDuplicateFields: true,
	}.Froze()

	type target2 struct {
		B Target `json:"b"`
		A Target `json:"a"`
	}

	data := `{"a": {"fieldA": "bla"}, "b": {"fieldA": "bar"}}`
	data += data

	d := api.NewDecoder(strings.NewReader(data))
	obj := &target2{}
	assert.Nil(t, d.Decode(obj))
	assert.Equal(t, "bla", obj.A.FieldA)
	assert.Equal(t, "bar", obj.B.FieldA)

	obj = &target2{}
	assert.Nil(t, d.Decode(obj))
	assert.Equal(t, "bla", obj.A.FieldA)
	assert.Equal(t, "bar", obj.B.FieldA)
}
