package main

import (
	"fmt"
	_ "image/png"
	"log"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/audio"
)

// Images
const (
	imageTitle = "title"
	soundLevel = "level0"
)

var (
	images     map[string]*ebiten.Image
	sounds     map[string][]byte
	imageNames = []string{
		imageTitle,
	}
	soundNames = []string{
		soundLevel,
	}
)

func init() {
	// it's easier to build up image names this way
	for i := 0; i <= 9; i++ {
		imageNames = append(imageNames, fmt.Sprintf("space%d", i))
	}
	for i := 0; i <= 3; i++ {
		imageNames = append(imageNames, fmt.Sprintf("bg%d", i))
		imageNames = append(imageNames, fmt.Sprintf("block%d", i))
	}
	for i := 0; i <= 4; i++ {
		for j := 0; j <= 2; j++ {
			imageNames = append(imageNames, fmt.Sprintf("fruit%d%d", i, j))
		}
	}
	for i := 0; i <= 1; i++ {
		for j := 0; j <= 6; j++ {
			imageNames = append(imageNames, fmt.Sprintf("pop%d%d", i, j))
		}
	}
}

func main() {
	var err error

	log.SetFlags(log.LstdFlags | log.Lshortfile)

	images, err = loadImages(imageNames)
	if err != nil {
		log.Fatal(err)
	}

	audioContext, err := audio.NewContext(SampleRate)
	if err != nil {
		log.Fatal(err)
	}

	sounds, err = loadSounds(audioContext, soundNames)
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
