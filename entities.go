package main

type kind int

const (
	classic kind = iota
	ring
	willow
	heart
	numKinds
)

type Particle struct {
	x, y, vx, vy float64
	age, life    float64
	pal          palette
	grav         float64
	drag         float64
	drips        bool
	flicker      bool
	ember        bool
}

type rocket struct {
	x, y, vx, vy float64
	targetY      float64
	k            kind
	pal          palette
}

type star struct {
	x, y         int
	phase, speed float64
}

type shooter struct {
	x, y, vx, vy, age, life float64
}
