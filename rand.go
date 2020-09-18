package main

import "math/rand"

// randomInt picks a number between low and high-1 (low included)
// this could seem like a strange behavior, but it allows for `randomInt(0, len(slice))`
func randomInt(low, high int) int {
	return rand.Intn(high-low) + low
}
