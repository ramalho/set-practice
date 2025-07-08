package intset

import (
	"fmt"
	"slices"
	"testing"
)

var (
	fibonacci = []int{0, 1, 2, 3, 5, 8}    // Fibonacci numbers < 10
	primes    = []int{2, 3, 5, 7}          // prime numbers < 10
	union     = []int{0, 1, 2, 3, 5, 7, 8} // fibonnacci ∪ prime
	empty     = []int{}
	single    = []int{100} // set with element > 63
)

func TestAdd(t *testing.T) {
	s := IntSet{}
	s.Add(3)
	for _, n := range fibonacci {
		s.Add(n)
	}
	want := "{0 1 2 3 5 8}"
	got := s.String()
	if want != got {
		t.Errorf("%s != %s", got, want)
	}
	wantLen := 6
	gotLen := s.Len()
	if wantLen != gotLen {
		t.Errorf("%d != %d", gotLen, wantLen)
	}
}

func TestFromSlice(t *testing.T) {
	s := FromSlice(primes)
	want := "{2 3 5 7}"
	got := s.String()
	if want != got {
		t.Errorf("%s != %s", got, want)
	}
}

func TestLen(t *testing.T) {
	s := FromSlice(fibonacci)
	want := 6
	got := s.Len()
	if want != got {
		t.Errorf("%d != %d", got, want)
	}
}

func TestLen_Union(t *testing.T) {
	want := FromSlice(union).Len()
	set := FromSlice(primes)
	set.UnionWith(FromSlice(fibonacci))
	got := set.Len()
	if want != got {
		t.Errorf("%d != %d", got, want)
	}

}

func TestLen_UnionSecondWord(t *testing.T) {
	// UnionWith set that uses 2 words: {100}
	set := FromSlice(fibonacci)
	set.UnionWith(FromSlice(single))
	want := 7
	got := set.Len()
	if want != got {
		t.Errorf("%d != %d", got, want)
	}

}

func TestHas(t *testing.T) {
	var testCases = []struct {
		n    int
		set  *IntSet
		want bool
	}{
		{3, FromSlice(primes), true},
		{4, FromSlice(primes), false},
		{1, FromSlice(empty), false},
	}
	for _, tc := range testCases {
		t.Run(fmt.Sprintf("%d in %v", tc.n, tc.set), func(t *testing.T) {
			got := tc.set.Has(tc.n)
			if got != tc.want {
				t.Errorf("tc.set.Has(%d) = %v; want %v", tc.n, got, tc.want)
			}
		})
	}
}

func TestBitCount(t *testing.T) {
	var testCases = []struct {
		word uint64
		want int
	}{
		{0, 0},
		{1, 1},
		{2, 1},
		{3, 2},
		{10, 2},
		{11, 3},
		{12, 2},
		{13, 3},
		{15, 4},
	}
	for _, tc := range testCases {
		t.Run(fmt.Sprintf("%b has %d", tc.word, tc.want), func(t *testing.T) {
			got := bitCount(tc.word)
			if got != tc.want {
				t.Errorf("bitCount(%b) = %d; want %d", tc.word, got, tc.want)
			}
		})
	}
}

func TestIter(t *testing.T) {
	var testCases = []struct {
		name  string
		input []int
	}{
		{"primes", primes},
		{"empty", empty},
		{"single", single},
		{"fibonacci", fibonacci},
		{"large_set", []int{0, 1, 63, 64, 650, 1270, 1280, 2000}},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			s := FromSlice(tc.input)
			var got []int

			for value := range s.Iter() {
				got = append(got, value)
			}

			if !slices.Equal(got, tc.input) {
				t.Errorf("got %v, input %v", got, tc.input)
			}
		})
	}
}
