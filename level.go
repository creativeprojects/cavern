package main

import (
	"math"

	"github.com/hajimehoshi/ebiten"
)

const (
	totalLevels = 4
)

type Level struct {
	backgroundImages [totalLevels]*ebiten.Image
	blockImages      [totalLevels]*ebiten.Image
	id               int
	colour           int
}

func NewLevel() *Level {
	return &Level{
		id:     0,
		colour: 0,
		backgroundImages: [totalLevels]*ebiten.Image{
			images["bg0"],
			images["bg1"],
			images["bg2"],
			images["bg3"],
		},
		blockImages: [totalLevels]*ebiten.Image{
			images["block0"],
			images["block1"],
			images["block2"],
			images["block3"],
		},
	}
}

func (l *Level) Next() {
	l.id++
	l.colour = int(math.Mod(float64(l.id), 4))
}

// Draw level
func (l *Level) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	screen.DrawImage(l.backgroundImages[l.colour], nil)

	for y, line := range append(LevelsDefinition[l.colour], LevelsDefinition[l.colour][0]) {
		x := 50.0
		for _, char := range line {
			if char != ' ' {
				op.GeoM.Reset()
				op.GeoM.Translate(x, float64(y)*GridBlockSize)
				screen.DrawImage(l.blockImages[l.colour], op)
			}
			x += GridBlockSize
		}
	}
}
