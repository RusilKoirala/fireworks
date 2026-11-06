package main

import (
	"fmt"
	"math"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// colors
var (
	barBg    = lipgloss.Color("#161a2b")
	keyStyle = lipgloss.NewStyle().Background(barBg).Foreground(lipgloss.Color("#ffe066")).Bold(true)
	txtStyle = lipgloss.NewStyle().Background(barBg).Foreground(lipgloss.Color("#a5aec0"))
	autoOff  = lipgloss.NewStyle().Background(lipgloss.Color("#343a40")).Foreground(lipgloss.Color("#adb5bd")).Padding(0, 1)
	autoOn   = lipgloss.NewStyle().Background(lipgloss.Color("#2f9e44")).Foreground(lipgloss.Color("#ffffff")).Bold(true).Padding(0, 1)
)

// status bar for info
func (m *model) statusBar() string {
	sep := txtStyle.Render("  ")
	item := func(k, label string) string { return keyStyle.Render(k) + txtStyle.Render(" "+label) }

	badge := autoOff.Render("AUTO off")
	if m.auto {
		badge = autoOn.Render("AUTO on")
	}

	line := strings.Join([]string{
		txtStyle.Render(" ") + item("1", "burst"),
		item("2", "ring"),
		item("3", "willow"),
		item("4", "heart"),
		item("spc", "random"),
		keyStyle.Render("a") + txtStyle.Render(" ") + badge,
		item("←→", "wind "+windArrows(m.wind)),
		item("c", "clear"),
		item("q", "quit"),
	}, sep)

	if count := sep + txtStyle.Render(fmt.Sprintf("%d sparks", len(m.parts))); lipgloss.Width(line+count) <= m.w {
		line += count
	}

	if pad := m.w - lipgloss.Width(line); pad > 0 {
		line += txtStyle.Render(strings.Repeat(" ", pad))
	}
	return lipgloss.NewStyle().MaxWidth(m.w).Render(line)
}

// windows arrow func
func windArrows(wind float64) string {
	n := int(math.Min(4, math.Round(math.Abs(wind)/2.5)))
	switch {
	case n == 0:
		return "calm"
	case wind > 0:
		return strings.Repeat("→", n)
	default:
		return strings.Repeat("←", n)
	}
}
