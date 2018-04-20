package deckarep

import (
	"fmt"
	"github.com/deckarep/golang-set"
	"sort"
)

func Example_invertedIndex() {
	index := make(map[string]mapset.Set)
	index["CHESS"] = mapset.NewSet('♚', '♛', '♜', '♝', '♞', '♟', '♔', '♕', '♖', '♗', '♘', '♙')
	index["BLACK"] = mapset.NewSet('⚑', '■', '🖤', '★', '☎', '☻', '♚', '♛', '♜', '♝', '♞', '♟', '♠', '♣', '✂', '㉈')
	result := index["CHESS"].Intersect(index["BLACK"])
	// fmt.Println(result) outputs in random order:
	// Set{9818, 9820, 9821, 9823, 9819, 9822}
	// The next five lines make the result testable:
	list := []string{}
	for char := range result.Iter() {
		list = append(list, string(char.(rune)))
	}
	sort.Strings(list)
	fmt.Println(list)
	// Output:
	// [♚ ♛ ♜ ♝ ♞ ♟]
}
