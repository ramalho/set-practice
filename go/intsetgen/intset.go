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

// UnsignedInteger is a type constraint for all unsigned integer types
type UnsignedInteger interface {
	~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr
}

// An IntSet is a set of small unsigned integers.
// Its zero value represents the empty set.
type IntSet[T UnsignedInteger] struct {
	words []uint64
	len   int
}

// Has reports whether the set contains the value x.
func (s *IntSet[T]) Has(x T) bool {
	word, bit := uint64(x)/64, uint(uint64(x)%64)
	return word < uint64(len(s.words)) && (s.words[word]>>bit)&1 == 1
}

// Add adds the value x to the set.
func (s *IntSet[T]) Add(x T) {
	if s.Has(x) {
		return
	}
	word, bit := uint64(x)/64, uint(uint64(x)%64)
	for word >= uint64(len(s.words)) {
		s.words = append(s.words, 0)
	}
	s.words[word] |= 1 << bit
	s.len++
}

func (s *IntSet[T]) Len() int {
	return s.len
}

// UnionWith sets s to the union of s and t.
func (s *IntSet[T]) UnionWith(t *IntSet[T]) {
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
func (s *IntSet[T]) String() string {
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
				fmt.Fprintf(&buf, "%d", T(64*i+j))
			}
		}
	}
	buf.WriteByte('}')
	return buf.String()
}

//!-string

func FromSlice[T UnsignedInteger](slice []T) *IntSet[T] {
	s := IntSet[T]{}
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
func (s *IntSet[T]) Iter() func(func(T) bool) {
	return func(yield func(T) bool) {
		for i, word := range s.words {
			if word == 0 {
				continue
			}
			for j := 0; j < 64; j++ {
				if word&(1<<uint(j)) != 0 {
					value := T(64*i + j)
					if !yield(value) {
						return // Stop iteration if yield returns false
					}
				}
			}
		}
	}
}

// Example usage:
// var set1 IntSet[uint32]
// set1.Add(42)
// set1.Add(100)
//
// var set2 IntSet[uint8]
// set2.Add(5)
// set2.Add(255)
//
// set3 := FromSlice([]uint16{1, 2, 3, 4, 5})
