package main

import (
	"math"

	"github.com/hajimehoshi/ebiten"
)

type Orb struct {
	*Collide
	blowImages       []*ebiten.Image
	direction        float64
	active           bool
	floating         bool
	timer            int
	blownFrames      int
	trappedEnemyType RobotType
}

func NewOrb(level *Level) *Orb {
	return &Orb{
		Collide: NewCollide(level, NewSprite(XCentre, YBottom)),
		blowImages: []*ebiten.Image{images["orb0"], images["orb1"], images["orb2"],
			images["orb3"], images["orb4"], images["orb5"], images["orb6"]},
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
	o.trappedEnemyType = RobotNone
	o.direction = direction
	o.blownFrames = 6
	o.MoveTo(x, y)
	o.SetSequenceFunc(imageSequence).Animate(o.blowImages, nil, 3, true)
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
		ok := o.CollideMove(o.direction, 0, 4)
		if !ok {
			// can't go further (because of a wall)
			o.floating = true
		}
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

func imageSequence(timer int) int {
	if timer < 9 {
		return timer / 3
	}
	return 3 + int(math.Mod(math.Floor((float64(timer)-9)/8), 4))
}
