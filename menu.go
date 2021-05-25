package main

import "github.com/hajimehoshi/ebiten/v2"

type Menu struct {
}

func NewMenu() *Menu {
	return &Menu{}
}

func (m *Menu) Update() (GameState, error) {
	return nil, nil
}

func (m *Menu) Draw(screen *ebiten.Image) {

}
