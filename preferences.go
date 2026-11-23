package main

import (
	"fmt"
	"math"
)

const settingCount = 7

type preferences struct {
	Version     int     `json:"version"`
	Speed       float64 `json:"speed"`
	Quality     int     `json:"quality"`
	Gravity     float64 `json:"gravity"`
	StarDensity int     `json:"star_density"`
	Skyline     bool    `json:"skyline"`
	WindLocked  bool    `json:"wind_locked"`
	Wind        float64 `json:"wind"`
}

func defaultPreferences() preferences {
	return preferences{Version: 1, Speed: 1, Quality: 2, Gravity: 1, StarDensity: 1, Skyline: true}
}
func finite(value float64) bool { return !math.IsNaN(value) && !math.IsInf(value, 0) }

func (p preferences) validate() error {
	if p.Version != 1 {
		return fmt.Errorf("unsupported settings version")
	}
	validSpeed := false
	for _, speed := range []float64{0.25, 0.5, 1, 1.5, 2} {
		if p.Speed == speed {
			validSpeed = true
		}
	}
	if !validSpeed {
		return fmt.Errorf("invalid simulation speed")
	}
	if p.Quality < 0 || p.Quality > 2 {
		return fmt.Errorf("quality must be between 0 and 2")
	}
	if !finite(p.Gravity) || p.Gravity < 0.25 || p.Gravity > 2 {
		return fmt.Errorf("gravity must be between 0.25 and 2")
	}
	if p.StarDensity < 0 || p.StarDensity > 3 {
		return fmt.Errorf("star density must be between 0 and 3")
	}
	if !finite(p.Wind) || math.Abs(p.Wind) > 14 {
		return fmt.Errorf("wind must be between -14 and 14")
	}
	return nil
}

func (m *model) adjustSetting(direction int) {
	p := &m.studio.prefs
	rebuild := false
	switch m.studio.row {
	case 0:
		m.setSpeed(direction)
	case 1:
		p.Quality = max(0, min(2, p.Quality+direction))
		if len(m.parts) > m.particleLimit() {
			m.parts = m.parts[:m.particleLimit()]
		}
	case 2:
		p.Gravity = clamp(p.Gravity+float64(direction)*0.25, 0.25, 2)
	case 3:
		p.StarDensity = max(0, min(3, p.StarDensity+direction))
		rebuild = true
	case 4:
		p.Skyline = !p.Skyline
	case 5:
		p.WindLocked = !p.WindLocked
	case 6:
		p.Wind = clamp(p.Wind+float64(direction), -14, 14)
		p.WindLocked = true
	}
	if p.WindLocked {
		m.wind, m.windTarget = p.Wind, p.Wind
	}
	if rebuild {
		m.buildScenery()
	}
}

func (m *model) settingRows() []string {
	p := m.studio.prefs
	quality := [...]string{"Low / 1000", "Medium / 3000", "High / 6000"}
	return []string{fmt.Sprintf("Speed          %.2gx", p.Speed), fmt.Sprintf("Particle cap   %s", quality[p.Quality]), fmt.Sprintf("Gravity        %.2gx", p.Gravity), fmt.Sprintf("Star density   %d", p.StarDensity), fmt.Sprintf("Skyline        %s", onOff(p.Skyline)), fmt.Sprintf("Wind lock      %s", onOff(p.WindLocked)), fmt.Sprintf("Locked wind    %+.0f", p.Wind)}
}

func onOff(value bool) string {
	if value {
		return "on"
	}
	return "off"
}
