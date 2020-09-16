package main

import (
	"math/rand"

	"github.com/hajimehoshi/ebiten"
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
	sprite := NewSprite(XCentre, YBottom)
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

func (f *Fruit) Generate(extra bool) {
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
	f.sprite.
		MoveTo(float64(randomInt(70, 730)), float64(randomInt(75, 400))).
		Animate(f.Animation[f.Type], fruitAnimation, 6, true)
}

// Update returns true when the fruit just expired (and needs an animation)
func (f *Fruit) Update() bool {
	if f.HasExpired() {
		return false
	}
	f.sprite.Update()
	f.TTL--
	if f.TTL == 0 {
		// create pop animation
		return true
	}
	if f.landed {
		return false
	}
	f.UpdateFall()
	return false
}

func (f *Fruit) Draw(screen *ebiten.Image, timer float64) {
	if f.HasExpired() {
		return
	}
	f.sprite.Draw(screen)
}

// HasExpired returns true when TTL is down to zero meaning the fruit is no longer displayed
func (f *Fruit) HasExpired() bool {
	return f.TTL <= 0
}
