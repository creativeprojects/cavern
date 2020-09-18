package main

import (
	"fmt"
	"math"

	"github.com/hajimehoshi/ebiten"
)

type Player struct {
	sprite        *Sprite
	gravity       *Gravity
	imageStill    *ebiten.Image
	runLeft       []*ebiten.Image
	runRight      []*ebiten.Image
	jumpLeft      *ebiten.Image
	jumpRight     *ebiten.Image
	blowLeft      *ebiten.Image
	blowRight     *ebiten.Image
	landingSounds [][]byte
	blowSounds    [][]byte
	lives         int
	health        int
	score         int
	hurtTimer     int     // how long since we got hurt
	fireTimer     int     // how long since we fired an orb
	blowTimer     int     // how long since we're blowing an Orb
	movingX       float64 // player is moving (-1 for left, 1 for right)
	direction     float64 // direction the player is facing (during and after moving)
	blowingOrb    *Orb    // orb being blown right now / nil if none
}

func NewPlayer(level *Level) *Player {
	sprite := NewSprite(XCentre, YBottom)
	gravity := NewGravity(level, sprite)
	return &Player{
		sprite:        sprite,
		gravity:       gravity,
		imageStill:    images[imagePlayerStill],
		runLeft:       []*ebiten.Image{images["run00"], images["run01"], images["run02"], images["run03"]},
		runRight:      []*ebiten.Image{images["run10"], images["run11"], images["run12"], images["run13"]},
		jumpLeft:      images[imageJumpLeft],
		jumpRight:     images[imageJumpRight],
		blowLeft:      images[imageBlowLeft],
		blowRight:     images[imageBlowRight],
		landingSounds: [][]byte{sounds["land0"], sounds["land1"], sounds["land2"], sounds["land3"]},
		blowSounds:    [][]byte{sounds["blow0"] /*sounds["blow1"],*/, sounds["blow2"], sounds["blow3"]},
		lives:         PlayerStartLives,
	}
}

// Strung returns a debug string
func (p *Player) String() string {
	return fmt.Sprintf("score %d - health %d - lives %d - blow timer %d",
		p.score,
		p.health,
		p.lives,
		p.blowTimer,
	)
}

func (p *Player) Start(level *Level) *Player {
	p.gravity = NewGravity(level, p.sprite)
	p.health = PlayerStartHealth
	p.hurtTimer = PlayerStartInvulnerability
	p.sprite.MoveTo(WindowWidth/2, 100)
	p.sprite.Animate([]*ebiten.Image{p.imageStill}, nil, 8, true)
	return p
}

func (p *Player) Update(game *Game) {
	if p.fireTimer > 0 {
		p.fireTimer--
	}
	p.hurtTimer--
	p.blowTimer++
	if p.health == 0 {
		p.gravity.UpdateFreeFall()
	} else {
		landed := p.gravity.UpdateFall()
		if landed {
			game.RandomSoundEffect(p.landingSounds)
		}
	}
	switch {
	case p.blowingOrb != nil && p.direction == -1:
		p.sprite.Animation([]*ebiten.Image{p.blowLeft}, nil, 8, true)
	case p.blowingOrb != nil:
		p.sprite.Animation([]*ebiten.Image{p.blowRight}, nil, 8, true)
	case !p.gravity.landed && p.movingX < 0:
		p.sprite.Animation([]*ebiten.Image{p.jumpLeft}, nil, 8, true)

	case !p.gravity.landed && p.movingX > 0:
		p.sprite.Animation([]*ebiten.Image{p.jumpRight}, nil, 8, true)

	case p.movingX < 0:
		p.sprite.Animation(p.runLeft, nil, 8, true)

	case p.movingX > 0:
		p.sprite.Animation(p.runRight, nil, 8, true)

	default:
		p.sprite.Animation([]*ebiten.Image{p.imageStill}, nil, 8, true)

	}
	p.sprite.Update()
}

// Draw the player on the screen
func (p *Player) Draw(screen *ebiten.Image) {
	p.sprite.Draw(screen)
}

func (p *Player) Move(x, y, speed float64) {
	p.gravity.CollideMove(x, y, speed)
	p.movingX = x
	p.direction = x
}

// Still tells the player to stop moving
func (p *Player) Still() {
	p.movingX = 0
}

// Jump in the grid
func (p *Player) Jump() bool {
	if p.gravity.speedY != 0 || p.gravity.landed == false {
		return false
	}
	p.gravity.speedY = -16
	return true
}

// Eat a fruit or a bonus
func (p *Player) Eat(fruitType FruitType) {
	switch {
	case fruitType == ExtraHealth:
		p.health = min(PlayerStartHealth, p.health+1)
	case fruitType == ExtraLife:
		p.lives++
	default:
		p.score += (int(fruitType) + 1) * 100
	}
}

func (p *Player) StartBlowing(game *Game) {
	if p.fireTimer > 0 {
		return
	}
	p.blowingOrb = game.NewOrb()
	if p.blowingOrb == nil {
		return
	}
	p.blowTimer = 0
	direction := p.direction
	if direction == 0 {
		direction = 1
	}
	// x position will be 38 pixels in front of the player position, while ensuring it is within the bounds of the level
	x := math.Min(730, math.Max(70, p.sprite.X(XCentre)+direction*38))
	y := p.sprite.Y(YCentre) // -35
	p.blowingOrb.Start(x, y, direction)
	game.RandomSoundEffect(p.blowSounds)
}

// Blowing keeps pushing the orb a bit further
func (p *Player) Blowing(game *Game) {
	if p.blowingOrb == nil {
		return
	}
	p.blowingOrb.Blow()
}

func (p *Player) StopBlowing(game *Game) {
	if p.blowingOrb == nil {
		return
	}
	p.blowingOrb = nil
	// wait a bit until you can blow another one
	p.fireTimer = OrbFireTimer
}
