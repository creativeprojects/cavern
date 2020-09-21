package main

import (
	"fmt"
	"image/color"
	"log"
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/audio"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"github.com/hajimehoshi/ebiten/inpututil"
)

var (
	CharWidths = []int{27, 26, 25, 26, 25, 25, 26, 25, 12, 26, 26, 25, 33, 25, 26,
		25, 27, 26, 26, 25, 26, 26, 38, 25, 25, 25}
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
	orbs         []*Orb
	robots       []*Robot
	bolts        []*Bolt
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
		orbs:   make([]*Orb, MaxOrbs),
		robots: make([]*Robot, 0, 10),
		bolts:  make([]*Bolt, 0, 10),
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

	// create Orbs
	for i := 0; i < MaxOrbs; i++ {
		g.orbs[i] = NewOrb(g.level)
	}

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

		// count the enemies in game
		enemyCount := 0
		for _, robot := range g.robots {
			if robot.IsAlive() {
				enemyCount++
			}
		}
		pendingEnemyCount := g.level.PendingEnemies()

		if pendingEnemyCount > 0 && enemyCount < g.level.MaxEnemies() && math.Mod(g.timer, NewEnemyRate) == 0 {
			robotType := g.level.NextEnemy()
			if robotType > RobotNone {
				g.CreateRobot(robotType)
			}
		}

		if pendingEnemyCount+enemyCount > 0 && math.Mod(g.timer, NewFruitRate) == 0 {
			g.CreateFruit(false)
		}

		for _, pop := range g.pops {
			pop.Update()
		}

		for _, fruit := range g.fruits {
			fruit.Update(g)
		}

		for _, bolt := range g.bolts {
			bolt.Update(g)
		}

		for _, orb := range g.orbs {
			if orb.IsActive() {
				orb.Update(g)
			}
		}

		for _, robot := range g.robots {
			if robot.IsAlive() {
				robot.Update(g)
			}
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
		if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
			g.player.StartBlowing(g)
		}
		if inpututil.KeyPressDuration(ebiten.KeySpace) > 1 && inpututil.KeyPressDuration(ebiten.KeySpace) <= MaxBlowingTime {
			g.player.Blowing(g)
		}
		if inpututil.IsKeyJustReleased(ebiten.KeySpace) || inpututil.KeyPressDuration(ebiten.KeySpace) > MaxBlowingTime {
			g.player.StopBlowing(g)
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

		for _, fruit := range g.fruits {
			fruit.Draw(screen, g.timer)
		}

		for _, bolt := range g.bolts {
			bolt.Draw(screen)
		}

		for _, pop := range g.pops {
			pop.Draw(screen)
		}

		for _, robot := range g.robots {
			if robot.alive {
				robot.Draw(screen)
			}
		}

		for _, orb := range g.orbs {
			if orb.IsActive() {
				orb.Draw(screen)
			}
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

func (g *Game) CreateFruit(extra bool) *Fruit {
	// find a free fruit
	for _, fruit := range g.fruits {
		if fruit.TTL == 0 {
			fruit.Generate(extra)
			return fruit
		}
	}
	fruit := NewFruit(g.level, extra)
	g.fruits = append(g.fruits, fruit)
	return fruit
}

func (g *Game) CreateRobot(robotType RobotType) {
	// find a dead robot
	for _, robot := range g.robots {
		if !robot.alive {
			robot.Generate(robotType)
			return
		}
	}
	g.robots = append(g.robots, NewRobot(g.level).Generate(robotType))
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

// NewOrb creates a new orb
func (g *Game) NewOrb() *Orb {
	// assign an inactive Orb
	for _, orb := range g.orbs {
		if !orb.IsActive() {
			return orb.Reset()
		}
	}
	return nil
}

func (g *Game) Fire(directionX, x, y float64) {
	// reuse an existing bolt
	for _, bolt := range g.bolts {
		if !bolt.active {
			bolt.Fire(directionX, x, y)
			return
		}
	}
	// otherwise create a new one
	g.bolts = append(g.bolts, NewBolt(g.level).Fire(directionX, x, y))
}

func (g *Game) displayDebug(screen *ebiten.Image) {
	template := " TPS: %0.2f \n Level %d - Colour %d \n Fruits %d - Pops %d - Orbs %d - Robots %d - Bolts %d \n Player %s \n"
	msg := fmt.Sprintf(template,
		ebiten.CurrentTPS(),
		g.level.id,
		g.level.colour,
		len(g.fruits),
		len(g.pops),
		len(g.orbs),
		len(g.robots),
		len(g.bolts),
		g.player,
	)
	ebitenutil.DebugPrint(screen, msg)

	ebitenutil.DrawLine(screen, 70, 75, 70, 400, color.White)
	ebitenutil.DrawLine(screen, 730, 75, 730, 400, color.White)
	ebitenutil.DrawLine(screen, 70, 75, 730, 75, color.White)
	ebitenutil.DrawLine(screen, 70, 400, 730, 400, color.White)
}

// CharWidth returns width of given character. For characters other than the letters A to Z (i.e. space, and the digits 0 to 9),
// the width of the letter A is returned.
func CharWidth(char byte) int {
	i := int(char) - 65
	if i < 0 {
		i = 0
	}
	if i >= len(CharWidths) {
		log.Printf("character '%c'(%d) not in font", char, char)
	}
	return CharWidths[i]
}

func DrawTextCentre(screen *ebiten.Image, text []byte, y float64) {
	width := 0
	for _, c := range text {
		width += CharWidth(c)
	}
	x := (WindowWidth - width) / 2
	DrawText(screen, text, float64(x), y)
}

func DrawText(screen *ebiten.Image, text []byte, x, y float64) {
	op := &ebiten.DrawImageOptions{}
	for _, char := range text {
		image := images[fmt.Sprintf("font0%d", char)]
		if image == nil {
			log.Printf("character '%c'(%d) not available in font", char, char)
			return
		}
		op.GeoM.Reset()
		op.GeoM.Translate(x, y)
		screen.DrawImage(image, op)
		x += float64(CharWidth(char))
	}
}
