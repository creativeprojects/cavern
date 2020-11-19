package main

import (
	"fmt"
	"image"

	_ "image/png"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/vorbis"
	"github.com/markbates/pkger"
)

func loadImages(imageNames []string) (map[string]*ebiten.Image, error) {
	imagesMap := make(map[string]*ebiten.Image, len(imageNames))
	for _, imageName := range imageNames {
		file, err := pkger.Open("/images/" + imageName + ".png")
		if err != nil {
			return imagesMap, fmt.Errorf("%s: %w", imageName, err)
		}
		defer file.Close()
		img, _, err := image.Decode(file)
		if err != nil {
			return imagesMap, fmt.Errorf("%s: %w", imageName, err)
		}
		img2 := ebiten.NewImageFromImage(img)
		if err != nil {
			return imagesMap, fmt.Errorf("%s: %w", imageName, err)
		}
		imagesMap[imageName] = img2
	}
	return imagesMap, nil
}

func loadSounds(context *audio.Context, soundNames []string) (map[string][]byte, error) {
	soundsMap := make(map[string][]byte, len(soundNames))
	for _, soundName := range soundNames {
		file, err := pkger.Open("/sounds/" + soundName + ".ogg")
		if err != nil {
			return soundsMap, fmt.Errorf("%s: %w", soundName, err)
		}
		defer file.Close()
		snd, err := vorbis.Decode(context, file)
		if err != nil {
			return soundsMap, fmt.Errorf("%s: %w", soundName, err)
		}
		buf := make([]byte, snd.Length())
		_, err = snd.Read(buf)
		if err != nil {
			return soundsMap, fmt.Errorf("%s: %w", soundName, err)
		}
		soundsMap[soundName] = buf
	}
	return soundsMap, nil
}
