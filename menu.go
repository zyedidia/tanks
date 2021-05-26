package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type Menu struct {
}

func NewMenu() *Menu {
	return &Menu{}
}

func (m *Menu) Update() (GameState, error) {
	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		return NewMatch(), nil
	}

	return nil, nil
}

func (m *Menu) Draw(screen *ebiten.Image) {
	ebitenutil.DebugPrintAt(screen, "Turbo Tanks", screenWidth/2-30, screenHeight/2-3)
	ebitenutil.DebugPrintAt(screen, "Press Enter to Play", screenWidth/2-50, screenHeight/2+20)
}
