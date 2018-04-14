package main

import "github.com/standupdev/strset"

// Deduplicate returns a new []string with duplicate strings removed.
func Deduplicate(input []string) (out []string) {
	return strset.Make(input...).ToSlice()
}
