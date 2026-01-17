package ui

import (
	"fmt"
	"time"
	"context"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/rangaroo/2gis-friends/internal/state"
	"github.com/rangaroo/2gis-friends/internal/client"
    "github.com/rangaroo/2gis-friends/internal/config"
    "github.com/rangaroo/2gis-friends/internal/handler"
)

// for tracker commands
type trackerStatus int

const (
	trackerDisconnected trackerStatus = iota
	trackerConnecting
	trackerConnected
)

type trackerEndedMsg struct {
	Err error
}

type trackerReconnectMsg struct{}

// for table commands
type tickMsg time.Time


type Model struct {
	table table.Model

	store         *state.GlobalStore

	cfg           *config.Config	
	handler       *handler.Handler

	ctx           context.Context
	cancel        context.CancelFunc

	trackerStatus trackerStatus
	backoff       time.Duration
}

func NewModel(
	store *state.GlobalStore,
	cfg   *config.Config,
	h     *handler.Handler,
) Model {
	// create context that cancels when Ctrl+C is pressed
	ctx, cancel := context.WithCancel(context.Background())

	return Model{
		table:         NewTable(),
		store:         store,
		cfg:           cfg,
		handler:       h,
		ctx:           ctx,
		cancel:        cancel,
		trackerStatus: trackerDisconnected,
		backoff:       time.Second,
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		startTrackerCmd(m.ctx, m.cfg, m.handler, m.store),
		tickCmd(),
	)
}

func tickCmd() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func startTrackerCmd(
	ctx   context.Context,
	cfg   *config.Config,
	h     *handler.Handler,
	store *state.GlobalStore,
) tea.Cmd {
	return func() tea.Msg {
		ws, err := client.Connect(cfg)
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
		case <- ctx.Done():
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

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "q" || msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
	case trackerEndedMsg:
		m.trackerStatus = trackerDisconnected

		if m.ctx.Err() != nil {
			return m, nil
		}

		if m.backoff < 10 * time.Second {
			m.backoff *= 2
		}

		return m, reconnectAfterCmd(m.backoff)

	case trackerReconnectMsg:
		if m.ctx.Err() != nil {
			return m, nil
		}

		m.trackerStatus = trackerConnecting
		if m.backoff == 0 {
			m.backoff = time.Second
		}

		return m, startTrackerCmd(m.ctx, m.cfg, m.handler, m.store)

	case tickMsg:
		friends := m.store.GetViewData()
		rows := []table.Row{}
		for _, f := range friends {
			ago := time.Since(f.LastSeen).Round(time.Second)

			rows = append(rows, table.Row{
				f.Name,
				fmt.Sprintf("%.0f%%", f.Battery),
				fmt.Sprintf("%s", ago),
				fmt.Sprintf("%.4f, %.4f", f.Lat, f.Lon),
			})
		}
		m.table.SetRows(rows)

		return m, tickCmd()
	}

	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	s := "\n  📡 2GIS FRIEND TRACKER\n"
	s += "  ──────────────────────\n"

	s += baseStyle.Render(m.table.View()) + "\n"

	s += "\n (Press 'q' to quit)\n"
	return s
}
