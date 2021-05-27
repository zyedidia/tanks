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
	anims  map[string]*Animation
}

func LoadAssets() *AssetManager {
	images := loadImages()
	return &AssetManager{
		images: images,
		sounds: loadSounds(),
		anims:  loadAnimations(images),
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

func loadAnimations(images map[string]*ebiten.Image) map[string]*Animation {
	return map[string]*Animation{
		"ltrack": &Animation{
			loop:     true,
			slowdown: 5,
			img:      images["ltrack-sheet.png"],
			frames:   buildFrames(0, 0, 6, 20, 9),
		},
		"explosion": &Animation{
			loop:     false,
			slowdown: 8,
			img:      images["explosion-6.png"],
			frames:   buildFrames(0, 0, 48, 48, 8),
		},
	}
}

func buildFrames(x, y, w, h, nframes int) []image.Rectangle {
	frames := make([]image.Rectangle, nframes)
	for i := 0; i < nframes; i++ {
		frames[i] = image.Rect(x*w+w*i, y*h, x*w+w*i+w, y*h+h)
	}
	return frames
}
