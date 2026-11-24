package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"
)

func (m *model) studioStatusBar() string {
	if m.w <= 0 {
		return ""
	}

	s := &m.studio
	state := "LIVE"
	if s.paused {
		state = "PAUSED"
	}
	
	if s.recorder.active {
		state += " REC"
	}
	if s.replay.active {
		state += fmt.Sprintf(" REPLAY %d/%d", s.replay.next, len(s.tape.Events))
	}
	if s.show.active {
		state += " " + s.show.script.name
	}
	if m.auto {
		state += " AUTO"
	}
	text := fmt.Sprintf(" %s %.2gx | %d sparks | ? help", state, s.prefs.Speed, len(m.parts))
	if s.noticeTicks > 0 {
		text += " | " + s.notice
	} else if m.w >= 95 {
		text += " | 1-8 launch  p pause  s settings  v select  g show"
	}
	text = ansi.Truncate(text, m.w, "")
	text += strings.Repeat(" ", max(0, m.w-lipgloss.Width(text)))
	return txtStyle.Render(text)
}

func (m *model) studioPanel() string {
	var lines []string
	switch m.studio.panel {
	case "help":
		lines = helpLines()
	case "settings":
		lines = []string{"FIREWORKS / SETTINGS", "", "Up/down selects; left/right changes.", "b saves settings; d restores defaults."}
		for i, row := range m.settingRows() {
			prefix := "  "
			if i == m.studio.row {
				prefix = "> "
			}
			lines = append(lines, prefix+row)
		}
		lines = append(lines, "", "Lower particle caps can help slower terminals.", "Gravity changes falling sparks, not rocket thrust.")
	case "stats":
		lines = m.statisticsLines()
	default:
		return ""
	}
	return m.panelFrame(lines)
}

func (m *model) panelFrame(lines []string) string {
	if m.w <= 0 || m.h <= 0 {
		return ""
	}
	visible := max(1, m.h-2)
	limit := max(0, len(lines)-visible)
	start := min(m.studio.scroll, limit)
	if m.studio.panel == "settings" {
		selected := 4 + m.studio.row
		if selected < start {
			start = selected
		}
		if selected >= start+visible {
			start = selected - visible + 1
		}
	}
	base := lipgloss.NewStyle().Background(lipgloss.Color(skyColor)).Foreground(lipgloss.Color("#cbd5e1"))
	selectedStyle := base.Foreground(lipgloss.Color("#ffe066")).Bold(true)
	var out strings.Builder
	for y := 0; y < m.h; y++ {
		text := ""
		index := start + y
		if y < m.h-2 && index < len(lines) {
			text = lines[index]
		}
		if y == m.h-2 {
			text = "Simulation paused while this panel is open."
			if m.studio.noticeTicks > 0 {
				text = m.studio.notice
			}
		}
		if y == m.h-1 {
			text = "esc back | up/down navigate | q quit"
		}
		style := base
		if strings.HasPrefix(text, "> ") {
			style = selectedStyle
		}
		text = ansi.Truncate(text, m.w, "")
		text += strings.Repeat(" ", max(0, m.w-lipgloss.Width(text)))
		out.WriteString(style.Render(text))
		if y < m.h-1 {
			out.WriteByte('\n')
		}
	}
	return out.String()
}

func helpLines() []string {
	return []string{
		"FIREWORKS / CONTROLS", "", "LAUNCH", "1 classic   2 ring   3 willow   4 heart", "5 spiral    6 double ring   7 star   8 palm", "Space / enter     Launch a random shape", "Left mouse        Launch toward the clicked position", "a                 Toggle automatic launches", "Left / right      Nudge wind and unlock fixed wind", "c                 Clear sky and stop show/replay", "", "PLAYBACK", "p                 Pause / resume", "[ / ]             Slower / faster simulation", "v                 Select the next show", "g                 Start / stop selected show", "r                 Start / finish launch recording", "t                 Start / stop recorded replay", "w                 Save finished recording to disk", "o                 Load recording from disk", "", "PANELS", "s                 Settings", "i                 Session statistics", "?                 Help", "b                 Save current settings", "Esc               Close panel; quit from the sky", "q / ctrl+c        Quit", "", "Recording stores launch timing, shape and position.", "Colors, wind and particle randomness can vary on replay.", "Only one recording is stored; saving replaces it.", "Opening panels freezes simulation time.", "Recordings are limited to 20 minutes / 3000 launches.",
	}
}

func (m *model) statisticsLines() []string {
	s := &m.studio
	lines := []string{"FIREWORKS / SESSION STATISTICS", "", fmt.Sprintf("Simulation time     %.1f seconds", m.clock), fmt.Sprintf("Delivered ticks     %d", s.frames), fmt.Sprintf("Rockets launched    %d", s.launched), fmt.Sprintf("Rockets exploded    %d", s.exploded), fmt.Sprintf("Active rockets      %d", len(m.rockets)), fmt.Sprintf("Active sparks       %d", len(m.parts)), fmt.Sprintf("Peak spark count    %d", s.peak), fmt.Sprintf("Particle cap        %d", m.particleLimit()), fmt.Sprintf("Current wind        %+.2f", m.wind), fmt.Sprintf("Selected show       %s", showNames[s.selected]), "", fmt.Sprintf("Recorded launches   %d", len(s.tape.Events)), fmt.Sprintf("Recording active    %s", onOff(s.recorder.active)), fmt.Sprintf("Replay active       %s", onOff(s.replay.active))}
	if s.show.active {
		lines = append(lines, fmt.Sprintf("Show progress       %.1f / %.1f seconds", s.show.elapsed, s.show.script.duration))
	}
	if s.replay.active {
		lines = append(lines, fmt.Sprintf("Replay progress     %.1f / %.1f seconds", s.replay.elapsed, s.tape.Duration))
	}
	if path, err := studioPath("settings.json"); err == nil {
		lines = append(lines, "", "Settings file:", path)
	}
	if path, err := studioPath("recording.json"); err == nil {
		lines = append(lines, "", "Recording file:", path)
	}
	return lines
}
