package ui

import (
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/rangaroo/2gis-friends/internal/state"
)

type Model struct {
	store *state.GlobalStore
	table table.Model
}

func NewModel(store *state.GlobalStore) Model {
	columns := []table.Column{
		{Title: "Name", Width: 20},
		{Title: "Battery", Width: 10},
		{Title: "Last updated", Width: 15},
		{Title: "Coordinates", Width: 20},
	}

	rows := []table.Row{table.Row{"Waiting..."}}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(15),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)
	t.SetStyles(s)

	return Model{
		store:     store,
		table:     t,
		Tabs       []string,
		TabContent []string,
		activeTab  int,
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
	case "right", "l", "n", "tab":
		m.activeTab = min(m.activeTab+1, len(m.Tabs)-1)
		return m, nil
	case "left", "h", "p", "shift+tab":
		m.activeTab = max(m.activeTab-1, 0)
		return m, nil
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

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

func (m Model) View() string {
	s := "\n  📡 2GIS FRIEND TRACKER\n"
	s += "  ──────────────────────\n"


	s += baseStyle.Render(m.table.View()) + "\n"

	s += "\n (Press 'q' to quit)\n"
	return s
}
