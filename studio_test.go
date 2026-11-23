package main

import (
	"math"
	"strings"
	"testing"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"
)

func testStudio() *model {
	m := newModel()
	m.w, m.h = 100, 35
	m.buildScenery()
	return m
}

func TestPauseAndSpeed(t *testing.T) {
	m := testStudio()
	m.studio.paused = true
	m.advanceStudio()
	if m.clock != 0 {
		t.Fatal("paused simulation advanced")
	}
	m.studio.paused = false
	m.studio.prefs.Speed = 0.5
	m.advanceStudio()
	if m.clock != 0 {
		t.Fatal("half speed advanced too soon")
	}
	m.advanceStudio()
	if math.Abs(m.clock-dt) > 1e-9 {
		t.Fatal("half speed did not advance after two ticks")
	}
}

func TestRecordAndReplay(t *testing.T) {
	m := testStudio()
	m.toggleRecording()
	m.launch(launchMsg{k: heart, xf: 0.3, yAim: -1})
	for i := 0; i < fps; i++ {
		m.advanceStudio()
	}
	m.launch(launchMsg{k: spiral, xf: 0.7, yAim: 10})
	m.finishRecording()
	if err := m.studio.tape.validate(); err != nil {
		t.Fatal(err)
	}
	before := m.studio.launched
	m.toggleReplay()
	for i := 0; i < fps*3; i++ {
		m.advanceStudio()
	}
	if m.studio.replay.active || m.studio.launched-before != 2 {
		t.Fatal("replay did not produce the recorded launches")
	}
}

func TestShowsFinish(t *testing.T) {
	for index, name := range showNames {
		t.Run(name, func(t *testing.T) {
			m := testStudio()
			m.studio.selected = index
			m.toggleShow()
			expected := len(m.studio.show.script.events)
			for i := 0; i < fps*34; i++ {
				m.advanceStudio()
			}
			if m.studio.show.active || m.studio.launched != expected {
				t.Fatal("show did not finish with all scheduled launches")
			}
		})
	}
}

func TestPanelsFitTerminal(t *testing.T) {
	m := testStudio()
	m.w, m.h = 36, 9
	for _, panel := range []string{"help", "settings", "stats"} {
		m.openPanel(panel)
		rows := strings.Split(m.View(), "\n")
		if len(rows) != m.h {
			t.Fatalf("%s panel has wrong height", panel)
		}
		for _, row := range rows {
			if lipgloss.Width(ansi.Strip(row)) != m.w {
				t.Fatalf("%s panel has wrong width", panel)
			}
		}
	}
}
