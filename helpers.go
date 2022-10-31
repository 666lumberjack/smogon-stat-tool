package main

import (
	"math"

	"golang.org/x/exp/constraints"
)

type Number interface {
	constraints.Float | constraints.Integer
}

// contains takes a value of comparable type and a slice containing value of that type and returns
// true if the slice contains that value, or false if not
func contains[T comparable](slice []T, value T) bool {
	for _, item := range slice {
		if item == value {
			return true
		}
	}
	return false
}

// closestIndex takes an number and a slice of numbers and returns the int64 index of the slice value closest to the provided int
func closestIndex[N Number](slice []N, value N) int {
	closestIndex := 0
	for i, v := range slice {
		if math.Abs(float64(v-value)) < math.Abs(float64(slice[closestIndex]-value)) {
			closestIndex = i
		}
	}
	return closestIndex
}
