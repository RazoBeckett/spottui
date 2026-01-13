package tui

import (
	"regexp"
	"strconv"
	"strings"
)

// SyncedLyricLine represents a single line of synced lyrics with timestamp
type SyncedLyricLine struct {
	TimeMs int
	Text   string
}

var lrcRegex = regexp.MustCompile(`\[(\d{2}):(\d{2})\.(\d{2,3})\](.*)`)

func parseLRC(lrc string) []SyncedLyricLine {
	var lines []SyncedLyricLine
	for _, line := range strings.Split(lrc, "\n") {
		matches := lrcRegex.FindStringSubmatch(line)
		if len(matches) == 5 {
			min, _ := strconv.Atoi(matches[1])
			sec, _ := strconv.Atoi(matches[2])
			msStr := matches[3]
			ms, _ := strconv.Atoi(msStr)
			if len(msStr) == 2 {
				ms *= 10
			}
			timeMs := min*60*1000 + sec*1000 + ms
			text := strings.TrimSpace(matches[4])
			lines = append(lines, SyncedLyricLine{TimeMs: timeMs, Text: text})
		}
	}
	return lines
}
