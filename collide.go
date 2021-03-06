package main

import (
	"math"

	"github.com/creativeprojects/cavern/lib"
)

// Collide manages collision between a sprite and a grid
type Collide struct {
	*lib.Sprite
	level *Level
}

// NewCollide creates a new collision detection between a sprite and a grid
func NewCollide(level *Level, sprite *lib.Sprite) *Collide {
	return &Collide{
		Sprite: sprite,
		level:  level,
	}
}

// CollideMove sets the new position, and returns true if the move was successful (returns false if there was a wall)
func (c *Collide) CollideMove(dx, dy, speed float64) bool {
	newX, newY := c.X(lib.XCentre), c.Y(lib.YBottom)

	// movement is done 1 pixel at a time
	for i := 0; i < int(speed); i++ {
		newX += dx
		newY += dy
		if newX < 70 || newX > 730 {
			// collided with edges of the grid
			return false
		}
		// we check for a block: in the direction we're moving to, and only when we're at the edge of one
		// we don't check for collision when the item is moving up
		if (dy > 0 && math.Mod(float64(newY), GridBlockSize) == 0 ||
			dx > 0 && math.Mod(float64(newX), GridBlockSize) == 0 ||
			dx < 0 && math.Mod(float64(newX), GridBlockSize) == GridBlockSize-1) && c.level.Block(int(newX), int(newY)) {
			return false
		}
		// register the move
		c.MoveToType(newX, newY, lib.XCentre, lib.YBottom)
	}
	return true
}
