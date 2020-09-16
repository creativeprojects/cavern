package main

import "math"

type Gravity struct {
	*Collide
	speed  float64
	landed bool
}

func NewGravity(level *Level, sprite *Sprite) *Gravity {
	return &Gravity{
		Collide: NewCollide(level, sprite),
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
	if g.sprite.Y(YTop) >= WindowHeight {
		// fallen off the bottom, reappear at the top
		g.sprite.MoveTo(g.sprite.x, g.sprite.y-WindowHeight)
	}
}

func (g *Gravity) UpdateFreeFall() {
	// Apply gravity, without going over the maximum fall speed
	g.speed = math.Min(g.speed+1, MaxFallSpeed)
	// collision detection disabled
	g.sprite.Move(0, g.speed)
}
