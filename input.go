package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/zyedidia/turbotanks/input"
)

const (
	ActionLeftDrive input.Action = iota
	ActionLeftReverse
	ActionRightDrive
	ActionRightReverse
	ActionShoot
	ActionReload
)

var DefaultKeyboard = map[input.Action]ebiten.Key{
	ActionLeftDrive:    ebiten.KeyW,
	ActionLeftReverse:  ebiten.KeyS,
	ActionRightDrive:   ebiten.KeyI,
	ActionRightReverse: ebiten.KeyK,
	ActionShoot:        ebiten.KeySpace,
	ActionReload:       ebiten.KeyM,
}

var DefaultGamepad = map[input.Action]input.GamepadInput{
	ActionLeftDrive:    input.AxisLAUp,
	ActionLeftReverse:  input.AxisLADown,
	ActionRightDrive:   input.AxisRAUp,
	ActionRightReverse: input.AxisRADown,
	ActionShoot:        input.ButtonR1,
	ActionReload:       input.ButtonL1,
}
