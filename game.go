package main

import (
	"fmt"
	"math"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/audio"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"github.com/hajimehoshi/ebiten/inpututil"
)

// Game contains the current game state
type Game struct {
	audioContext *audio.Context
	musicPlayer  *AudioPlayer
	state        GameState
	debug        bool
	timer        float64
	level        *Level
	fruits       []*Fruit
}

// NewGame creates a new game instance and prepares a demo AI game
func NewGame(audioContext *audio.Context) (*Game, error) {

	m, err := NewAudioPlayer(audioContext)
	if err != nil {
		return nil, err
	}
	g := &Game{
		audioContext: audioContext,
		musicPlayer:  m,
		state:        StateMenu,
		level:        NewLevel(),
		fruits:       make([]*Fruit, 0, 10),
	}

	return g, nil
}

// Layout defines the size of the game in pixels
func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return WindowWidth, WindowHeight
}

// Start initializes a new game
func (g *Game) Start() *Game {
	g.timer = -1
	g.level = NewLevel()
	g.level.Next()

	g.state = StatePlaying
	return g
}

// NextLevel loads the next level
func (g *Game) NextLevel() {
	g.SoundEffect(sounds[soundLevel])
	g.level.Next()
}

// Update game events
func (g *Game) Update(screen *ebiten.Image) error {
	g.timer++

	// Debug
	if inpututil.IsKeyJustPressed(ebiten.KeyD) {
		g.debug = !g.debug
	}

	if g.state == StateMenu {
		if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
			g.Start()
		}
		return nil
	}
	if g.state == StatePlaying {
		// skip to next level
		if inpututil.IsKeyJustPressed(ebiten.KeyN) {
			g.NextLevel()
		}

		if math.Mod(g.timer, 100) == 0 {
			// x 70 to 730, y 75 to 400
			g.fruits = append(g.fruits, NewFruit(g.level, true))
		}
		for _, fruit := range g.fruits {
			fruit.Update()
		}
		return nil
	}

	if g.state == StatePaused {
		// un-pause
		if inpututil.IsKeyJustPressed(ebiten.KeyP) || inpututil.IsKeyJustPressed(ebiten.KeySpace) {
			g.state = StatePlaying
		}
		return nil
	}

	if g.state == StateGameOver {
		// un-pause
		if inpututil.IsKeyJustPressed(ebiten.KeyEnter) || inpututil.IsKeyJustPressed(ebiten.KeySpace) || inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
			g.Reset()
			g.state = StateMenu
		}
		return nil
	}
	return nil
}

// Reset game ready for a new one
func (g *Game) Reset() {
	g.timer = -1
}

// Draw game events
func (g *Game) Draw(screen *ebiten.Image) {

	switch g.state {
	case StateMenu:
		screen.DrawImage(images[imageTitle], nil)
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(130, 280)
		frame := math.Min((math.Floor(math.Mod(g.timer+60, 160)) / 4), 9)
		screen.DrawImage(images[fmt.Sprintf("space%.0f", frame)], op)

		for _, fruit := range g.fruits {
			fruit.Draw(screen, g.timer)
		}

	case StatePlaying:
		g.level.Draw(screen)

	case StateGameOver:
	}

	if g.debug {
		g.displayDebug(screen)
	}
}

// SoundEffect plays a sound in the game
func (g *Game) SoundEffect(se []byte) {
	PlaySE(g.audioContext, se)
}

func (g *Game) displayDebug(screen *ebiten.Image) {
	template := " TPS: %0.2f \n Level %d \n Colour %d \n"
	msg := fmt.Sprintf(template,
		ebiten.CurrentTPS(),
		g.level.id,
		g.level.colour,
	)
	ebitenutil.DebugPrint(screen, msg)
}
