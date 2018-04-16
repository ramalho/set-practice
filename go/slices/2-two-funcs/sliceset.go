package sliceset

func Contains(slice []string, needle string) bool {
	for _, elem := range slice {
		if needle == elem {
			return true
		}
	}
	return false
}

func ContainsAll(slice, subslice []string) bool {
	for _, needle := range subslice {
		if !Contains(slice, needle) {
			return false
		}
	}
	return true
}
