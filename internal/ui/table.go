package ui

import (
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// bubbletea message and command that run every second
type tickMsg time.Time

func tickCmd() tea.Cmd {
	return tea.Every(time.Second, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func NewTable() table.Model {
	// define column names
	columns := []table.Column{
		{Title: "Name", Width: 20},
		{Title: "Battery", Width: 10},
		{Title: "Speed", Width: 10},
		{Title: "Last updated", Width: 15},
		{Title: "Coordinates", Width: 20},
	}

	// row values are updated in update hook
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

	return t
}

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

func (m *Model) updateTable(msg tea.Msg) (tea.Cmd, bool) {
	switch msg.(type) {
	case tickMsg:
		friends := m.state.GetViewData()

		rows := []table.Row{}

		for _, f := range friends {
			ago := time.Since(f.LastSeen).Round(time.Second)
			row := table.Row{
				f.Name,
				fmt.Sprintf("%.0f%%", f.Battery),
				fmt.Sprintf("%.1f", f.Speed),
				fmt.Sprintf("%s", ago),
				fmt.Sprintf("%.4f, %.4f", f.Lat, f.Lon),
			}

			rows = append(rows, row)
		}
		m.table.SetRows(rows)

		// run next tick
		return tickCmd(), true
	}

	return nil, false
}
