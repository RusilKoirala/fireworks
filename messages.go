package main

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type tickMsg time.Time

type launchMsg struct {
	k    kind
	xf   float64
	yAim float64
}

func tick() tea.Cmd {
	return tea.Tick(time.Second/fps, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func waitLaunch(ch <-chan launchMsg) tea.Cmd {
	return func() tea.Msg {
		return <-ch
	}
}
