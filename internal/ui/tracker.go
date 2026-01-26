package ui

import (
	"context"
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/rangaroo/2gis-friends/internal/core"
)

// bubbletea message and command
type trackerStatus int

// TODO: display status in footer
const (
	trackerDisconnected trackerStatus = iota
	trackerConnecting
	trackerConnected
)

type trackerEndedMsg struct {
	Err error
}

type trackerReconnectMsg struct{}

func startTrackerCmd(
	ctx context.Context,
	cfg *core.Config,
	h *core.Handler,
	state *core.GlobalState,
) tea.Cmd {
	return func() tea.Msg {
		ws, err := core.Connect(cfg)
		if err != nil {
			return trackerEndedMsg{Err: err}
		}
		defer ws.Close()

		done := make(chan error, 1)

		go func() {
			err := ws.ReadMessages(h.HandleMessage)
			done <- err
		}()

		select {
		case <-ctx.Done():
			return trackerEndedMsg{Err: nil}
		case err := <-done:
			return trackerEndedMsg{Err: err}
		}
	}
}

func reconnectAfterCmd(d time.Duration) tea.Cmd {
	return func() tea.Msg {
		time.Sleep(d)
		return trackerReconnectMsg{}
	}
}

func (m *Model) updateTracker(msg tea.Msg) (tea.Cmd, bool) {
	switch msg.(type) {
	case trackerEndedMsg:
		m.trackerStatus = trackerDisconnected

		if m.ctx.Err() != nil { // user ended the application
			return nil, true
		}

		if m.backoff < 10*time.Second {
			m.backoff *= 2
		}

		return reconnectAfterCmd(m.backoff), true

	case trackerReconnectMsg:
		if m.ctx.Err() != nil { // user ended the application
			return nil, true
		}

		m.trackerStatus = trackerConnecting
		if m.backoff == 0 {
			m.backoff = time.Second
		}

		return startTrackerCmd(m.ctx, m.cfg, m.handler, m.state), true
	}

	return nil, false
}
