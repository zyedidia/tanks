package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

const (
	// screenWidth  = 1920 / 4
	// screenHeight = 1080 / 4

	screenWidth  = 570
	screenHeight = 400
)

var assets *AssetManager

type Game struct {
	state GameState
}

func (g *Game) Update() error {
	for _, id := range inpututil.JustConnectedGamepadIDs() {
		log.Printf("gamepad connected: id: %d", id)
	}
	// for id := range g.gamepadIDs {
	// 	if inpututil.IsGamepadJustDisconnected(id) {
	// 		log.Printf("gamepad disconnected: id: %d", id)
	// 	}
	// }

	s, err := g.state.Update()
	if s != nil {
		g.state = s
	}
	return err
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.state.Draw(screen)

	// ebitenutil.DebugPrint(screen, fmt.Sprintf("FPS: %0.2f", ebiten.CurrentTPS()))
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowResizable(true)
	ebiten.SetWindowTitle("Turbo Tanks")
	ebiten.SetVsyncEnabled(true)

	assets = LoadAssets()

	g := &Game{
		state: NewMenu(),
	}

	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
