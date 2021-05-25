package input

import "github.com/hajimehoshi/ebiten/v2"

type GamepadButton struct {
	button ebiten.GamepadButton
}

var (
	ButtonA       = GamepadButton{ebiten.GamepadButton0}
	ButtonB       = GamepadButton{ebiten.GamepadButton1}
	ButtonX       = GamepadButton{ebiten.GamepadButton2}
	ButtonY       = GamepadButton{ebiten.GamepadButton3}
	ButtonL1      = GamepadButton{ebiten.GamepadButton4}
	ButtonR1      = GamepadButton{ebiten.GamepadButton5}
	ButtonUp      = GamepadButton{ebiten.GamepadButton11}
	ButtonLeft    = GamepadButton{ebiten.GamepadButton14}
	ButtonRight   = GamepadButton{ebiten.GamepadButton12}
	ButtonDown    = GamepadButton{ebiten.GamepadButton13}
	ButtonStart   = GamepadButton{ebiten.GamepadButton7}
	ButtonOptions = GamepadButton{ebiten.GamepadButton6}
	ButtonHome    = GamepadButton{ebiten.GamepadButton8}
)

const (
	EAxisLAHorizontal = 0
	EAxisLAVertical   = 1
	EAxisRAHorizontal = 3
	EAxisRAVertical   = 4
	EAxisL2           = 2
	EAxisR2           = 5
)

type GamepadAxis struct {
	axis     int
	scale    float64
	shift    float64
	deadzone float64
}

const defaultDeadzone = 0.20

var (
	AxisLAUp    = GamepadAxis{EAxisLAVertical, -1.0, 0.0, defaultDeadzone}
	AxisLADown  = GamepadAxis{EAxisLAVertical, 1.0, 0.0, defaultDeadzone}
	AxisLALeft  = GamepadAxis{EAxisLAHorizontal, -1.0, 0.0, defaultDeadzone}
	AxisLARight = GamepadAxis{EAxisLAHorizontal, 1.0, 0.0, defaultDeadzone}
	AxisRAUp    = GamepadAxis{EAxisRAVertical, -1.0, 0.0, defaultDeadzone}
	AxisRADown  = GamepadAxis{EAxisRAVertical, 1.0, 0.0, defaultDeadzone}
	AxisRALeft  = GamepadAxis{EAxisRAHorizontal, -1.0, 0.0, defaultDeadzone}
	AxisRARight = GamepadAxis{EAxisRAHorizontal, 1.0, 0.0, defaultDeadzone}
	AxisL2      = GamepadAxis{EAxisL2, 0.5, 1.0, defaultDeadzone}
	AxisR2      = GamepadAxis{EAxisR2, 0.5, 1.0, defaultDeadzone}
)

func (g GamepadAxis) isGamepadInput()   {}
func (g GamepadButton) isGamepadInput() {}

type GamepadInput interface {
	isGamepadInput()
}

type Gamepad struct {
	id       ebiten.GamepadID
	bindings map[Action]GamepadInput
}

func NewGamepad(id ebiten.GamepadID, bindings map[Action]GamepadInput) *Gamepad {
	return &Gamepad{
		id:       id,
		bindings: bindings,
	}
}

func (g *Gamepad) Get(a Action) float64 {
	input := g.bindings[a]
	switch t := input.(type) {
	case GamepadAxis:
		value := ebiten.GamepadAxis(g.id, t.axis)
		value += t.shift
		value *= t.scale
		if value >= t.deadzone {
			return value
		}
	case GamepadButton:
		if ebiten.IsGamepadButtonPressed(g.id, t.button) {
			return 1.0
		}
	}
	return 0.0
}

func (g *Gamepad) String() string {
	return ebiten.GamepadName(g.id)
}
