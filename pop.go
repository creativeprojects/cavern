package main

import (
	"github.com/hajimehoshi/ebiten"
)

type PopType int

const (
	PopFruit PopType = iota
	PopOrb
)

// Pop animation
type Pop struct {
	images [2][]*ebiten.Image
	Type   PopType
	sprite *Sprite
}

// NewPop creates a new blank pop animation.
func NewPop() *Pop {
	i := &Pop{
		images: [2][]*ebiten.Image{
			{images["pop00"], images["pop01"], images["pop02"], images["pop03"], images["pop04"], images["pop05"], images["pop06"]},
			{images["pop10"], images["pop11"], images["pop12"], images["pop13"], images["pop14"], images["pop15"], images["pop16"]},
		},
		sprite: NewSprite(XCentre, YBottom),
	}
	return i
}

// Start (and restart) the pop animation on coordinates from another sprite (X centre & Y bottom)
func (i *Pop) Start(popType PopType, x, y float64) *Pop {
	i.Type = popType
	i.sprite.MoveTo(x, y).Animate(i.images[i.Type], nil, 2, false)
	return i
}

func (i *Pop) Update() {
	if i.HasExpired() {
		return
	}
	i.sprite.Update()
}

func (i *Pop) Draw(screen *ebiten.Image) {
	if i.HasExpired() {
		return
	}
	i.sprite.Draw(screen)
}

// HasExpired returns true when the animation is finished
func (i *Pop) HasExpired() bool {
	return i.sprite.IsFinished()
}
