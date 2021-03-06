package main

import (
	"github.com/creativeprojects/cavern/lib"
	"github.com/hajimehoshi/ebiten/v2"
)

type Bolt struct {
	*Collide
	leftImages  []*ebiten.Image
	rightImages []*ebiten.Image
	directionX  float64
	active      bool
}

func NewBolt(level *Level) *Bolt {
	sprite := lib.NewSprite(lib.XCentre, lib.YCentre)
	return &Bolt{
		Collide:     NewCollide(level, sprite),
		leftImages:  []*ebiten.Image{images["bolt00"], images["bolt01"]},
		rightImages: []*ebiten.Image{images["bolt00"], images["bolt01"]},
	}
}

func (b *Bolt) Fire(directionX, x, y float64) *Bolt {
	b.directionX = directionX
	b.active = true
	b.MoveTo(x, y)
	if b.directionX == -1 {
		b.Animate(b.leftImages, nil, 4, true)
	} else if b.directionX == 1 {
		b.Animate(b.rightImages, nil, 4, true)
	}
	return b
}

func (b *Bolt) Update(game *Game) {
	if !b.IsActive() {
		return
	}

	// Move horizontally and check to see if we've collided with a block
	ok := b.CollideMove(b.directionX, 0, BoltSpeed)
	if !ok {
		b.active = false
		return
	}
	// collision with an orb
	for _, orb := range game.orbs {
		if orb.Hit(b.X(lib.XCentre), b.Y(lib.YCentre)) {
			b.active = false
			return
		}
	}
	// collision with a player
	if game.player.Hit(b.X(lib.XCentre), b.Y(lib.YCentre), b.directionX, game) {
		b.active = false
		return
	}

	b.Sprite.Update()
}

func (b *Bolt) Draw(screen *ebiten.Image) {
	if !b.IsActive() {
		return
	}
	b.Sprite.Draw(screen)
}

func (b *Bolt) IsActive() bool {
	return b.active
}
