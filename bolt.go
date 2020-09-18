package main

import "github.com/hajimehoshi/ebiten"

type Bolt struct {
	*Collide
}

func NewBolt(level *Level) *Bolt {
	sprite := NewSprite(XCentre, YCentre)
	return &Bolt{
		Collide: NewCollide(level, sprite),
	}
}

func (b *Bolt) Update(game *Game) {
	b.Sprite.Update()
}

func (b *Bolt) Draw(screen *ebiten.Image) {
	b.Sprite.Draw(screen)
}
