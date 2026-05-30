package tui

import (
	"time"

	"github.com/zmb3/spotify/v2"
)

// NavState tracks the active view and back-navigation history.
type NavState struct {
	Current View
	Stack   []View
}

// Push navigates to v, remembering the current view for Pop.
func (n *NavState) Push(v View) {
	n.Stack = append(n.Stack, n.Current)
	n.Current = v
}

// Pop returns to the previously pushed view, if any.
func (n *NavState) Pop() {
	if len(n.Stack) == 0 {
		return
	}
	n.Current = n.Stack[len(n.Stack)-1]
	n.Stack = n.Stack[:len(n.Stack)-1]
}

// Previous peeks at the view Pop would return to, or Current if the stack is empty.
func (n *NavState) Previous() View {
	if len(n.Stack) == 0 {
		return n.Current
	}
	return n.Stack[len(n.Stack)-1]
}

// UIState holds transient, presentation-only state.
type UIState struct {
	Width  int
	Height int

	ErrMsg    string
	ShowError bool

	NotifyMsg  string
	ShowNotify bool

	Fetching     bool
	FetchingDots int
	Searching    bool
}

// PlaybackState holds the now-playing snapshot plus locally tracked
// progress and pending-seek state shared across views.
type PlaybackState struct {
	State *spotify.PlayerState

	LocalProgress  int
	LastProgressAt time.Time
	IsPlaying      bool

	PendingSeek     int
	SeekPending     bool
	LastSeekRequest time.Time
}

// nextRepeatState returns the next mode in the off -> context -> track cycle.
func nextRepeatState(current string) string {
	switch current {
	case "off":
		return "context"
	case "context":
		return "track"
	default:
		return "off"
	}
}

// adjustVolumeOptimistic applies an immediate, clamped volume delta to the
// local playback snapshot so the UI reflects the change before the API responds.
func (m *Model) adjustVolumeOptimistic(delta int) {
	if m.Playback.State == nil || m.Playback.State.Device.ID == "" {
		return
	}
	v := int(m.Playback.State.Device.Volume) + delta
	if v < 0 {
		v = 0
	}
	if v > 100 {
		v = 100
	}
	m.Playback.State.Device.Volume = spotify.Numeric(v)
}
