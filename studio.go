package main

type studioState struct {
	prefs preferences
	panel string
	row int
	scroll int
	paused bool
	credit float64
	notice string
	noticeTricks int
	launched int
	exploded int
	peak int
	frames int 
	selected int
	show showPlayer
	tape recording
	recorder recorderState
	replay replayState
}

func newStudio() studioState {
	return studioState{
		prefs: defaultPreferences(),
	}
}


func (m *model) tell(message string) {
	m.studio.notice =message
	m.studio.noticeTricks = fps* 5

}

func (m *model) studioFrozen() bool {
	return  m.studio.paused || m.studio.panel != ""|| m.w < 4 || m.h < 6
}

func (m *model) advanceStudio() {
	s := &m.studio
	s.frames++

	if s.noticeTricks > 0 {
		s.noticeTricks--
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
			m.tell("recording stopped : duration limit reached")
		}
	}
}

func (m *model) clearSky() {
	m.rockets = nil
	m.parts = nil
	m.shooter = nil 
}

