package main

// Deduplicate returns a new []string with duplicate strings removed.
func Deduplicate(input []string) (out []string) {
	seen := make(map[string]struct{})
	for _, elem := range input {
		if _, found := seen[elem]; !found {
			out = append(out, elem)
			seen[elem] = struct{}{}
		}
	}
	return out
}
