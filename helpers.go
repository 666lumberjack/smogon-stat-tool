package main

// contains takes a value of any type and a slice containing value of that type and returns
// true if the slice contains that value, or false if not
func contains[T comparable](slice []T, value T) bool {
	for _, item := range slice {
		if item == value {
			return true
		}
	}
	return false
}
