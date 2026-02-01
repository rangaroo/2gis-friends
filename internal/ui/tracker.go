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
	return func() tea.Msg {
		conn, err := core.ConnectToWebSocket(cfg)
		if err != nil {
			return trackerEndedMsg{err: err}
		}
		
		return trackerConnectedMsg{conn: conn}
	}
}

func produceMessagesCmd(conn *core.WebSocketConn, sub chan interface{}) tea.Cmd {
	return func() tea.Msg {
		defer conn.Close()

		for {
			msg, err := conn.ReadMessages()
			if err != nil {
				sub <- trackerEndedMsg{err: err} // research how to deal with this situation
				return nil
			}

			sub <- trackerDataMsg(msg)
		}
	}
}

func waitForMessages(sub chan interface{}) tea.Cmd {
	return func() tea.Msg {
		return <-sub
	}
}

func reconnectAfterCmd(timeout time.Duration) tea.Cmd {
	return tea.Tick(timeout,
		func(_, time.Time) tea.Msg {
			return trackerReconnectMsg{}
		}
	)
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
