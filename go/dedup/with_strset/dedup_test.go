package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDeduplicate(t *testing.T) {
	testCases := []struct {
		given []string
		want  []string
	}{
		{[]string{}, []string{}},
		{[]string{"a"}, []string{"a"}},
		{[]string{"a", "a"}, []string{"a"}},
		{[]string{"a", "b"}, []string{"a", "b"}},
		{[]string{"a", "b", "a"}, []string{"a", "b"}},
	}
	for _, tc := range testCases {
		t.Run(fmt.Sprintf("%v deduplicated is %v", tc.given, tc.want), func(t *testing.T) {
			got := Deduplicate(tc.given)
			assert.ElementsMatch(t, tc.want, got)
		})
	}
}
