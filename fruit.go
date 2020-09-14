package main

import (
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten"
)

const (
	totalFruits      = 5
	totalFruitFrames = 3
)

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
	Animation [totalFruits][totalFruitFrames]*ebiten.Image
	op        *ebiten.DrawImageOptions
	TTL       int
}

var (
	fruitAnimation = [...]int{0, 1, 2, 1}
)

// NewFruit creates a new random fruit. If extra is true there's a small chance to also create an extra life and extra health fruit.
func NewFruit(level *Level, extra bool) *Fruit {
	var fruitType FruitType
	if !extra {
		fruitType = FruitType(rand.Intn(2))
	} else {
		// 00 to 09 => apple
		// 10 to 19 => raspberry
		// 20 to 29 => lemon
		// 30 to 38 => extra health
		// 39       => extra life
		pick := rand.Intn(39)
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
	return &Fruit{
		Gravity: NewGravity(level),
		Type:    fruitType,
		Animation: [totalFruits][totalFruitFrames]*ebiten.Image{
			{images["fruit00"], images["fruit01"], images["fruit02"]},
			{images["fruit10"], images["fruit11"], images["fruit12"]},
			{images["fruit20"], images["fruit21"], images["fruit22"]},
			{images["fruit30"], images["fruit31"], images["fruit32"]},
			{images["fruit40"], images["fruit41"], images["fruit42"]},
		},
		op:  &ebiten.DrawImageOptions{},
		TTL: 500,
	}
}

func (f *Fruit) Update() {
	f.TTL--
	if f.TTL == 0 {
		// create pop animation
	}
}

func (f *Fruit) Draw(screen *ebiten.Image, timer float64) {
	frame := int(math.Mod(math.Floor(timer/6), 4))
	f.op.GeoM.Reset()
	f.op.GeoM.Translate(f.x, f.y)
	screen.DrawImage(f.Animation[f.Type][frame], f.op)
}
