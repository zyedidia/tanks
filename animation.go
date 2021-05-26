package main

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

type Animation struct {
	img      *ebiten.Image
	frames   []image.Rectangle
	slowdown int
	loop     bool
}

type AnimImage struct {
	anim  *Animation
	count int
}

func (a *AnimImage) Forward() {
	a.count++
}

func (a *AnimImage) Backward() {
	a.count--
	if a.count < 0 {
		a.count = len(a.anim.frames)*a.anim.slowdown - 1
	}
}

func (a *AnimImage) Frame() int {
	frame := a.count / a.anim.slowdown
	if a.anim.loop {
		return frame % len(a.anim.frames)
	}
	return min(frame, len(a.anim.frames))
}

func (a *AnimImage) Done() bool {
	return !a.anim.loop && a.Frame() >= len(a.anim.frames)
}

func (a *AnimImage) Image() *ebiten.Image {
	return a.anim.img.SubImage(a.anim.frames[a.Frame()]).(*ebiten.Image)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
