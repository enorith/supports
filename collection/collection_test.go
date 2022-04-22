package collection_test

import (
	"fmt"
	"testing"

	"github.com/enorith/supports/collection"
)

type StructFoo struct {
	Name string
	Age  int
}

var itemsFoo []StructFoo

func TestReduce(t *testing.T) {
	res := collection.Reduce(itemsFoo, func(r int, t StructFoo) int {
		return r + t.Age
	}, 0)

	t.Log(res)
}

func TestIndexOf(t *testing.T) {
	res := collection.IndexOf(itemsFoo, StructFoo{Name: "foo", Age: 4})

	t.Log(res)
}

func TestSort(t *testing.T) {
	data := collection.SortBy(itemsFoo, func(a, b StructFoo) bool {
		return a.Age < b.Age
	})

	t.Log(data, itemsFoo)
}

func TestGroupBy(t *testing.T) {
	data := collection.GroupBy(itemsFoo, func(t StructFoo) int {
		return t.Age % 4
	})

	for k, v := range data {
		fmt.Println(k, v)
	}
}

func init() {
	itemsFoo = []StructFoo{
		{Name: "foo1", Age: 1},
		{Name: "foo3", Age: 3},
		{Name: "bar2", Age: 2},
		{Name: "bar4", Age: 4},
		{Name: "foo6", Age: 6},
		{Name: "bar5", Age: 5},
		{Name: "foo7", Age: 7},
		{Name: "baz8", Age: 8},
		{Name: "baz11", Age: 11},
		{Name: "bar8", Age: 9},
	}
}
