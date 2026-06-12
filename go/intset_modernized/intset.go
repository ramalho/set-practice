// Package intset provides a set of integers based on a bit vector.
package intset

import (
	"bytes"
	"fmt"
	"iter"
	"math/bits"
)

//!+intset

// An IntSet is a set of small non-negative integers.
// Its zero value represents the empty set.
type IntSet struct {
	words []uint64
	len   int
}

// Has reports whether the set contains the non-negative value x.
func (s *IntSet) Has(x int) bool {
	word, bit := x/64, uint(x%64)
	return word < len(s.words) && (s.words[word]>>bit)&1 == 1
}

// Add adds the non-negative value x to the set.
// It panics if x is negative.
func (s *IntSet) Add(x int) {
	if x < 0 {
		panic(fmt.Sprintf("intset: negative value %d", x))
	}
	if s.Has(x) {
		return
	}
	word, bit := x/64, uint(x%64)
	for word >= len(s.words) {
		s.words = append(s.words, 0)
	}
	s.words[word] |= 1 << bit
	s.len++
}

func (s *IntSet) Len() int {
	return s.len
}

// UnionWith sets s to the union of s and t.
func (s *IntSet) UnionWith(t *IntSet) {
	for i, tword := range t.words {
		if tword == 0 {
			continue
		}
		for i >= len(s.words) {
			s.words = append(s.words, 0)
		}
		before := bits.OnesCount64(s.words[i])
		s.words[i] |= tword
		s.len += bits.OnesCount64(s.words[i]) - before
	}
}

//!-intset

//!+string

// String returns the set as a string of the form "{1 2 3}".
func (s *IntSet) String() string {
	var buf bytes.Buffer
	buf.WriteByte('{')
	for x := range s.Iter() {
		if buf.Len() > len("{") {
			buf.WriteByte(' ')
		}
		fmt.Fprintf(&buf, "%d", x)
	}
	buf.WriteByte('}')
	return buf.String()
}

//!-string

func FromSlice(slice []int) *IntSet {
	s := IntSet{}
	for _, n := range slice {
		s.Add(n)
	}
	return &s
}

// Iter returns an iterator usable with range-over-function (Go 1.23+).
// It yields each integer in the set in ascending order.
func (s *IntSet) Iter() iter.Seq[int] {
	return func(yield func(int) bool) {
		for i, word := range s.words {
			if word == 0 {
				continue
			}
			for j := range 64 {
				if word&(1<<j) != 0 {
					if !yield(64*i + j) {
						return // Stop iteration if yield returns false
					}
				}
			}
		}
	}
}
