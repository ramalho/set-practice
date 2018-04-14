package main

// Deduplicate returns a new []string with duplicate strings removed.
func Deduplicate(input []string) []string {
	unique := make(map[string]struct{})
	for _, elem := range input {
		unique[elem] = struct{}{}
	}
	out := make([]string, len(unique))
	i := 0
	for k := range unique {
		out[i] = k
		i++
	}
	return out
}
