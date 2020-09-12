package main

import "math"

type Gravity struct {
	*Collide
	speed  float64
	landed bool
}

func NewGravity(level *Level) *Gravity {
	return &Gravity{
		Collide: NewCollide(level),
		speed:   0,
		landed:  false,
	}
}

func (g *Gravity) UpdateFall() {
	// Apply gravity, without going over the maximum fall speed
	g.speed = math.Min(g.speed+1, MaxFallSpeed)
	dy := 1.0
	if g.speed <= 0 {
		dy = -1.0
	}
	if !g.Move(0, dy, math.Abs(g.speed)) {
		// we have landed on a block
		g.speed = 0
		g.landed = true
	}
}

func (g *Gravity) UpdateFreeFall() {
	// Apply gravity, without going over the maximum fall speed
	g.speed = math.Min(g.speed+1, MaxFallSpeed)
	// collision detection disabled
	g.y += g.speed
}
