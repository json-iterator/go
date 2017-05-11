package jsoniter

import (
	"fmt"
	"testing"
	"github.com/json-iterator/go/require"
	"unsafe"
	"strconv"
)

func Test_bind_api_demo(t *testing.T) {
	should := require.New(t)
	val := []int{}
	err := UnmarshalFromString(`[0,1,2,3]  `, &val)
	should.Nil(err)
	should.Equal([]int{0, 1, 2, 3}, val)
}

func Test_iterator_api_demo(t *testing.T) {
	iter := ParseString(`[0,1,2,3]`)
	total := 0
	for iter.ReadArray() {
		total += iter.ReadInt()
	}
	fmt.Println(total)
}

type DocumentMatch struct {
	Index string   `json:"index,omitempty"`
	ID    string   `json:"id"`
	Score float64  `json:"score"`
	Sort  []string `json:"sort,omitempty"`
}

type DocumentMatchCollection []*DocumentMatch

type SearchResult struct {
	Hits DocumentMatchCollection `json:"hits"`
}

func Test2(t *testing.T) {
	RegisterTypeEncoder("float64", func(ptr unsafe.Pointer, stream *Stream) {
		t := *((*float64)(ptr))
		stream.WriteRaw(strconv.FormatFloat(t, 'E', -1, 64))
	})
	hits := []byte(`{"hits":[{"index":"geo","id":"firehouse_grill_brewery","score":3.584608106366055e-07,
							"sort":[" \u0001@\t\u0007\u0013;a\u001b}W"]},
							{"index":"geo","id":"jack_s_brewing","score":2.3332790568885077e-07,
							"sort":[" \u0001@\u0013{w?.\"0\u0010"]},
							{"index":"geo","id":"brewpub_on_the_green","score":2.3332790568885077e-07,
							"sort":[" \u0001@\u0014\u0017+\u00137QZG"]}]}`)
	var h SearchResult
	err := Unmarshal(hits, &h)
	fmt.Printf("SR %+v \n", h.Hits[0])
	b, err := Marshal(h.Hits[0])
	if err != nil {
		fmt.Printf("error marshalling search res: %v", err)
		//return
	}
	fmt.Printf("SR %s \n", string(b))
}
