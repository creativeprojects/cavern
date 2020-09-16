package main

import "math/rand"

func randomInt(low, high int) int {
	return rand.Intn(high-low) + low
}
