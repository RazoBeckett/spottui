package styles

import (
	"strings"
	"testing"

	"charm.land/lipgloss/v2"
)

func TestApplyForegroundGradColors(t *testing.T) {
	out := ApplyForegroundGrad("ABCDEF", lipgloss.Color("#1DB954"), lipgloss.Color("#1ED760"))
	if !strings.Contains(out, "\x1b[") {
		t.Fatalf("expected ANSI color codes, got %q", out)
	}
	if !strings.Contains(out, "A") || !strings.Contains(out, "F") {
		t.Fatalf("expected original glyphs preserved")
	}
}

func TestGradientLogoPerLine(t *testing.T) {
	logo := "AB\nCD"
	out := GradientLogo(logo, lipgloss.Color("#1DB954"), lipgloss.Color("#1ED760"))
	if strings.Count(out, "\n") != 1 {
		t.Fatalf("expected 1 newline preserved, got %q", out)
	}
}
