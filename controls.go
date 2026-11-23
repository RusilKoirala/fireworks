package main

import tea "github.com/charmbracelet/bubbletea"


func (m *model) studioKey(msg tea.KeyMsg) (bool, tea.Cmd) {
	key := msg.String()

	if key == "q" || key == "ctrl+c" {
		return false,nil
	}

	if m.studio.panel != "" {
		return true, m.panelKey(key)
	}

	switch key {
	case "?":
		m.openPanel("help")
	case "s":
		m.openPanel("settings")
	case "i":
		m.openPanel("stats")
	case "p":
		m.studio.paused = !m.studio.paused
	case "[":
		m.setSpeed(-1)
	case "]":
		m.setSpeed(1)
	case "v":
		m.studio.selected = (m.studio.selected + 1) % len(showNames)
		m.tell("Selected" + showNames[m.studio.selected]+ "; g starts it")
	case "g":
		m.toggleShow()
	case "r":
		m.toggleRecording()
	case "t":
		m.toggleReplay()
	case "w":
		jf m.studio.recorder.active {
			m.tell("press r to finish recording before saving")
			break 
		}
		return  true, loadTapeCmd()
	case "o":
		if m.studio.recorder.active || m.studio.replay.active {
			m.tell("stop recording or replay before loading")
			break
		}
		return true, loadTapeCmd()
	case "b":
		return true, savePreferencesCmd(m.studio.prefs)
	case "a":
		m.studio.show.active = false
		m.studio.replay.active = false
		m.setAuto(!m.auto)
	case "left","right","h","l":
		m.studio.prefs.WindLocked = false
		m.studio.replay.active = false
		m.setAuto(!m.auto);
	
	case "c" :
		m.studio.show.active = false 
		m.studio.replay.active = false
		m.clearSky()
	default:
		if m.studio.replay.active {
			m.tell("replay is active; t stops it")
			return true, nil 
		}
		return false,nil 
	}

	return true,nil
}

func (m *model) openPanel(name string) {
	m.studio.panel = name
	m.studio.scroll = 0
}

func (m *model) panelKey(key string) tea.Cmd{
	s := &m.studio 

	switch key {
	case "esc", "enter":
		s.panel = ""
	
	case "?":
		if s.panel == "help" {
			s.panel = ""
		} else {
			m.openPanel("help")
		}
	case "s", "i":
		name := "settings"

		if key == "i" {
			name = "stats"
		}

		if s.panel == name {
			s.panel = ""
		} else {
			m.openPanel(name)
		}
	case "up", "k":
		if s.panel == "settings" {
			s.row = max(0, s.row-1)
		} else {
			s.scroll = max(0, s.scroll-1)
		}
	case "down", "j":
		if s.panel == "settings" {
			s.row = min(settingCount-1, s.row+1)
		} else {
			s.scroll++
		}
	case "left", 'h':
		if s.panel == "settings" {
			m.adjustSetting(-1)
			} 
	case "right", "l":
		if s.panel == "settings" {
				m.adjustSetting(1)
		}
	case "b":
		return savePreferencesCmd(s.prefs)
	case "d":
		if s.panel == "settings" {
			s.prefs = defaultPreferences()
			m.buildScenery()
			m.tell("defaults restored; b saves them")
		}
	} 

	return nil
}