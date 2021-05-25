package main

import (
	"embed"
	"errors"
	"image"
	_ "image/png"
	"io"
	"io/fs"
	"log"
	"path/filepath"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/vorbis"
)

//go:embed assets/img
var images embed.FS

//go:embed assets/sound
var sounds embed.FS

type AssetManager struct {
	images map[string]*ebiten.Image
	sounds map[string]*audio.Player
}

func LoadAssets() *AssetManager {
	return &AssetManager{
		images: loadImages(),
		sounds: loadSounds(),
	}
}

// TODO: error reporting
func loadImages() map[string]*ebiten.Image {
	loaded := make(map[string]*ebiten.Image)
	fs.WalkDir(images, ".", func(path string, info fs.DirEntry, err error) error {
		if info.IsDir() || filepath.Ext(path) != ".png" {
			return nil
		}

		f, err := images.Open(path)
		if err != nil {
			return err
		}
		img, _, err := image.Decode(f)
		if err != nil {
			return err
		}
		loaded[info.Name()] = ebiten.NewImageFromImage(img)
		log.Println("Loaded", path)

		return f.Close()
	})

	return loaded
}

func loadSounds() map[string]*audio.Player {
	loaded := make(map[string]*audio.Player)
	fs.WalkDir(sounds, ".", func(path string, info fs.DirEntry, err error) error {
		if info.IsDir() || filepath.Ext(path) != ".ogg" {
			return nil
		}

		f, err := sounds.Open(path)
		if err != nil {
			return err
		}
		fseek, ok := f.(io.ReadSeeker)
		if !ok {
			return errors.New("File is not a ReadSeeker")
		}
		s, err := vorbis.Decode(audioctx, fseek)
		if err != nil {
			return err
		}
		player, err := audio.NewPlayer(audioctx, s)
		if err != nil {
			return err
		}
		loaded[info.Name()] = player
		log.Println("Loaded", path)

		return f.Close()
	})
	loaded["explode.ogg"].SetVolume(0.5)

	return loaded
}
