package mhandler

type direction string

const (
	left  = "LEFT"
	right = "RIGHT"
	up    = "UP"
	down  = "DOWN"
)

type Motion struct {
	direction direction
	jump      int
}

func NewMotion(direction direction, jump int) *Motion {
	m := Motion{direction: direction, jump: jump}
	return &m
}
