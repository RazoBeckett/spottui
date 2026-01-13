package tui

import (
	"testing"
)

func TestParseLRC(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []SyncedLyricLine
	}{
		{
			name:  "single line",
			input: "[00:15.50]Hello world",
			expected: []SyncedLyricLine{
				{TimeMs: 15500, Text: "Hello world"},
			},
		},
		{
			name: "multiple lines",
			input: `[00:00.00]First line
[00:05.25]Second line
[01:30.00]Third line`,
			expected: []SyncedLyricLine{
				{TimeMs: 0, Text: "First line"},
				{TimeMs: 5250, Text: "Second line"},
				{TimeMs: 90000, Text: "Third line"},
			},
		},
		{
			name:  "three digit milliseconds",
			input: "[02:45.123]With milliseconds",
			expected: []SyncedLyricLine{
				{TimeMs: 165123, Text: "With milliseconds"},
			},
		},
		{
			name:  "two digit milliseconds",
			input: "[00:30.05]Two digit ms",
			expected: []SyncedLyricLine{
				{TimeMs: 30050, Text: "Two digit ms"},
			},
		},
		{
			name:     "empty input",
			input:    "",
			expected: nil,
		},
		{
			name:     "invalid format",
			input:    "Not a valid LRC line",
			expected: nil,
		},
		{
			name: "mixed valid and invalid",
			input: `[00:10.00]Valid line
Invalid line
[00:20.00]Another valid`,
			expected: []SyncedLyricLine{
				{TimeMs: 10000, Text: "Valid line"},
				{TimeMs: 20000, Text: "Another valid"},
			},
		},
		{
			name:  "line with extra whitespace",
			input: "[00:05.00]  Trimmed text  ",
			expected: []SyncedLyricLine{
				{TimeMs: 5000, Text: "Trimmed text"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseLRC(tt.input)

			if len(result) != len(tt.expected) {
				t.Fatalf("parseLRC() returned %d lines, want %d", len(result), len(tt.expected))
			}

			for i, line := range result {
				if line.TimeMs != tt.expected[i].TimeMs {
					t.Errorf("line %d TimeMs = %d, want %d", i, line.TimeMs, tt.expected[i].TimeMs)
				}
				if line.Text != tt.expected[i].Text {
					t.Errorf("line %d Text = %q, want %q", i, line.Text, tt.expected[i].Text)
				}
			}
		})
	}
}
