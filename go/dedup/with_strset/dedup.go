package main

import "github.com/standupdev/strset"

// Deduplicate returns a new []string with duplicate strings removed.
func Deduplicate(input []string) (out []string) {
	seen := strset.Make()
	for _, elem := range input {
		if !seen.Has(elem) {
			out = append(out, elem)
			seen.Add(elem)
		}
	}
	return out
}
