package ui

import (
	"time"
	"context"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/rangaroo/2gis-friends/internal/state"
    "github.com/rangaroo/2gis-friends/internal/config"
    "github.com/rangaroo/2gis-friends/internal/handler"
)

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

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// try to handle msg with custom commands
	if cmd, handled := m.updateTable(msg); handled {
		return m, cmd
	}

	if cmd, handled := m.updateTracker(msg); handled {
		return m, cmd
	}

	// default key handling
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "q" || msg.String() == "ctrl+c" {
			// close the model's context
			if m.cancel != nil {
				m.cancel()
			}
			return m, tea.Quit
		}
	}
	var cmd tea.Cmd
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
