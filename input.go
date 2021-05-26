package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/zyedidia/turbotanks/input"
)

const (
	ActionDrive input.Action = iota
	ActionReverse
	ActionRight
	ActionLeft
	ActionShoot
)

var DefaultKeyboard1 = map[input.Action]ebiten.Key{
	ActionDrive:   ebiten.KeyW,
	ActionReverse: ebiten.KeyS,
	ActionRight:   ebiten.KeyD,
	ActionLeft:    ebiten.KeyA,
	ActionShoot:   ebiten.KeyControlLeft,
}

var DefaultKeyboard2 = map[input.Action]ebiten.Key{
	ActionDrive:   ebiten.KeyUp,
	ActionReverse: ebiten.KeyDown,
	ActionRight:   ebiten.KeyRight,
	ActionLeft:    ebiten.KeyLeft,
	ActionShoot:   ebiten.KeyShiftRight,
}

var DefaultGamepad = map[input.Action]input.GamepadInput{
	ActionDrive:   input.AxisLAUp,
	ActionReverse: input.AxisLADown,
	ActionRight:   input.AxisLARight,
	ActionLeft:    input.AxisLALeft,
	ActionShoot:   input.ButtonA,
}
