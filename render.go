package main

import (
	"math"
	"math/rand"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

type cell struct {
	r    rune
	c    string
	prio int
}

// canvas stuct
type canvas struct {
	w, h  int
	cells []cell
}

// draw a new canvas
func newCanvas(w, h int) *canvas {
	return &canvas{w: w, h: h, cells: make([]cell, w*h)}
}

// to put things
func (cv *canvas) put(fx, fy float64, r rune, c string, prio int) {
	x, y := int(math.Round(fx)), int(math.Round(fy))
	if x < 0 || x >= cv.w || y < 0 || y >= cv.h {
		return
	}
	g := &cv.cells[y*cv.w+x]
	if g.r == 0 || prio <= g.prio {
		*g = cell{r, c, prio}
	}
}

// to style
func (m *model) style(c string) lipgloss.Style {
	s, ok := m.styles[c]
	if !ok {
		s = lipgloss.NewStyle().Background(lipgloss.Color(skyColor))
		if c != "" {
			s = s.Foreground(lipgloss.Color(c))
		}
		m.styles[c] = s
	}
	return s
}

// to view
func (m *model) View() string {
	if m.w == 0 || m.h < 3 {
		return "lighting the fuse…"
	}
	cv := newCanvas(m.w, m.skyH())

	m.drawStars(cv)
	m.drawShooters(cv)
	m.drawParticles(cv)
	m.drawRockets(cv)
	m.drawSkyline(cv)

	return m.paint(cv) + m.statusBar()
}

// draw stars
func (m *model) drawStars(cv *canvas) {
	for _, s := range m.stars {
		x, y := float64(s.x), float64(s.y)
		b := math.Sin(m.clock*s.speed + s.phase)
		switch {
		case b > 0.92:
			cv.put(x, y, '*', "#ffffff", 90)
		case b > 0.5:
			cv.put(x, y, '+', "#ced4da", 91)
		case b > -0.2:
			cv.put(x, y, '.', "#868e96", 92)
		case b > -0.7:
			cv.put(x, y, '.', "#3d434a", 93)
		}
	}
}

// draw the shooters
func (m *model) drawShooters(cv *canvas) {
	for _, s := range m.shooter {
		tail := '-'
		if math.Abs(s.vx) < 3*math.Abs(s.vy) {
			if s.vx > 0 {
				tail = '\\'
			} else {
				tail = '/'
			}
		}
		fade := int((s.age / s.life) * 3) // the whole streak dims near the end
		cv.put(s.x, s.y, '*', tailGray[min(fade, len(tailGray)-1)], 50)
		for i := 1; i < len(tailGray); i++ {
			t := float64(i) * 0.022
			ci := min(i+fade, len(tailGray)-1)
			cv.put(s.x-s.vx*t, s.y-s.vy*t, tail, tailGray[ci], 50+i)
		}
	}
}

func (m *model) drawParticles(cv *canvas) {
	for _, p := range m.parts {
		f := p.age / p.life
		ci := min(int(f*float64(len(p.pal))), len(p.pal)-1)
		if p.flicker && f > 0.45 && rand.Float64() < 0.4 {
			continue
		}
		var r rune
		switch {
		case p.ember && f < 0.4:
			r = '·'
		case p.ember:
			r = '.'
		case f < 0.35:
			r = '*'
		case f < 0.7:
			r = '+'
		default:
			r = '.'
		}
		prio := ci * 2
		if p.ember {
			prio++
		}
		cv.put(p.x, p.y, r, p.pal[ci], prio)
	}
}

func (m *model) drawRockets(cv *canvas) {
	for _, r := range m.rockets {
		cv.put(r.x, r.y, '|', "#fff3bf", 0)
	}
}

func (m *model) drawSkyline(cv *canvas) {
	for x := 0; x < cv.w && x < len(m.skyline); x++ {
		for h := 0; h < m.skyline[x]; h++ {
			y := cv.h - 1 - h
			if y < 0 {
				break
			}
			c := "#1a1d2e"
			if m.windowLit(x, y, h) {
				c = "#e8b923"
			}
			cv.cells[y*cv.w+x] = cell{'█', c, -1}
		}
	}
}

func (m *model) paint(cv *canvas) string {
	var b strings.Builder
	run := make([]rune, 0, cv.w)
	for y := 0; y < cv.h; y++ {
		cur := ""
		flush := func() {
			if len(run) > 0 {
				b.WriteString(m.style(cur).Render(string(run)))
				run = run[:0]
			}
		}
		for x := 0; x < cv.w; x++ {
			c := cv.cells[y*cv.w+x]
			if c.r == 0 {
				run = append(run, ' ')
				continue
			}
			if c.c != cur {
				flush()
				cur = c.c
			}
			run = append(run, c.r)
		}
		flush()
		b.WriteByte('\n')
	}
	return b.String()
}
