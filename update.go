package main

import (
	"math/rand"

	tea "github.com/charmbracelet/bubbletea"
)

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.w, m.h = msg.Width, msg.Height
		m.buildScenery()
		return m, nil

	case tea.KeyMsg:
		return m.handleKey(msg)

	case tea.MouseMsg:
		if msg.Action == tea.MouseActionPress && msg.Button == tea.MouseButtonLeft && m.w > 0 {
			m.launch(launchMsg{
				k:    randomKind(),
				xf:   float64(msg.X) / float64(m.w),
				yAim: float64(msg.Y),
			})
		}
		return m, nil

	case launchMsg:
		m.launch(msg)
		return m, waitLaunch(m.launches)

	case tickMsg:
		m.step()
		return m, tick()
	}
	return m, nil
}

func (m *model) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q", "esc", "ctrl+c":
		m.setAuto(false)
		return m, tea.Quit
	case "1", "2", "3", "4":
		m.launch(launchMsg{k: kind(msg.String()[0] - '1'), xf: randX(), yAim: -1})
	case " ", "space", "enter":
		m.launch(launchMsg{k: randomKind(), xf: randX(), yAim: -1})
	case "a":
		m.setAuto(!m.auto)
	case "left", "h":
		m.nudgeWind(-2)
	case "right", "l":
		m.nudgeWind(+2)
	case "c":
		m.rockets, m.parts, m.shooter = nil, nil, nil
	}
	return m, nil
}

func (m *model) nudgeWind(d float64) {
	m.windTarget = clamp(m.windTarget+d, -14, 14)
	m.windTimer = 10
}

func randomKind() kind {
	return kind(rand.Intn(int(numKinds)))
}

func randX() float64 {
	return 0.15 + rand.Float64()*0.7
}
