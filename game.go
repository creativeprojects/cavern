package main

import (
	"fmt"
	"log"
	"math"
	"math/rand"

	"github.com/creativeprojects/cavern/lib"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
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
	space        *lib.Sprite
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

	g := &Game{
		audioContext: audioContext,
		musicPlayer:  m,
		state:        StateMenu,
		slow:         false,
		space: lib.NewSprite(lib.XCentre, lib.YCentre).MoveTo(400, 280+45).Animate([]*ebiten.Image{
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
	}

	return g.Initialize(), nil
}

// Initialize a new game
func (g *Game) Initialize() *Game {
	g.timer = -1
	g.level = NewLevel()
	g.level.Next()
	g.fruits = make([]*Fruit, 0, 10)
	g.pops = make([]*Pop, 0, 10)
	g.orbs = make([]*Orb, MaxOrbs)
	g.robots = make([]*Robot, 0, 10)
	g.bolts = make([]*Bolt, 0, 10)

	// create Orbs
	for i := 0; i < MaxOrbs; i++ {
		g.orbs[i] = NewOrb(g.level)
	}
	return g
}

// Layout defines the size of the game in pixels
func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return WindowWidth, WindowHeight
}

// Start a new game
func (g *Game) Start() *Game {
	g.Initialize()
	g.player = NewPlayer().Start(g.level, false)
	g.state = StatePlaying
	return g
}

// NextLevel loads the next level
func (g *Game) NextLevel() {
	g.SoundEffect(sounds[soundLevel])
	g.level.Next()
}

// Update game events
func (g *Game) Update() error {
	g.timer++

	// Debug screen
	if Debug && inpututil.IsKeyJustPressed(ebiten.KeyD) {
		g.debug = !g.debug
	}

	if g.state == StateMenu {
		g.space.Update()

		if len(g.robots) < 4 && math.Mod(g.timer, NewEnemyRate) == 0 {
			robotType := g.level.NextEnemy()
			if robotType > RobotNone {
				g.CreateRobot(robotType)
			}
		}

		if math.Mod(g.timer, NewFruitRate) == 0 {
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
			orb.Update(g)
		}

		for _, robot := range g.robots {
			robot.Update(g)
		}

		if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
			g.Start()
		}
		return nil
	}
	if g.state == StatePlaying {
		// skip to next level
		if Debug && inpututil.IsKeyJustPressed(ebiten.KeyN) {
			g.NextLevel()
		}

		// toggle between slow and normal speed mode
		if Debug && inpututil.IsKeyJustPressed(ebiten.KeyS) {
			g.slow = !g.slow
			if g.slow {
				ebiten.SetMaxTPS(GameSlowSpeed)
			} else {
				ebiten.SetMaxTPS(GameNormalSpeed)
			}
		}

		// instant game over
		if Debug && inpututil.IsKeyJustPressed(ebiten.KeyO) {
			g.state = StateGameOver
		}

		if inpututil.IsKeyJustPressed(ebiten.KeyP) {
			g.state = StatePaused
		}

		// count the enemies in game
		enemyCount := 0
		for _, robot := range g.robots {
			if robot.IsAlive() {
				enemyCount++
			}
		}
		pendingEnemyCount := g.level.PendingEnemies()

		if pendingEnemyCount+enemyCount == 0 {
			// end of the level when all the fruits are gone
			fruitCount := 0
			for _, fruit := range g.fruits {
				if !fruit.HasExpired() {
					fruitCount++
				}
			}
			// also check the trapped enemies
			trapped := 0
			for _, orb := range g.orbs {
				if orb.IsActive() && orb.EnemyTrapped() {
					trapped++
				}
			}
			if fruitCount == 0 && trapped == 0 {
				g.NextLevel()
				return nil
			}
		}
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
			orb.Update(g)
		}

		for _, robot := range g.robots {
			robot.Update(g)
		}

		dx := 0.0
		// player actions
		if g.player.CanMove() {
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
			g.Initialize()
			g.state = StateMenu
		}
		return nil
	}
	return nil
}

// Draw game events
func (g *Game) Draw(screen *ebiten.Image) {

	g.level.Draw(screen)

	for _, fruit := range g.fruits {
		fruit.Draw(screen)
	}

	for _, bolt := range g.bolts {
		bolt.Draw(screen)
	}

	for _, pop := range g.pops {
		pop.Draw(screen)
	}

	for _, robot := range g.robots {
		robot.Draw(screen)
	}

	for _, orb := range g.orbs {
		orb.Draw(screen)
	}

	if g.state == StatePlaying {
		g.player.Draw(screen)
	}

	if g.debug {
		g.displayDebug(screen)
	}

	if g.state == StateMenu {
		screen.DrawImage(images[imageTitle], nil)
		g.space.Draw(screen)
		return
	}

	if g.state == StateGameOver {
		screen.DrawImage(images[imageOver], nil)
		return
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
		if fruit.HasExpired() {
			return fruit.Generate(extra)
		}
	}
	fruit := NewFruit(g.level, extra)
	g.fruits = append(g.fruits, fruit)
	return fruit
}

func (g *Game) CreateRobot(robotType RobotType) {
	// find a dead robot
	for _, robot := range g.robots {
		if !robot.IsAlive() {
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

// Fire generates a new bolt
func (g *Game) Fire(directionX, x, y float64) {
	// reuse an existing bolt
	for _, bolt := range g.bolts {
		if !bolt.IsActive() {
			bolt.Fire(directionX, x, y)
			return
		}
	}
	// otherwise create a new one
	g.bolts = append(g.bolts, NewBolt(g.level).Fire(directionX, x, y))
}

func (g *Game) Orbs() []*Orb {
	return g.orbs
}

func (g *Game) ActiveOrbs() []*Orb {
	orbs := make([]*Orb, 0, len(g.orbs))
	for _, orb := range g.orbs {
		if orb.IsActive() {
			orbs = append(orbs, orb)
		}
	}
	return orbs
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
