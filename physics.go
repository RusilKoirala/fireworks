package main

import "math/rand"

// step advances the whole simulation by one frame.
func (m *model) step() {
	if m.w == 0 {
		return
	}
	m.clock += dt
	m.updateWind()
	m.updateRockets()
	m.updateParticles()
	m.updateShooters()
}

// updateWind lets the wind wander slowly toward a new target every few seconds.
func (m *model) updateWind() {
	m.windTimer -= dt
	if m.windTimer <= 0 {
		m.windTarget = rand.Float64()*8 - 4
		m.windTimer = 4 + rand.Float64()*6
	}
	m.wind += (m.windTarget - m.wind) * 0.6 * dt
}

// updateRockets: rockets climb, shed sparks, and burst near their apex.
func (m *model) updateRockets() {
	alive := m.rockets[:0]
	for _, r := range m.rockets {
		r.vy += rocketG * dt
		r.x += (r.vx + m.wind*0.2) * dt
		r.y += r.vy * dt
		if rand.Float64() < 0.8 {
			m.add(Particle{
				x: r.x, y: r.y + 0.5,
				vx: rand.Float64()*2 - 1, vy: 1 + rand.Float64()*2,
				life: 0.3 + rand.Float64()*0.4, pal: emberPalette,
				grav: 0.4, drag: 2, ember: true,
			})
		}
		if r.vy >= -1.5 || r.y <= r.targetY {
			m.explode(r)
		} else {
			alive = append(alive, r)
		}
	}
	m.rockets = alive
}

// updateParticles: gravity pulls sparks down, air drag pulls their velocity
// toward the wind (so everything drifts together), and age fades them out.
func (m *model) updateParticles() {
	next := m.parts[:0]
	var drips []Particle
	W, H := float64(m.w), float64(m.h)
	for _, p := range m.parts {
		p.vx += (m.wind - p.vx) * p.drag * dt
		p.vy += gravity*p.grav*dt - p.vy*p.drag*dt
		p.x += p.vx * dt
		p.y += p.vy * dt
		p.age += dt

		if p.drips && rand.Float64() < 0.35 {
			drips = append(drips, Particle{
				x: p.x, y: p.y, vx: p.vx * 0.3, vy: 0.5,
				life: 0.6 + rand.Float64()*0.7, pal: emberPalette,
				grav: 0.7, drag: 2, ember: true,
			})
		}
		if p.age < p.life && p.x > -4 && p.x < W+4 && p.y < H+2 {
			next = append(next, p)
		}
	}
	m.parts = next
	for _, d := range drips {
		m.add(d)
	}
}

// updateShooters spawns rare shooting stars and moves existing ones.
func (m *model) updateShooters() {
	if rand.Float64() < 0.006 {
		dir := 1.0
		if rand.Intn(2) == 0 {
			dir = -1
		}
		m.shooter = append(m.shooter, shooter{
			x: rand.Float64() * float64(m.w), y: rand.Float64() * float64(m.skyH()) * 0.3,
			vx: dir * (35 + rand.Float64()*30), vy: 6 + rand.Float64()*6,
			life: 0.7 + rand.Float64()*0.5,
		})
	}
	alive := m.shooter[:0]
	for _, s := range m.shooter {
		s.x += s.vx * dt
		s.y += s.vy * dt
		s.age += dt
		if s.age < s.life {
			alive = append(alive, s)
		}
	}
	m.shooter = alive
}
