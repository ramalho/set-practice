package sliceset

func ContainsAll(slice, subslice []string) bool {
	set := make(map[string]struct{})
	for _, elem := range slice {
		set[elem] = struct{}{}
	}
	for _, needle := range subslice {
		if _, found := set[needle]; !found {
			return false
		}
	}
	return true
}
