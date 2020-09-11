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
)

var (
	images     map[string]*ebiten.Image
	sounds     map[string][]byte
	imageNames = []string{
		imageTitle,
	}
	soundNames = []string{}
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
