package main

import (
	"fmt"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

type IconType int

const (
	IconLife IconType = iota
	IconPlus
	IconHealth
)

var (
	iconWidths = []float64{
		44,
		40,
		40,
	}
)

type Player struct {
	sprite        *Sprite
	gravity       *Gravity
	imageBlank    *ebiten.Image
	imageStill    *ebiten.Image
	runLeft       []*ebiten.Image
	runRight      []*ebiten.Image
	jumpLeft      *ebiten.Image
	jumpRight     *ebiten.Image
	blowLeft      *ebiten.Image
	blowRight     *ebiten.Image
	recoilLeft    *ebiten.Image
	recoilRight   *ebiten.Image
	imagesFall    [2]*ebiten.Image
	iconImages    [3]*ebiten.Image
	landingSounds [][]byte
	blowSounds    [][]byte
	ouchSounds    [][]byte
	dieSound      []byte
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
		imageBlank:    images[imagePlayerBlank],
		imageStill:    images[imagePlayerStill],
		runLeft:       []*ebiten.Image{images["run00"], images["run01"], images["run02"], images["run03"]},
		runRight:      []*ebiten.Image{images["run10"], images["run11"], images["run12"], images["run13"]},
		jumpLeft:      images[imageJumpLeft],
		jumpRight:     images[imageJumpRight],
		blowLeft:      images[imageBlowLeft],
		blowRight:     images[imageBlowRight],
		recoilLeft:    images[imageRecoilLeft],
		recoilRight:   images[imageRecoilRight],
		imagesFall:    [2]*ebiten.Image{images["fall0"], images["fall1"]},
		iconImages:    [3]*ebiten.Image{images["life"], images["plus"], images["health"]},
		landingSounds: [][]byte{sounds["land0"], sounds["land1"], sounds["land2"], sounds["land3"]},
		blowSounds:    [][]byte{sounds["blow0"] /*sounds["blow1"],*/, sounds["blow2"], sounds["blow3"]},
		ouchSounds:    [][]byte{sounds["ouch0"], sounds["ouch1"], sounds["ouch2"], sounds["ouch3"]},
		dieSound:      sounds["die0"],
		lives:         PlayerStartLives,
	}
}

// String returns a debug string
func (p *Player) String() string {
	return fmt.Sprintf("score %d - health %d - lives %d - blow timer %d - hurt timer %d\nPlayer coordinates %s: %.3f %s: %.3f\n",
		p.score,
		p.health,
		p.lives,
		p.blowTimer,
		p.hurtTimer,
		p.sprite.xType.String(),
		p.sprite.x,
		p.sprite.yType.String(),
		p.sprite.y,
	)
}

func (p *Player) Start(level *Level) *Player {
	p.gravity = NewGravity(level, p.sprite)
	p.Reset()
	p.sprite.Animate([]*ebiten.Image{p.imageStill}, nil, 8, true)
	return p
}

func (p *Player) Reset() {
	p.health = PlayerStartHealth
	p.hurtTimer = PlayerStartInvulnerability
	p.sprite.MoveTo(WindowWidth/2, 100)
}

// Hit tests if the coordinates collide with us and returns yes if it does
func (p *Player) Hit(x, y, directionX float64, game *Game) bool {
	collided := p.sprite.CollidePoint(x, y) && p.hurtTimer < 0
	if collided {
		p.hurtTimer = PlayerStartInvulnerability
		p.health--
		p.gravity.speedY = -12
		p.gravity.landed = false
		p.direction = directionX
		if p.health >= 0 {
			game.RandomSoundEffect(p.ouchSounds)
		} else {
			game.SoundEffect(p.dieSound)
		}
	}
	return collided
}

func (p *Player) Update(game *Game) {
	if p.fireTimer >= 0 {
		p.fireTimer--
	}
	if p.hurtTimer >= 0 {
		p.hurtTimer--
	}
	p.blowTimer++

	if p.gravity.landed {
		// hurt timer starts at 200, but drops to 100 once the player has landed
		p.hurtTimer = min(p.hurtTimer, 100)
	}

	if p.hurtTimer > 100 && p.health > 0 {
		// sideway motion if just being knocked by a bolt
		p.sprite.Move(p.direction*4, 0) // FIXME! this code sends the player inside the walls
	}
	if p.hurtTimer > 100 && p.health <= 0 {
		p.gravity.UpdateFreeFall()
		if p.gravity.Y(YCentre) >= WindowHeight*1.5 {
			p.lives--
			if p.lives >= 0 {
				p.Reset()
			} else {
				game.state = StateGameOver
			}
		}
	} else {
		landed := p.gravity.UpdateFall()
		if landed {
			game.RandomSoundEffect(p.landingSounds)
		}
	}
	switch {
	case p.hurtTimer > 0 && p.hurtTimer%2 == 0:
		p.sprite.Animation([]*ebiten.Image{p.imageBlank}, nil, 8, true)

	case p.hurtTimer > 100 && p.health > 0 && p.direction == -1:
		p.sprite.Animation([]*ebiten.Image{p.recoilLeft}, nil, 8, true)
	case p.hurtTimer > 100 && p.health > 0:
		p.sprite.Animation([]*ebiten.Image{p.recoilRight}, nil, 8, true)

	case p.hurtTimer > 100 && p.health <= 0:
		p.sprite.Animation(p.imagesFall[:], nil, 4, true)

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

	// Draw player score
	scoreBytes := []byte(fmt.Sprintf("%d", p.score))
	DrawText(screen, scoreBytes, float64(WindowWidth-2-(CharWidth(0)*len(scoreBytes))), 451)

	// Draw player score
	p.DrawHealth(screen)
}

func (p *Player) DrawHealth(screen *ebiten.Image) {
	icons := make([]IconType, 0, 6)
	switch {
	case p.lives == 1:
		icons = append(icons, IconLife)
	case p.lives == 2:
		icons = append(icons, IconLife, IconLife)
	case p.lives > 2:
		icons = append(icons, IconLife, IconLife, IconPlus)
	}
	for i := 0; i < p.health; i++ {
		icons = append(icons, IconHealth)
	}
	x := 0.0
	op := &ebiten.DrawImageOptions{}
	for _, icon := range icons {
		op.GeoM.Reset()
		op.GeoM.Translate(x, 450)
		screen.DrawImage(p.iconImages[icon], op)
		x += iconWidths[icon]
	}
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
