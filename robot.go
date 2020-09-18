package main

import (
	"math/rand"

	"github.com/hajimehoshi/ebiten"
)

type RobotType int

const (
	RobotNone RobotType = iota
	RobotNormal
	RobotAggressive
)

type Robot struct {
	*Gravity
	imagesLeft           [2][]*ebiten.Image
	imagesRight          [2][]*ebiten.Image
	trapSounds           [][]byte
	robotType            RobotType
	alive                bool
	directionX           float64
	speed                float64
	changeDirectionTimer int
	fireTimer            int
}

func NewRobot(level *Level) *Robot {
	sprite := NewSprite(XCentre, YBottom)
	return &Robot{
		Gravity: NewGravity(level, sprite),
		imagesLeft: [2][]*ebiten.Image{
			{images["robot000"], images["robot001"], images["robot002"], images["robot003"], images["robot004"]},
			{images["robot100"], images["robot101"], images["robot102"], images["robot103"], images["robot104"]},
		},
		imagesRight: [2][]*ebiten.Image{
			{images["robot010"], images["robot011"], images["robot012"], images["robot013"], images["robot014"]},
			{images["robot110"], images["robot111"], images["robot112"], images["robot113"], images["robot114"]},
		},
		trapSounds: [][]byte{sounds["trap0"], sounds["trap1"], sounds["trap2"], sounds["trap3"]},
	}
}

func (r *Robot) Generate(robotType RobotType) *Robot {
	r.alive = true
	r.robotType = robotType
	r.speed = float64(randomInt(1, 4))
	r.directionX = 1
	r.changeDirectionTimer = 0
	r.fireTimer = 100
	x := r.level.GetRobotSpawnX()
	y := -30.0
	r.Sprite.MoveTo(x, y)
	return r
}

func (r *Robot) Update(game *Game) {
	r.changeDirectionTimer--
	r.fireTimer++
	// move in current direction, change direction if we hit a wall
	if !r.CollideMove(r.directionX, 0, r.speed) {
		r.changeDirectionTimer = 0
	}
	if r.changeDirectionTimer <= 0 {
		// randomly choose a direction to move in
		// if there's a player, there's two thirds chance that we'll move towards them
		directions := []float64{-1, 1}
		r.directionX = directions[rand.Intn(len(directions))]
		r.changeDirectionTimer = randomInt(100, 251)
		if r.directionX == -1 {
			r.Sprite.Animate(r.imagesLeft[r.robotType-1], nil, 4, true)
		} else {
			r.Sprite.Animate(r.imagesRight[r.robotType-1], nil, 4, true)
		}
	}
	r.Gravity.UpdateFall()

	// am I colliding with an Orb? if so, become trapped in it
	for _, orb := range game.orbs {
		if orb.trappedEnemyType == RobotNone && r.CollidePoint(orb.X(XCentre), orb.Y(YCentre)) {
			r.alive = false
			orb.floating = true
			orb.trappedEnemyType = r.robotType
			game.RandomSoundEffect(r.trapSounds)
			// no need to go further
			return
		}
	}
	r.Sprite.Update()
}

func (r *Robot) Draw(screen *ebiten.Image) {
	r.Sprite.Draw(screen)
}
