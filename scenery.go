package main

import (
	"math"
	"math/rand"
)

// to build secenery
func (m *model) buildScenery() {
	sky := m.skyH()

	if m.w <= 0 || sky <= 0 {
		return
	}

	n := m.w * sky / 32
	m.stars = make([]star, n)
	for i := range m.stars {
		m.stars[i] = star{
			x: rand.Intn(m.w), y: rand.Intn(max(1, sky*3/4)), phase: rand.Float64() * 2 * math.Pi, speed: 0.4 + rand.Float64()*2.5,
		}
	}

	maxH := int(clamp(float64(sky)/7, 2, 7))
	m.skyline = make([]int, m.w)
	m.edge = make([]bool, m.w)
	for x := 0; x < m.w; {
		bw := 4 + rand.Intn(7)
		bh := 1 + rand.Intn(maxH)
		if rand.Float64() < 0.15 {
			bh, bw = 0, 1+rand.Intn(3)
		}
		for i := 0; i < bw && x+i < m.w; i++ {
			m.skyline[x+i] = bh
			m.edge[x+i] = i == 0 || i == bw-1
		}
		x += bw
	}
}

// to make those yellow windwss
func (m *model) windowLit(x, y, h int) bool {
	return !m.edge[x] && h%2 == 1 && x%2 == 0 && hash(x, y)%3 == 0
}
