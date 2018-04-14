// Package cptset implements a set of Codepoint (a.k.a. rune) elements
package cptset

import (
	"bytes"
	"sort"
)

// +gen set
type Codepoint rune

func (c Codepoint) String() string {
	return string(c)
}

// ToStringSlice returns elements as a slice of strings
func (set CodepointSet) ToStringSlice() []string {
	var s []string
	for v := range set {
		s = append(s, string(v))
	}
	return s
}


// String returns a string representation of the set
func (set CodepointSet) String() string {
	elems := set.ToStringSlice()
	sort.Strings(elems)
	var buf bytes.Buffer
	buf.WriteString("CodepointSet{")
	for i, elem := range elems {
		if i > 0 {
			buf.WriteByte(' ')
		}
		buf.WriteString(elem)
	}
	buf.WriteByte('}')
	return buf.String()
}
