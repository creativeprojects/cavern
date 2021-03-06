package main

import (
	_ "image/png"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
)

// Images
const (
	imageTitle       = "title"
	imagePlayerBlank = "blank"
	imagePlayerStill = "still"
	imageJumpLeft    = "jump0"
	imageJumpRight   = "jump1"
	imageBlowLeft    = "blow0"
	imageBlowRight   = "blow1"
	imageRecoilLeft  = "recoil0"
	imageRecoilRight = "recoil1"
	imageLife        = "life"
	imageHealth      = "health"
	imagePlus        = "plus"
	imageOver        = "over"

	soundLevel = "level0"
	soundJump  = "jump0"
	soundScore = "score0"
	soundBonus = "bonus0"
	soundLife  = "life0"
	soundOver  = "over0"
)

var (
	images     map[string]*ebiten.Image
	sounds     map[string][]byte
	imageNames = []string{
		imageTitle,
		imagePlayerBlank,
		imagePlayerStill,
		imageJumpLeft,
		imageJumpRight,
		imageBlowLeft,
		imageBlowRight,
		imageRecoilLeft,
		imageRecoilRight,
		imageLife,
		imageHealth,
		imagePlus,
		imageOver,
	}
	soundNames = []string{
		soundLevel,
		soundJump,
		soundScore,
		soundBonus,
		soundLife,
		soundOver,
	}
)

func main() {
	var err error

	log.SetFlags(log.LstdFlags | log.Lshortfile)

	images, err = loadImages()
	if err != nil {
		log.Fatal(err)
	}

	audioContext := audio.NewContext(SampleRate)

	sounds, err = loadSounds(audioContext)
	if err != nil {
		log.Fatal(err)
	}

	ebiten.SetRunnableOnUnfocused(true)
	ebiten.SetWindowSize(WindowWidth, WindowHeight)
	ebiten.SetWindowTitle(WindowTitle)
	game, err := NewGame(audioContext)
	if err != nil {
		log.Fatal(err)
	}
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
