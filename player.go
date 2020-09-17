package main

import "github.com/hajimehoshi/ebiten"

type Player struct {
	sprite        *Sprite
	gravity       *Gravity
	imageStill    *ebiten.Image
	runLeft       []*ebiten.Image
	runRight      []*ebiten.Image
	jumpLeft      *ebiten.Image
	jumpRight     *ebiten.Image
	landingSounds [][]byte
	lives         int
	health        int
	score         int
	movingX       float64
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
		landingSounds: [][]byte{sounds["land0"], sounds["land1"], sounds["land2"], sounds["land3"]},
	}
}

func (p *Player) Start(level *Level) *Player {
	p.gravity = NewGravity(level, p.sprite)
	p.health = PlayerStartHealth
	p.sprite.MoveTo(WindowWidth/2, 100)
	p.sprite.Animate([]*ebiten.Image{p.imageStill}, nil, 8, true)
	return p
}

func (p *Player) Update(game *Game) {
	if p.health == 0 {
		p.gravity.UpdateFreeFall()
	} else {
		landed := p.gravity.UpdateFall()
		if landed {
			game.RandomSoundEffect(p.landingSounds)
		}
	}
	switch {
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

func (p *Player) Draw(screen *ebiten.Image) {
	p.sprite.Draw(screen)
}

func (p *Player) Move(x, y, speed float64) {
	p.gravity.Move(x, y, speed)
	p.movingX = x
}

func (p *Player) Still() {
	p.movingX = 0
}

func (p *Player) Jump() bool {
	if p.gravity.speedY != 0 || p.gravity.landed == false {
		return false
	}
	p.gravity.speedY = -16
	return true
}
