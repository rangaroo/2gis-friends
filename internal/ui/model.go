package ui

import (
	"context"
	"time"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/rangaroo/2gis-friends/internal/core"
)

type Model struct {
	// ui
	table table.Model

	// global state
	state *core.GlobalState

	cfg     core.Config
	handler *core.Handler

	ctx    context.Context
	cancel context.CancelFunc

	// concurrent tracker
	trackerStatus trackerStatus
	reconnectTimeout time.Duration
	sub chan interface{}
}

func NewModel(cfg core.Config, db *core.DatabaseClient) Model {
	state := core.NewState()

	h := core.NewHandler(db, state)

	// create context that cancels when Ctrl+C is pressed
	ctx, cancel := context.WithCancel(context.Background())

	return Model{
		table:         NewTable(),
		state:         state,
		cfg:           cfg,
		handler:       h,
		ctx:           ctx,
		cancel:        cancel,
		trackerStatus: trackerDisconnected,
		reconnectTimeout:       time.Second,
		sub: make(chan interface{}),
	}
}

func (m Model) Init() tea.Cmd {
	//return connectToTrackerCmd(m.cfg)
	return tea.Batch(
		connectToTrackerCmd(m.cfg),
		tickCmd(),
	)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// try to handle msg with custom commands first
	if cmd, handled := m.updateTable(msg); handled {
		return m, cmd
	}

	if cmd, handled := m.updateTracker(msg); handled {
		return m, cmd
	}

	// key press handling
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
	s := baseStyle.Render(m.table.View()) + "\n"

	s += "\n (Press 'q' to quit)\n"
	return s
}
