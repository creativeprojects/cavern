package main

func min(value1, value2 int) int {
	if value1 < value2 {
		return value1
	}
	return value2
}

func abs(value int) int {
	if value < 0 {
		return -value
	}
	return value
}
