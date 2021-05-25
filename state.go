package main

import "github.com/hajimehoshi/ebiten/v2"

type GameState interface {
	Update() (GameState, error)
	Draw(s *ebiten.Image)
}
