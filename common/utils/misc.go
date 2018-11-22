package utils

func IsInSlice(str string, slice []string) bool {
	for _, value := range slice {
		if value == str {
			return true
		}
	}
	return false
}

