package main

import (
	"math"
	"math/rand"
)



func (m *model) burstSpiral(r rocket, s float64) {
	const arms, points = 3, 30
	offset :=  rand.Float64() * 2 * math.Pi

	for arm:= 0; arm< arms; arm++ {
		for i := 0; i < points; i++ {
			t := float64(i)/float64(points-1)
			angle := offset + float64(arm) * 2 * math.Pi/arms + t * 1.5* math.Pi
			speed := (0.15 + 0.85 *t ) * 18 *s;

			m.add(Particle{
				x: r.x, y:r.y,
				vx: math.Cos(angle) * speed,
				vy: math.Sin(angle) * speed *aspect,
				life: 2.2, pal: r.pal,
				grav: 0.15, drag:0.7, 
			})
		}
	}
}


func (m *model) burstDoubleRing(r rocket, s float64) {
	offset := rand.Float64() * 2 * math.Pi;
	innerPalette := randPal()
	
	for band := 0;band < 2; band++ {
		count := 64
		speed := 17 * s
		colors := r.pal

		if band == 1 {
			count = 40
			speed = 9 * s
			colors = innerPalette
		}

		for i := 0; i < count; i++ {
			angle := offset + float64(i)/float64(count)*2*math.Pi

			m.add(Particle{
				x: r.x, y: r.y,
				vx: math.Cos(angle) * speed,
				vy: math.Sin(angle) * speed * aspect,
				life: 2.1, pal:colors,
				grav: 0.2, drag: 0.8,
			})
		}
	}
}