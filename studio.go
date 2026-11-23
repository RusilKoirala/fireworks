package main

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

type studioState struct {
	prefs       preferences
	panel       string
	row         int
	scroll      int
	paused      bool
	credit      float64
	notice      string
	noticeTicks int
	launched    int
	exploded    int
	peak        int
	frames      int
	selected    int
	show        showPlayer
	tape        recording
	recorder    recorderState
	replay      replayState
}

func newStudio() studioState { return studioState{prefs: defaultPreferences()} }

func (m *model) tell(message string) {
	m.studio.notice = message
	m.studio.noticeTicks = fps * 5
}

func (m *model) studioFrozen() bool {
	return m.studio.paused || m.studio.panel != "" || m.w < 4 || m.h < 6
}

func (m *model) advanceStudio() {
	s := &m.studio
	s.frames++
	if s.noticeTicks > 0 {
		s.noticeTicks--
	}
	if m.studioFrozen() {
		return
	}
	s.credit += s.prefs.Speed
	for s.credit >= 1 {
		s.credit--
		m.step()
		m.advanceShow()
		m.advanceReplay()
		if len(m.parts) > s.peak {
			s.peak = len(m.parts)
		}
		if s.recorder.active && m.clock-s.recorder.started >= maxTapeSeconds {
			m.finishRecording()
			m.tell("Recording stopped: duration limit reached")
		}
	}
}

func (m *model) clearSky() { m.rockets, m.parts, m.shooter = nil, nil, nil }

func (m *model) particleLimit() int {
	limits := [...]int{1000, 3000, maxParticles}
	return limits[m.studio.prefs.Quality]
}

func (m *model) setSpeed(direction int) {
	speeds := [...]float64{0.25, 0.5, 1, 1.5, 2}
	index := 2
	for i, speed := range speeds {
		if speed == m.studio.prefs.Speed {
			index = i
		}
	}
	index = max(0, min(len(speeds)-1, index+direction))
	m.studio.prefs.Speed = speeds[index]
	m.tell(fmt.Sprintf("Simulation speed: %.2gx", speeds[index]))
}

func (m *model) handleStudioMessage(msg tea.Msg) (tea.Cmd, bool) {
	switch msg := msg.(type) {
	case preferencesLoadedMsg:
		if msg.err != nil {
			m.tell("Settings: " + msg.err.Error())
		} else if !msg.missing {
			m.studio.prefs = msg.prefs
			m.buildScenery()
			m.tell("Saved settings loaded")
		}
		return nil, true
	case tapeLoadedMsg:
		if msg.err != nil {
			m.tell("Load failed: " + msg.err.Error())
		} else if m.studio.recorder.active || m.studio.replay.active {
			m.tell("Stop recording or replay before loading")
		} else {
			m.studio.tape = msg.tape
			m.tell(fmt.Sprintf("Loaded %d launches; press t to replay", len(msg.tape.Events)))
		}
		return nil, true
	case fileSavedMsg:
		if msg.err != nil {
			m.tell("Save failed: " + msg.err.Error())
		} else {
			m.tell("Saved " + msg.name)
		}
		return nil, true
	}
	return nil, false
}
