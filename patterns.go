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


func (m *model) burstStar(r rocket, s float64) {
	var vertices [10][2]float64
	
	for i := range vertices {
		angle := -math.Pi/2 + float64(i)*math.Pi/5
		radius := 18 *s

		if i%2 ==  1 {
			radius *= 0.42
		}

		vertices[i] = [2]float64 {
			math.Cos(angle) * radius,
			math.Sin(angle) * radius,
		}
	}

	for edge := range vertices {
		start := vertices[edge]
		end := vertices[(edge+1)%len(vertices)]

		for i := 0; i< 10; i++ {
			t:= float64(i) / 10 
			vx := start[0] + (end[0]-start[0])*t
			vy := start[1] + (end[1]-start[1]) * t

			m.add(Particle{
				x: r.x, y: r.y,
				vx: vx, vy:vy * aspect,
				life: 2.3, pal: r.pal,
				grav: 0.1, drag: 0.7,
			})
		}
	}
}


func (m *model) burstPalm(r rocket, s float64) {
	const fronds , points = 7,12

	for frond :=0; frond < fronds; frond++ {
		angle := -0.9*math.Pi + float64(frond)/float64(fronds-1)*0.8*math.Pi;

		for i := 0; i<points;i++ {
			t:= float64(i) / float64(points-1)
			speed := (0.35+0.65*t) * 19 * s

			m.add(Particle{
				x: r.x, y:r.y,
				vx: math.Cos(angle) * speed,
				vy: math.Sin(angle) * speed * aspect,
				life: 2.6+ rand.Float64()*0.6,
				pal: gold, grav: 0.8, drag: 0.45,
				drips: true,
			})
		}
	}
}