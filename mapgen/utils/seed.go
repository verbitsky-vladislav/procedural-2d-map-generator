package utils

import "math/rand"

func GenerateSeed(size int) []int {

	seed := make([]int, size)

	for i := 0; i < size; i++ {
		seed[i] = rand.Int()
	}

	return seed
}
