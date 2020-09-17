package main

import "github.com/hajimehoshi/ebiten"

type Orb struct {
	*Collide
	blowImages  []*ebiten.Image
	direction   float64
	active      bool
	floating    bool
	timer       int
	blownFrames int
}

func NewOrb(level *Level) *Orb {
	return &Orb{
		Collide:    NewCollide(level, NewSprite(XCentre, YBottom)), // TODO YCentre???
		blowImages: []*ebiten.Image{images["orb0"], images["orb1"], images["orb2"]},
	}
}

func (o *Orb) Reset() *Orb {
	o.active = true
	o.timer = 0
	return o
}

func (o *Orb) Start(x, y, direction float64) *Orb {
	o.active = true
	o.floating = false
	o.blownFrames = 6
	o.MoveTo(x, y)
	o.Animate(o.blowImages, nil, 3, false)
	return o
}

func (o *Orb) Blow() {
	o.blownFrames += 4
}

func (o *Orb) IsActive() bool {
	return o.active
}

func (o *Orb) Update(game *Game) {
	if !o.active {
		return
	}
	o.timer++
	if o.floating {
		o.CollideMove(0, -1, float64(randomInt(1, 2)))
	} else {
		o.CollideMove(1, 0, 4)
	}
	if o.timer == o.blownFrames {
		o.floating = true
	}
	if o.timer > OrbMaxTimer {
		o.active = false
		game.StartPop(PopOrb, o.X(XCentre), o.Y(YBottom))
		return
	}
	o.Sprite.Update()
}

func (o *Orb) Draw(screen *ebiten.Image) {
	o.Sprite.Draw(screen)
}
