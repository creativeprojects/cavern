package main

import (
	"math"

	"github.com/creativeprojects/cavern/lib"
)

type Gravity struct {
	*Collide
	speedY float64
	landed bool
}

func NewGravity(level *Level, sprite *lib.Sprite) *Gravity {
	return &Gravity{
		Collide: NewCollide(level, sprite),
		speedY:  0,
		landed:  false,
	}
}

// UpdateFall updates the coordinates when falling down.
// it returns true when the sprite just landed on a block
func (g *Gravity) UpdateFall() bool {
	// Apply gravity, without going over the maximum fall speed
	g.speedY = math.Min(g.speedY+1, MaxFallSpeed)
	dy := 1.0
	if g.speedY <= 0 {
		dy = -1.0
	}
	moved := g.CollideMove(0, dy, math.Abs(g.speedY))
	if moved {
		g.landed = false
	} else {
		g.speedY = 0
		if !g.landed {
			// just landed on a block
			g.landed = true
			return true
		}
		return false
	}
	if g.Y(lib.YTop) >= WindowHeight {
		// fallen off the bottom, reappear at the top
		g.MoveTo(g.RawX(), g.RawY()-WindowHeight)
	}
	return false
}

// UpdateFreeFall updates the coordinates when falling down without stopping on any block
func (g *Gravity) UpdateFreeFall() {
	// Apply gravity, without going over the maximum fall speed
	g.speedY = math.Min(g.speedY+1, MaxFallSpeed)
	// collision detection disabled
	g.Move(0, g.speedY)
}
