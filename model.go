package main

import (
	"math/rand"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	w, h int

	rockets []rocket
	parts   []Particle
	stars   []star
	shooter []shooter
	skyline []int
	edge    []bool

	wind, windTarget, windTimer float64

	clock float64

	auto     bool
	stopAuto chan struct{}
	launches chan launchMsg

	styles map[string]lipgloss.Style
}

func newModel() *model {
	return &model{
		launches:   make(chan launchMsg, 32),
		styles:     map[string]lipgloss.Style{},
		windTarget: rand.Float64()*4 - 2,
		windTimer:  5,
	}
}
func (m *model) skyH() int {
	return m.h - 1
}

func (m *model) Init() tea.Cmd {
	return tea.Batch(tick(), waitLaunch((m.launches)))
}
