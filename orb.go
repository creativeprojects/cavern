package main

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

type Orb struct {
	*Collide
	blowImages       []*ebiten.Image
	trapImages       [2][]*ebiten.Image
	popSounds        [][]byte
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
		trapImages: [2][]*ebiten.Image{
			{images["trap00"], images["trap01"], images["trap02"], images["trap03"], images["trap04"], images["trap05"], images["trap06"], images["trap07"]},
			{images["trap10"], images["trap11"], images["trap12"], images["trap13"], images["trap14"], images["trap15"], images["trap16"], images["trap17"]},
		},
		popSounds: [][]byte{sounds["pop0"], sounds["pop1"], sounds["pop2"], sounds["pop3"]},
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

func (o *Orb) TrapEnemy(robotType RobotType) {
	o.trappedEnemyType = robotType
	o.floating = true
	o.SetSequenceFunc(nil).Animate(o.trapImages[robotType-1], nil, 4, true)
}

func (o *Orb) EnemyTrapped() bool {
	return o.trappedEnemyType > RobotNone
}

// Hit tests if the coordinates collide with us and returns yes if it does
func (o *Orb) Hit(x, y float64) bool {
	collided := o.CollidePoint(x, y)
	if collided {
		o.timer = OrbMaxTimer - 1
	}
	return collided
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
	if o.timer > OrbMaxTimer || o.Y(YBottom) <= -40 {
		o.active = false
		game.StartPop(PopOrb, o.X(XCentre), o.Y(YBottom))
		// create an extra fruit if an enemy was trapped in it
		if o.trappedEnemyType > RobotNone {
			fruit := game.CreateFruit(true)
			fruit.MoveTo(o.X(XCentre), math.Ceil(o.Y(YBottom)))
			game.RandomSoundEffect(o.popSounds)
		}
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
