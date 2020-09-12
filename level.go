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
	op               *ebiten.DrawImageOptions
	id               int
	colour           int
	grid             []string
}

// NewLevel creates an empty level. Please call Next() to load the first level
func NewLevel() *Level {
	return &Level{
		id:     -1,
		colour: -1,
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
		op: &ebiten.DrawImageOptions{},
	}
}

// Next changes the color and loads the grid for the next level
func (l *Level) Next() {
	l.id++
	l.colour = int(math.Mod(float64(l.colour+1), 4))
	gridID := int(math.Mod(float64(l.id), float64(len(LevelsDefinition))))
	l.grid = append(LevelsDefinition[gridID], LevelsDefinition[gridID][0])
}

// Draw level
func (l *Level) Draw(screen *ebiten.Image) {
	screen.DrawImage(l.backgroundImages[l.colour], nil)

	for y, line := range l.grid {
		x := LeftGridOffset
		for _, char := range line {
			if char != ' ' {
				l.op.GeoM.Reset()
				l.op.GeoM.Translate(x, float64(y)*GridBlockSize)
				screen.DrawImage(l.blockImages[l.colour], l.op)
			}
			x += GridBlockSize
		}
	}
}

// Block returns true if there's a grid block at these coordinates
func (l *Level) Block(x, y float64) bool {
	gridX := math.Floor((x - LeftGridOffset) / GridBlockSize)
	gridY := math.Floor(y / GridBlockSize)
	if gridY > 0 && gridY < NumRows {
		row := l.grid[int(gridY)]
		return gridX >= 0 && gridX < NumColumns && len(row) > 0 && row[int(gridX)] != ' '
	}
	return false
}
