// Package venn provides a generic Set type inspired by Python's set,
// designed around the Go 1.23+ iterator conventions.
//
// Methods with the Seq suffix accept any iterable as an [iter.Seq]:
// use [slices.Values] for slices, [maps.Keys] for maps,
// [strings.FieldsSeq] and friends for text, or [Set.All] for another set.
//
// Every Seq method treats its argument as if it were first collected
// into a set (duplicates are irrelevant), consumes it at most once
// (single-use sequences are safe), and requires it to be finite.
//
// This is the "core + new-set ops" candidate API; package vennp is the
// "full parity" candidate, built for side-by-side comparison.
package venn

import (
	"fmt"
	"iter"
	"maps"
	"slices"
	"strings"
)

// Set is a set of elements of some comparable type.
// Set is a map type: create one with Make or Collect before adding
// elements, and range over it directly to iterate.
type Set[E comparable] map[E]struct{}

// Make returns a new Set containing the given elements.
func Make[E comparable](elems ...E) Set[E] {
	s := make(Set[E], len(elems))
	s.AddAll(elems...)
	return s
}

// Collect returns a new Set containing the elements yielded by seq.
func Collect[E comparable](seq iter.Seq[E]) Set[E] {
	s := Set[E]{}
	s.AddSeq(seq)
	return s
}

// Add adds element e to the set.
func (s Set[E]) Add(e E) {
	s[e] = struct{}{}
}

// AddAll adds the given elements to the set.
func (s Set[E]) AddAll(elems ...E) {
	for _, e := range elems {
		s[e] = struct{}{}
	}
}

// Remove removes element e from the set if it is present.
func (s Set[E]) Remove(e E) {
	delete(s, e)
}

// RemoveAll removes the given elements from the set, ignoring
// elements that are not present.
func (s Set[E]) RemoveAll(elems ...E) {
	for _, e := range elems {
		delete(s, e)
	}
}

// Pop removes and returns an arbitrary element of the set.
// The second result is false if the set is empty.
func (s Set[E]) Pop() (E, bool) {
	for e := range s {
		delete(s, e)
		return e, true
	}
	var zero E
	return zero, false
}

// Contains reports whether e is in the set.
func (s Set[E]) Contains(e E) bool {
	_, found := s[e]
	return found
}

// ContainsAll reports whether all the given elements are in the set.
func (s Set[E]) ContainsAll(elems ...E) bool {
	for _, e := range elems {
		if !s.Contains(e) {
			return false
		}
	}
	return true
}

// Len returns the number of elements in the set.
func (s Set[E]) Len() int {
	return len(s)
}

// Clear removes all elements from the set.
func (s Set[E]) Clear() {
	clear(s)
}

// Copy returns a shallow copy of the set.
func (s Set[E]) Copy() Set[E] {
	c := make(Set[E], len(s))
	maps.Copy(c, s)
	return c
}

// ToSlice returns the elements of the set as a new slice,
// in no particular order.
func (s Set[E]) ToSlice() []E {
	return slices.Collect(s.All())
}

// All is an iterator over the elements of the set,
// in no particular order.
func (s Set[E]) All() iter.Seq[E] {
	return maps.Keys(s)
}

// String returns the set formatted as "{a b c}",
// in no particular order.
func (s Set[E]) String() string {
	var b strings.Builder
	b.WriteByte('{')
	first := true
	for e := range s {
		if !first {
			b.WriteByte(' ')
		}
		first = false
		fmt.Fprintf(&b, "%v", e)
	}
	b.WriteByte('}')
	return b.String()
}

// Equal reports whether s and z contain the same elements.
func (s Set[E]) Equal(z Set[E]) bool {
	return maps.Equal(s, z)
}

// SubsetOf reports whether every element of s is in z.
func (s Set[E]) SubsetOf(z Set[E]) bool {
	if len(s) > len(z) {
		return false
	}
	for e := range s {
		if !z.Contains(e) {
			return false
		}
	}
	return true
}

// SupersetOf reports whether every element of z is in s.
func (s Set[E]) SupersetOf(z Set[E]) bool {
	return z.SubsetOf(s)
}

// Union returns a new set with the elements of s and z.
func (s Set[E]) Union(z Set[E]) Set[E] {
	u := make(Set[E], max(len(s), len(z)))
	maps.Copy(u, s)
	maps.Copy(u, z)
	return u
}

// Intersection returns a new set with the elements common to s and z.
func (s Set[E]) Intersection(z Set[E]) Set[E] {
	small, large := s, z
	if len(z) < len(s) {
		small, large = z, s
	}
	i := Set[E]{}
	for e := range small {
		if large.Contains(e) {
			i[e] = struct{}{}
		}
	}
	return i
}

// Difference returns a new set with the elements of s that are not in z.
func (s Set[E]) Difference(z Set[E]) Set[E] {
	d := Set[E]{}
	for e := range s {
		if !z.Contains(e) {
			d[e] = struct{}{}
		}
	}
	return d
}

// SymmetricDifference returns a new set with the elements that are
// in s or in z, but not in both.
func (s Set[E]) SymmetricDifference(z Set[E]) Set[E] {
	d := s.Difference(z)
	for e := range z {
		if !s.Contains(e) {
			d[e] = struct{}{}
		}
	}
	return d
}

// UnionUpdate adds the elements of z to s.
func (s Set[E]) UnionUpdate(z Set[E]) {
	maps.Copy(s, z)
}

// IntersectionUpdate removes from s the elements not in z.
func (s Set[E]) IntersectionUpdate(z Set[E]) {
	for e := range s {
		if !z.Contains(e) {
			delete(s, e)
		}
	}
}

// DifferenceUpdate removes the elements of z from s.
func (s Set[E]) DifferenceUpdate(z Set[E]) {
	for e := range z {
		delete(s, e)
	}
}

// SymmetricDifferenceUpdate updates s to hold the elements that are
// in s or in z, but not in both.
func (s Set[E]) SymmetricDifferenceUpdate(z Set[E]) {
	for e := range z {
		if s.Contains(e) {
			delete(s, e)
		} else {
			s[e] = struct{}{}
		}
	}
}

// AddSeq adds the elements yielded by seq to s.
// It is the iterable counterpart of UnionUpdate (Python: s.update(it)).
func (s Set[E]) AddSeq(seq iter.Seq[E]) {
	for e := range seq {
		s[e] = struct{}{}
	}
}

// RemoveSeq removes the elements yielded by seq from s.
// It is the iterable counterpart of DifferenceUpdate
// (Python: s.difference_update(it)).
func (s Set[E]) RemoveSeq(seq iter.Seq[E]) {
	for e := range seq {
		delete(s, e)
	}
}

// UnionSeq returns a new set with the elements of s and
// the elements yielded by seq.
func (s Set[E]) UnionSeq(seq iter.Seq[E]) Set[E] {
	u := s.Copy()
	u.AddSeq(seq)
	return u
}

// IntersectionSeq returns a new set with the elements of s that are
// yielded by seq.
func (s Set[E]) IntersectionSeq(seq iter.Seq[E]) Set[E] {
	i := Set[E]{}
	for e := range seq {
		if s.Contains(e) {
			i[e] = struct{}{}
		}
	}
	return i
}

// DifferenceSeq returns a new set with the elements of s that are
// not yielded by seq.
func (s Set[E]) DifferenceSeq(seq iter.Seq[E]) Set[E] {
	d := s.Copy()
	d.RemoveSeq(seq)
	return d
}

// SymmetricDifferenceSeq returns a new set with the elements that are
// in s or yielded by seq, but not both.
func (s Set[E]) SymmetricDifferenceSeq(seq iter.Seq[E]) Set[E] {
	return s.SymmetricDifference(Collect(seq))
}
