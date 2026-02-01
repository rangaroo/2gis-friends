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

func connectToTrackerCmd(cfg core.Config) tea.Cmd {
	
}

func (m *Model) updateTracker(msg tea.Msg) (tea.Cmd, bool) {
	switch msg := msg.(type) {
	case trackerReconnectMsg:
		m.trackerStatus = trackerConnecting
		return connectToTrackerCmd(m.cfg), true

	case trackerConnectedMsg:
		m.trackerStatus = trackerConnected
		m.reconnectTimeout= 0

		return tea.Batch(
			produceMessagesCmd(msg.conn, m.sub),
			waitForMessages(m.sub),
		), true

	case trackerDataMsg:
		m.handler.HandleMessage(msg)
		
		return waitForMessages(m.sub), true

	case trackerEndedMsg:
		m.trackerStatus = trackerDisconnected

		if m.reconnectTimeout == 0 {
			m.reconnectTimeout = time.Second
		} else if m.reconnectTimeout < 10 * time.Second {
			m.reconnectTimeout *= 2
		}

		return reconnectAfterCmd(m.reconnectTimeout), true
	}

	return nil, false
}
