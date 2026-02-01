package ui

import (
	"context"
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/rangaroo/2gis-friends/internal/core"
)

// bubbletea message and command
type trackerStatus int // TODO: display status in footer

const (
	trackerDisconnected trackerStatus = iota
	trackerConnecting
	trackerConnected
)

type trackerConnectedMsg struct {
	conn *core.WebSocketConn
}

type trackerDataMsg []byte

type trackerEndedMsg struct {
	err error
}

type trackerReconnectMsg struct {}

func (m *Model) updateTracker(msg tea.Msg) (tea.Cmd, bool) {
	switch msg := msg.(type) {
	case trackerReconnectMsg:
		m.trackerStatus = trackerConnecting


	case trackerConnectedMsg:
		m.trackerStatus = trackerConnected

	case trackerDataMsg:


	case trackerEndedMsg:
		m.trackerStatus = trackerDisconnected
	}
}
