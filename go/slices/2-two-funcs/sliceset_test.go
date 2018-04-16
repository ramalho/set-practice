package sliceset

import (
	"fmt"
	"testing"
)

func TestContainsAll(t *testing.T) {
	var testCases = []struct {
		s    []string
		sub  []string
		want bool
	}{
		{[]string{"a", "b"}, []string{"a"}, true},
		{[]string{"a"}, []string{"a", "b"}, false},
		{[]string{"a", "b", "c"}, []string{"b", "a"}, true},
		{[]string{"a", "b", "c"}, []string{"b", "z"}, false},
		{[]string{"a"}, []string{}, true},
		{[]string{}, []string{}, true},
		{[]string{}, []string{"a"}, false},
	}
	for _, tc := range testCases {
		t.Run(fmt.Sprintf("ContainsAll(%v, %v)", tc.s, tc.sub), func(t *testing.T) {
			got := ContainsAll(tc.s, tc.sub)
			if got != tc.want {
				t.Errorf("ContainsAll(%v, %v) = %v; want %v",
					tc.s, tc.sub, got, tc.want)
			}
		})
	}
}