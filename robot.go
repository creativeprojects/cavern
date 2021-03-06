package main

import (
	"math"
	"math/rand"

	"github.com/creativeprojects/cavern/lib"
	"github.com/hajimehoshi/ebiten/v2"
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
	imagesLeftFire       [2][]*ebiten.Image
	imagesRightFire      [2][]*ebiten.Image
	trapSounds           [][]byte
	laserSounds          [][]byte
	robotType            RobotType
	alive                bool
	directionX           float64
	speed                float64
	changeDirectionTimer int
	fireTimer            int
}

func NewRobot(level *Level) *Robot {
	sprite := lib.NewSprite(lib.XCentre, lib.YBottom)
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
		imagesLeftFire: [2][]*ebiten.Image{
			{images["robot005"], images["robot006"], images["robot007"]},
			{images["robot105"], images["robot106"], images["robot107"]},
		},
		imagesRightFire: [2][]*ebiten.Image{
			{images["robot015"], images["robot016"], images["robot017"]},
			{images["robot115"], images["robot116"], images["robot117"]},
		},
		trapSounds:  [][]byte{sounds["trap0"], sounds["trap1"], sounds["trap2"], sounds["trap3"]},
		laserSounds: [][]byte{sounds["laser0"], sounds["laser1"], sounds["laser2"], sounds["laser3"]},
	}
}

// Generate a new robot of type robotType
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

func (r *Robot) IsAlive() bool {
	return r.alive
}

func (r *Robot) Update(game *Game) {
	if !r.IsAlive() {
		return
	}
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

	// the more powerful type of robot can deliberately shoot at orbs - turning to face them if necessary
	if r.robotType == RobotAggressive && r.fireTimer >= 24 {
		// go through all the orbs to see if any can be shot
		for _, orb := range game.orbs {
			// the orb must be at our height, and within 200 pixels on the x axis
			if orb.IsActive() &&
				orb.Y(lib.YCentre) >= r.Y(lib.YTop) &&
				orb.Y(lib.YCentre) < r.Y(lib.YBottom) &&
				math.Abs(orb.X(lib.XCentre)-r.X(lib.XCentre)) < 200 {
				if orb.X(lib.XCentre)-r.X(lib.XCentre) < 0 {
					r.directionX = -1
				} else {
					r.directionX = 1
				}
				r.fireTimer = 0
				break
			}
		}
	}
	// check to see if we can fire at player
	if r.fireTimer >= 12 {
		// random chance of firing each frame. Likehood increases 10 times if player is at the same height as us
		probability := game.level.FireProbability()
		if r.Y(lib.YTop) < game.player.sprite.Y(lib.YBottom) && r.Y(lib.YBottom) > game.player.sprite.Y(lib.YTop) {
			probability *= 10
		}
		if rand.Float64() < probability {
			r.fireTimer = 1
			game.RandomSoundEffect(r.laserSounds)
			// change animation
			if r.directionX == -1 {
				r.Sprite.Animate(r.imagesLeftFire[r.robotType-1], nil, 4, false)
			} else {
				r.Sprite.Animate(r.imagesRightFire[r.robotType-1], nil, 4, false)
			}
		}
	} else if r.fireTimer == 8 {
		// once the fire timer has been set to 0, it will count up - frame 8 of the animation is when the actual bolt is fired
		game.Fire(r.directionX, r.X(lib.XCentre)+r.directionX*20, r.Y(lib.YCentre))
	}
	// am I colliding with an Orb? if so, become trapped in it
	for _, orb := range game.orbs {
		if orb.IsActive() && !orb.EnemyTrapped() && r.CollidePoint(orb.X(lib.XCentre), orb.Y(lib.YCentre)) {
			r.alive = false
			orb.TrapEnemy(r.robotType)
			game.RandomSoundEffect(r.trapSounds)
			// no need to go further
			return
		}
	}

	// change animation back to normal after firing
	if r.fireTimer == 12 {
		// put normal animation back
		if r.directionX == -1 {
			r.Sprite.Animate(r.imagesLeft[r.robotType-1], nil, 4, true)
		} else {
			r.Sprite.Animate(r.imagesRight[r.robotType-1], nil, 4, true)
		}
	}

	r.Sprite.Update()
}

func (r *Robot) Draw(screen *ebiten.Image) {
	if !r.IsAlive() {
		return
	}
	r.Sprite.Draw(screen)
}
