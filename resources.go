package main

import (
	"bytes"
	"embed"
	"fmt"
	"image"
	"io/fs"
	"path"
	"strings"

	_ "image/png"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/vorbis"
)

//go:embed images sounds music
var embededFiles embed.FS

func loadImages() (map[string]*ebiten.Image, error) {
	imageNames, err := fs.Glob(embededFiles, "images/*.png")
	if err != nil {
		return nil, err
	}
	imagesMap := make(map[string]*ebiten.Image, len(imageNames))
	for _, imageName := range imageNames {
		file, err := embededFiles.Open(imageName)
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
		imageName = path.Base(imageName)
		imageName = strings.TrimSuffix(imageName, path.Ext(imageName))
		imagesMap[imageName] = img2
	}
	return imagesMap, nil
}

func loadSounds(context *audio.Context) (map[string][]byte, error) {
	soundNames, err := fs.Glob(embededFiles, "sounds/*.ogg")
	if err != nil {
		return nil, err
	}
	soundsMap := make(map[string][]byte, len(soundNames))
	for _, soundName := range soundNames {
		// annoyingly, fs.File does not implement io.ReadSeeker,
		// so we need to load it first and create a reader from the buffer
		buffer, err := embededFiles.ReadFile(soundName)
		if err != nil {
			return soundsMap, fmt.Errorf("%s: %w", soundName, err)
		}
		reader := bytes.NewReader(buffer)
		snd, err := vorbis.Decode(context, reader)
		if err != nil {
			return soundsMap, fmt.Errorf("%s: %w", soundName, err)
		}
		buf := make([]byte, snd.Length())
		_, err = snd.Read(buf)
		if err != nil {
			return soundsMap, fmt.Errorf("%s: %w", soundName, err)
		}
		soundName = path.Base(soundName)
		soundName = strings.TrimSuffix(soundName, path.Ext(soundName))
		soundsMap[soundName] = buf
	}
	return soundsMap, nil
}
