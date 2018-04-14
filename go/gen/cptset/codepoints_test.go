package cptset

import "fmt"

var (
	vowels = []Codepoint{'A', 'E', 'I', 'O', 'U'}
	five   = []Codepoint{'A', 'B', 'C', 'D', 'E'}
)

func ExampleString() {
	vowelSet := NewCodepointSet(vowels...)
	fmt.Println(vowels)
	fmt.Println(vowelSet)
	// Output: 
	// [A E I O U]
	// CodepointSet{A E I O U}
}


func ExampleIntersect() {
	vowelSet := NewCodepointSet(vowels...)
	fiveSet := NewCodepointSet(five...)
	fmt.Println(fiveSet.Intersect(vowelSet))
	// Output: 
	// CodepointSet{A E}
}


func ExampleDifference() {
	vowelSet := NewCodepointSet(vowels...)
	fiveSet := NewCodepointSet(five...)
	fmt.Println(fiveSet.Difference(vowelSet))
	// Output: 
	// CodepointSet{B C D}
}
