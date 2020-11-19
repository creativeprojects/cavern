package main

import (
	"fmt"
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
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
	pendingEnemies   []RobotType
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
	l.createPendingEnemies()
}

// ID is the current level number (starting at zero)
func (l *Level) ID() int {
	return l.id
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
	// Level
	DrawTextCentre(screen, []byte(fmt.Sprintf("LEVEL %d", l.id+1)), 451)
}

// Block returns true if there's a grid block at these coordinates
func (l *Level) Block(x, y int) bool {
	gridX := (x - LeftGridOffset) / GridBlockSize
	gridY := y / GridBlockSize
	if gridY > 0 && gridY < NumRows {
		row := l.grid[int(gridY)]
		return gridX >= 0 && gridX < NumColumns && len(row) > 0 && row[int(gridX)] != ' '
	}
	return false
}

// MaxEnemies returns the maximum number of enemies on-screen at once
func (l *Level) MaxEnemies() int {
	return min((l.id+6)/2, 8)
}

// FireProbability returns the likehood per frame of each robot firing a bolt
func (l *Level) FireProbability() float64 {
	return 0.001 + (0.0001 * math.Min(100, float64(l.id)))
}

// NextEnemy returns the type of the next enemy to create, if any
func (l *Level) NextEnemy() RobotType {
	if len(l.pendingEnemies) == 0 {
		return RobotNone
	}
	enemy := l.pendingEnemies[0]
	l.pendingEnemies = l.pendingEnemies[1:]
	return enemy
}

// PendingEnemies returns a count of enemies left to spawn on this level
func (l *Level) PendingEnemies() int {
	return len(l.pendingEnemies)
}

// GetRobotSpawnX return an x coordinate where an enemy can appear from
func (l *Level) GetRobotSpawnX() float64 {
	// Find a spawn location for a robot, by checking the top row of the grid for empty spots
	// Start by choosing a random grid column
	r := rand.Intn(NumColumns)

	for i := 0; i < NumColumns; i++ {
		// Keep looking at successive columns (wrapping round if we go off the right-hand side) until
		// we find one where the top grid column is unoccupied
		gridX := math.Mod(float64(r+i), NumColumns)
		if l.grid[0][int(gridX)] == ' ' {
			return GridBlockSize*gridX + LeftGridOffset + 12
		}
	}
	// surely we should have found a hole...
	return WindowWidth / 2
}

func (l *Level) createPendingEnemies() {
	// At the start of each level we create a list of pending enemies - enemies to be created as the level plays out.
	// When this list is empty, we have no more enemies left to create, and the level will end once we have destroyed
	// all enemies currently on-screen. Each element of the list will be either 0 or 1, where 0 corresponds to
	// a standard enemy, and 1 is a more powerful enemy.
	// First we work out how many total enemies and how many of each type to create
	numEnemies := 10 + l.id
	numStrongEnemies := 1 + int(float64(l.id)/1.5)
	numWeakEnemies := numEnemies - numStrongEnemies
	l.pendingEnemies = make([]RobotType, numEnemies)
	for i := 0; i < numStrongEnemies; i++ {
		l.pendingEnemies[i] = RobotAggressive
	}
	for i := 0; i < numWeakEnemies; i++ {
		l.pendingEnemies[i+numStrongEnemies] = RobotNormal
	}
	// randomize the list
	rand.Shuffle(len(l.pendingEnemies), func(i, j int) {
		l.pendingEnemies[i], l.pendingEnemies[j] = l.pendingEnemies[j], l.pendingEnemies[i]
	})
}
