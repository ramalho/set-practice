// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 165.

// Package intset provides a set of integers based on a bit vector.
package intset

import (
	"bytes"
	"fmt"
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
func (s *IntSet) Add(x int) {
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
		if tword > 0 {
			if i < len(s.words) {
				before := bitCount(s.words[i])
				s.words[i] |= tword
				s.len += bitCount(s.words[i]) - before
			} else {
				s.words = append(s.words, tword)
				s.len += bitCount(tword)
			}
		}
	}
}

//!-intset

//!+string

// String returns the set as a string of the form "{1 2 3}".
func (s *IntSet) String() string {
	var buf bytes.Buffer
	buf.WriteByte('{')
	for i, word := range s.words {
		if word == 0 {
			continue
		}
		for j := 0; j < 64; j++ {
			if word&(1<<uint(j)) != 0 {
				if buf.Len() > len("{") {
					buf.WriteByte(' ')
				}
				fmt.Fprintf(&buf, "%d", 64*i+j)
			}
		}
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

func bitCount(word uint64) int {
	count := 0
	for bit := range uint(64) {
		count += int(word>>bit) & 1
	}
	return count
}

// Iter returns an iterator function that can be used with range-over-function in Go 1.24+
// The iterator yields each integer in the set in ascending order.
func (s *IntSet) Iter() func(func(int) bool) {
	return func(yield func(int) bool) {
		for i, word := range s.words {
			if word == 0 {
				continue
			}
			for j := 0; j < 64; j++ {
				if word&(1<<uint(j)) != 0 {
					value := 64*i + j
					if !yield(value) {
						return // Stop iteration if yield returns false
					}
				}
			}
		}
	}
}
