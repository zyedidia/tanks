package main

import "github.com/hajimehoshi/ebiten/v2"

type Explosion struct {
	img  *AnimImage
	x, y float64
}

func (e *Explosion) Update() {
	e.img.Forward()
}

func (e *Explosion) Draw(screen *ebiten.Image) {
	if e.img.Done() {
		return
	}

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(-48/2, -48/2)
	op.GeoM.Translate(e.x, e.y)
	screen.DrawImage(e.img.Image(), op)
}
