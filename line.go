package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/jakecoffman/cp"
)

type Line struct {
	p1, p2 cp.Vector
}

func (l *Line) Draw(screen *ebiten.Image) {
	ebitenutil.DrawLine(screen, l.p1.X, l.p1.Y, l.p2.X, l.p2.Y, color.White)
}
