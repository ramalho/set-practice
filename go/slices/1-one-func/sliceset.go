package sliceset

func ContainsAll(slice, subslice []string) bool {
	for _, needle := range subslice {
		found := false
		for _, elem := range slice {
			if needle == elem {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}

