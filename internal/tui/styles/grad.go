package styles

import (
	"image/color"
	"strings"

	"charm.land/lipgloss/v2"
	"github.com/rivo/uniseg"
)

// ForegroundGrad renders input with a horizontal foreground gradient from
// color1 to color2, returning one styled string per grapheme cluster.
// Adapted from charmbracelet/crush (internal/ui/styles/grad.go).
func ForegroundGrad(input string, color1, color2 color.Color) []string {
	if input == "" {
		return []string{""}
	}

	var clusters []string
	gr := uniseg.NewGraphemes(input)
	for gr.Next() {
		clusters = append(clusters, string(gr.Runes()))
	}

	ramp := lipgloss.Blend1D(len(clusters), color1, color2)
	for i, c := range ramp {
		clusters[i] = lipgloss.NewStyle().Foreground(c).Render(clusters[i])
	}
	return clusters
}

// ApplyForegroundGrad renders a string with a horizontal foreground gradient.
func ApplyForegroundGrad(input string, color1, color2 color.Color) string {
	if input == "" {
		return ""
	}
	return strings.Join(ForegroundGrad(input, color1, color2), "")
}

// GradientLogo renders multi-line ASCII art with a per-line horizontal
// gradient between two colors.
func GradientLogo(art string, color1, color2 color.Color) string {
	lines := strings.Split(art, "\n")
	for i, line := range lines {
		lines[i] = ApplyForegroundGrad(line, color1, color2)
	}
	return strings.Join(lines, "\n")
}
