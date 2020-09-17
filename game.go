package main

import (
	"fmt"
	"image/color"
	"math"
	"math/rand"

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
	slow         bool
	debug        bool
	timer        float64
	space        *Sprite
	level        *Level
	player       *Player
	fruits       []*Fruit
	pops         []*Pop
}

// NewGame creates a new game instance and prepares a demo AI game
func NewGame(audioContext *audio.Context) (*Game, error) {

	m, err := NewAudioPlayer(audioContext)
	if err != nil {
		return nil, err
	}
	level := NewLevel()
	g := &Game{
		audioContext: audioContext,
		musicPlayer:  m,
		state:        StateMenu,
		slow:         false,
		space: NewSprite(XCentre, YCentre).MoveTo(400, 280+45).Animate([]*ebiten.Image{
			images["space0"], images["space1"], images["space2"], images["space3"], images["space4"],
			images["space5"], images["space6"], images["space7"], images["space8"], images["space9"],
		}, nil, 4, true).SetSequenceFunc(func(counter int) int {
			// Draw "Press SPACE" animation, which has 10 frames numbered 0 to 9
			// The first part gives us a number between 0 and 159, based on the game timer
			// Dividing by 4 means we go to a new animation frame every 4 frames
			// We enclose this calculation in the min function, with the other argument being 9, which results in the
			// animation staying on frame 9 for three quarters of the time. Adding 40 to the game timer is done to alter
			// which stage the animation is at when the game first starts
			return int(math.Min((math.Floor(math.Mod(float64(counter)+40, 160)) / 4), 9))
		}),
		level:  level,
		player: NewPlayer(level),
		fruits: make([]*Fruit, 0, 10),
		pops:   make([]*Pop, 0, 10),
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

	g.player.Start(g.level)

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
		g.space.Update()
		return nil
	}
	if g.state == StatePlaying {
		// skip to next level
		if inpututil.IsKeyJustPressed(ebiten.KeyN) {
			g.NextLevel()
		}

		// toggle between slow and normal speed mode
		if inpututil.IsKeyJustPressed(ebiten.KeyS) {
			g.slow = !g.slow
			if g.slow {
				ebiten.SetMaxTPS(GameSlowSpeed)
			} else {
				ebiten.SetMaxTPS(GameNormalSpeed)
			}
		}

		if math.Mod(g.timer, NewFruitEvery) == 0 {
			g.GenerateFruit(true)
		}

		for _, pop := range g.pops {
			pop.Update()
		}

		for _, fruit := range g.fruits {
			fruit.Update(g)
		}

		dx := 0.0
		// player actions
		if ebiten.IsKeyPressed(ebiten.KeyLeft) {
			dx = -1
		}
		if ebiten.IsKeyPressed(ebiten.KeyRight) {
			dx = 1
		}
		if dx != 0 {
			g.player.Move(dx, 0, PlayerDefaultSpeed)
		} else {
			g.player.Still()
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyUp) {
			if g.player.Jump() {
				g.SoundEffect(sounds[soundJump])
			}
		}
		g.player.Update(g)

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
		g.space.Draw(screen)

	case StatePlaying:
		g.level.Draw(screen)

		for _, pop := range g.pops {
			pop.Draw(screen)
		}

		for _, fruit := range g.fruits {
			fruit.Draw(screen, g.timer)
		}

		g.player.Draw(screen)

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

// RandomSoundEffect plays a random sound effect from a list
func (g *Game) RandomSoundEffect(sounds [][]byte) {
	if sounds == nil || len(sounds) == 0 {
		return
	}
	soundID := rand.Intn(len(sounds))
	PlaySE(g.audioContext, sounds[soundID])
}

func (g *Game) GenerateFruit(extra bool) {
	// find a free fruit
	for _, fruit := range g.fruits {
		if fruit.TTL == 0 {
			fruit.Generate(true)
			return
		}
	}
	g.fruits = append(g.fruits, NewFruit(g.level, true))
}

func (g *Game) StartPop(popType PopType, x, y float64) {
	// find a free pop
	for _, pop := range g.pops {
		if pop.HasExpired() {
			pop.Start(popType, x, y)
			return
		}
	}
	// we need a new one
	pop := NewPop()
	pop.Start(popType, x, y)
	g.pops = append(g.pops, pop)
}

func (g *Game) displayDebug(screen *ebiten.Image) {
	template := " TPS: %0.2f \n Level %d - Colour %d \n Fruits %d - Pops %d \n"
	msg := fmt.Sprintf(template,
		ebiten.CurrentTPS(),
		g.level.id,
		g.level.colour,
		len(g.fruits),
		len(g.pops),
	)
	ebitenutil.DebugPrint(screen, msg)

	ebitenutil.DrawLine(screen, 70, 75, 70, 400, color.White)
	ebitenutil.DrawLine(screen, 730, 75, 730, 400, color.White)
	ebitenutil.DrawLine(screen, 70, 75, 730, 75, color.White)
	ebitenutil.DrawLine(screen, 70, 400, 730, 400, color.White)
}
