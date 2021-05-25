package input

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type Keyboard struct {
	bindings map[Action]ebiten.Key
}

func NewKeyboard(bindings map[Action]ebiten.Key) *Keyboard {
	return &Keyboard{
		bindings: bindings,
	}
}

func (k *Keyboard) Get(a Action) float64 {
	key := k.bindings[a]
	if ebiten.IsKeyPressed(key) {
		return 1.0
	}
	return 0.0
}

func (k *Keyboard) String() string {
	return "Keyboard"
}
