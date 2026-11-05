package main

import (
	"math"
	"math/rand"
)

func (m *model) launch(l launchMsg) {
	if m.w == 0 {
		return
	}

	sky := float64(m.skyH())
	startY := sky - 1
	target := l.yAim
	if target < 0 {
		target = sky * (0.12 * rand.Float64() * 0.33)
	}
	target = clamp(target, 2, sky*0.75)

	vy := -math.Sqrt(2 * rocketG * (startY - target))
	x := clamp(l.xf, 0.02, 0.98) * float64(m.w)

	m.rockets = append(m.rockets, rocket{
		x: x, y: startY, vx: rand.Float64()*4 - 2, vy: vy,
		targetY: target, k: l.k, pal: randPal(),
	})
}

func (m *model) add(p particle) {
	if len(m.parts) < maxParticles {
		m.parts = append(m.parts, p)
	}
}

func (m *model) explode(r rocket) {

	s := clamp(float64(m.skyH())/26, 0.6, 1.7)

	switch r.k {
	case classic:
		m.burstClassic(r, s)
	case ring:
		m.burstRing(r, s)
	case willow:
		m.burstWillow(r, s)
	case heart:
		m.burstHeart(r, s)
	}
}

func (m *model) burstClassic(r rocket, s float64) {
	n := 80 + rand.Intn(50)
	alt := randPal()
	mixed := rand.Float64() < 0.35
	for i := 0; i < n; i++ {
		a := rand.Float64() * 2 * math.Pi
		sp := (0.25 + 0.75*math.Sqrt(rand.Float64())) * 17 * s
		pal := r.pal
		if mixed && i%2 == 0 {
			pal = alt
		}
		m.add(particle{
			x: r.x, y: r.y,
			vx: math.Cos(a) * sp, vy: math.Sin(a) * sp * aspect,
			life: 1.2 + rand.Float64()*1.0, pal: pal,
			grav: 1, drag: 1.0, flicker: rand.Float64() < 0.3,
		})
	}
}

func (m *model) burstRing(r rocket, s float64) {
	n := 56
	sp := 15 * s
	offset := rand.Float64() * 2 * math.Pi
	ringOf := func(count int, speed, life float64, pal palette) {
		for i := 0; i < count; i++ {
			a := offset + float64(i)/float64(count)*2*math.Pi
			m.add(particle{
				x: r.x, y: r.y,
				vx: math.Cos(a) * speed, vy: math.Sin(a) * speed * aspect,
				life: life, pal: pal, grav: 0.8, drag: 1.0,
			})
		}
	}
	ringOf(n, sp, 1.7+rand.Float64()*0.2, r.pal)
	if rand.Float64() < 0.6 { // inner ring in a contrasting colour
		ringOf(n/2, sp*0.5, 1.5, randPal())
	}
}

func (m *model) burstWillow(r rocket, s float64) {
	for i := 0; i < 70; i++ {
		a := rand.Float64() * 2 * math.Pi
		sp := (0.4 + 0.6*rand.Float64()) * 12 * s
		m.add(particle{
			x: r.x, y: r.y,
			vx: math.Cos(a) * sp, vy: math.Sin(a) * sp * aspect,
			life: 3.0 + rand.Float64()*1.6, pal: gold,
			grav: 1.1, drag: 1.5, drips: true,
		})
	}
}

func (m *model) burstHeart(r rocket, s float64) {
	n := 80
	k := 0.95 * s
	pal := heartPals[rand.Intn(len(heartPals))]
	for i := 0; i < n; i++ {
		t := float64(i) / float64(n) * 2 * math.Pi
		hx := 16 * math.Pow(math.Sin(t), 3)
		hy := 13*math.Cos(t) - 5*math.Cos(2*t) - 2*math.Cos(3*t) - math.Cos(4*t)
		m.add(particle{
			x: r.x, y: r.y,
			vx: hx * k, vy: -hy * k * aspect,

			life: 2.0, pal: pal, grav: 0.3, drag: 1.2,
		})
	}
}
