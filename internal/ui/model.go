package ui

import (
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/rangaroo/2gis-friends/internal/state"
)

type Model struct {
	store *state.GlobalStore
	table table.Model
}

func NewModel(store *state.GlobalStore) Model {

	return Model{
		store:     store,
		table:     NewTable(),
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return t
	})
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "q" || msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
	case time.Time:
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

		return m, tea.Tick(time.Second, func (t time.Time) tea.Msg {
			return t
		})
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
