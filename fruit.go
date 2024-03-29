package main

import (
	"math/rand"

	"github.com/creativeprojects/cavern/lib"
	"github.com/hajimehoshi/ebiten/v2"
)

const (
	totalFruits = 5
)

// FruitType describes the type of sprite (apple, raspberry, lemon, health or life)
type FruitType int

// Fruit type
const (
	Apple FruitType = iota
	Raspberry
	Lemon
	ExtraHealth
	ExtraLife
)

type Fruit struct {
	*Gravity
	Type      FruitType
	Animation [totalFruits][]*ebiten.Image
	op        *ebiten.DrawImageOptions
	TTL       int
}

var (
	fruitAnimation = []int{0, 1, 2, 1}
)

// NewFruit creates a new random fruit. If extra is true there's a small chance to also create an extra life and extra health fruit.
func NewFruit(level *Level, extra bool) *Fruit {
	sprite := lib.NewSprite(lib.XCentre, lib.YBottom)
	f := &Fruit{
		Gravity: NewGravity(level, sprite),
		Animation: [totalFruits][]*ebiten.Image{
			{images["fruit00"], images["fruit01"], images["fruit02"]},
			{images["fruit10"], images["fruit11"], images["fruit12"]},
			{images["fruit20"], images["fruit21"], images["fruit22"]},
			{images["fruit30"], images["fruit31"], images["fruit32"]},
			{images["fruit40"], images["fruit41"], images["fruit42"]},
		},
		op: &ebiten.DrawImageOptions{},
	}
	f.Generate(extra)
	return f
}

// Generate a new fruit. If extra is set to yes, a health or life can be generated.
func (f *Fruit) Generate(extra bool) *Fruit {
	var fruitType FruitType
	if !extra {
		fruitType = FruitType(rand.Intn(2))
	} else {
		// 00 to 09 => apple
		// 10 to 19 => raspberry
		// 20 to 29 => lemon
		// 30 to 38 => extra health
		// 39       => extra life
		pick := rand.Intn(40)
		switch {
		case pick <= 9:
			fruitType = Apple
		case pick <= 19:
			fruitType = Raspberry
		case pick <= 29:
			fruitType = Lemon
		case pick <= 38:
			fruitType = ExtraHealth
		default:
			fruitType = ExtraLife
			break
		}
	}

	f.Type = fruitType
	f.landed = false
	f.TTL = FruitTTL
	f.
		MoveTo(float64(randomInt(70, 730)), float64(randomInt(75, 400))).
		Animate(f.Animation[f.Type], fruitAnimation, 6, true)
	return f
}

// Update fruit gravity, expiration, and collision with player
func (f *Fruit) Update(game *Game) {
	if f.HasExpired() {
		return
	}
	f.Sprite.Update()
	f.TTL--
	if f.TTL == 0 {
		// create pop animation
		game.StartPop(PopFruit, f.X(lib.XCentre), f.Y(lib.YBottom))
		return
	}
	if game.player != nil && game.player.sprite.CollidePoint(f.X(lib.XCentre), f.Y(lib.YCentre)) {
		f.TTL = 0
		switch f.Type {
		case ExtraHealth:
			game.SoundEffect(sounds[soundBonus])
		case ExtraLife:
			game.SoundEffect(sounds[soundLife])
		default:
			game.SoundEffect(sounds[soundScore])
		}
		game.player.Eat(f.Type)
	}
	if f.landed {
		return
	}
	f.UpdateFall()
	return
}

func (f *Fruit) Draw(screen *ebiten.Image) {
	if f.HasExpired() {
		return
	}
	f.Sprite.Draw(screen)
}

// HasExpired returns true when TTL is down to zero meaning the fruit is no longer displayed
func (f *Fruit) HasExpired() bool {
	return f.TTL <= 0
}
