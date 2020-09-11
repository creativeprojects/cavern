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
	}

	return g, nil
}

// Layout defines the size of the game in pixels
func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return WindowWidth, WindowHeight
}

// Start initializes a game with a number of players
func (g *Game) Start(players int) *Game {
	if players < 0 || players > 2 {
		players = 0
	}
	g.timer = -1

	g.state = StatePlaying
	return g
}

// Update game events
func (g *Game) Update(screen *ebiten.Image) error {
	g.timer++

	// Debug
	if inpututil.IsKeyJustPressed(ebiten.KeyD) {
		g.debug = !g.debug
	}

	if g.state == StateMenu {

		return nil
	}
	if g.state == StatePlaying {

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

	case StatePlaying:

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
	template := " TPS: %0.2f \n "
	msg := fmt.Sprintf(template,
		ebiten.CurrentTPS(),
	)
	ebitenutil.DebugPrint(screen, msg)
}
