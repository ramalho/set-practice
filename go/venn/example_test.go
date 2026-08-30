package venn_test

import (
	"fmt"
	"maps"
	"slices"
	"strings"

	"venn"
)

// sorted returns the elements of s in order, for deterministic output.
func sorted(s venn.Set[string]) []string {
	return slices.Sorted(s.All())
}

func ExampleCollect() {
	s := venn.Collect(strings.FieldsSeq("the quick brown fox the lazy dog the end"))
	fmt.Println(sorted(s))
	// Output: [brown dog end fox lazy quick the]
}

func ExampleSet_AddSeq() {
	inventory := venn.Make("apple", "banana")
	prices := map[string]int{"banana": 5, "cherry": 7}
	inventory.AddSeq(maps.Keys(prices)) // Python: inventory.update(prices)
	fmt.Println(sorted(inventory))
	// Output: [apple banana cherry]
}

func ExampleSet_RemoveSeq() {
	s := venn.Make("a", "b", "c", "d")
	s.RemoveSeq(slices.Values([]string{"b", "d", "x"}))
	fmt.Println(sorted(s))
	// Output: [a c]
}

func ExampleSet_UnionSeq() {
	s := venn.Make("a", "b")
	z := venn.Make("b", "c")
	u := s.UnionSeq(z.All()) // another set is iterable via All
	fmt.Println(sorted(u))
	// Output: [a b c]
}

func ExampleSet_IntersectionSeq() {
	vowels := venn.Make("a", "e", "i", "o", "u")
	found := vowels.IntersectionSeq(strings.SplitSeq("b-a-n-a-n-a", "-"))
	fmt.Println(sorted(found))
	// Output: [a]
}
