package ui

import (
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/rangaroo/2gis-friends/internal/state"
)

type Model struct {
	store *state.GlobalStore
}

func NewModel(store *state.GlobalStore) Model {
	return Model{store: store}
}

func (m Model) Init() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return t
	})
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "q" || msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
	case time.Time:
		return m, tea.Tick(time.Second, func (t time.Time) tea.Msg {
			return t
		})
	}
	return m, nil
}

func (m Model) View() string {
	s := "\n  📡 2GIS FRIEND TRACKER\n"
	s += "  ──────────────────────\n\n"

	// friends list
	friends := m.store.GetViewData()

	if len(friends) == 0 {
		s += " Waiting for data...\n"
	}

	for _, f := range friends {
		ago := time.Since(f.LastSeen).Round(time.Second)

		s += fmt.Sprintf(" %s | %3.0f%% | %6s ago | %.4f, %.4f\n",
		f.Name, f.Battery, ago, f.Lat, f.Lon)
	}

	s += "\n (Press 'q' to quit)\n"
	return s
}
