package input

type Action int

type Controller interface {
	Get(a Action) float64
}
