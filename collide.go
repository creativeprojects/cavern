package main

import "math"

type Collide struct {
	level *Level
	x     float64
	y     float64
}

func NewCollide(level *Level) *Collide {
	return &Collide{
		level: level,
	}
}

// Move sets the new position, and returns true if the move was successful (returns false if there was a wall)
func (c *Collide) Move(dx, dy, speed float64) bool {
	newX, newY := c.x, c.y

	// movement is done 1 pixel at a time
	for i := 0; i < int(speed); i++ {
		newX += dx
		newY += dy
		if newX < 70 || newX > 730 {
			// collided with edges of the grid
			return false
		}
		if (dy > 0 && math.Mod(newY, GridBlockSize) == 0 ||
			dx > 0 && math.Mod(newX, GridBlockSize) == 0 ||
			dx < 0 && math.Mod(newX, GridBlockSize) == GridBlockSize-1) && c.level.Block(newX, newY) {
			return false
		}
		// register the move
		c.x, c.y = newX, newY
	}
	return true
}
