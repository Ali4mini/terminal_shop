package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	coffeeModels []string
	cursor       int
	selected     map[int]struct{}
	err          error
}

func initialModel() model {
	return model{
		coffeeModels: []string{"arch", "debian", "manjaro"},
		selected:     make(map[int]struct{}),
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) View() string {
	var lines []string

	lines = append(lines, headerStyle.Render("hello world"))
	lines = append(lines, "")

	for i, coffee := range m.coffeeModels {
		cursor := " "
		if i == m.cursor {
			cursor = ">"
		}

		if _, ok := m.selected[i]; ok {
			row := XStyle.Render(fmt.Sprintf("%s [X] %s", cursor, coffee))

			lines = append(lines, row)
		} else {

			row := productStyle.Render(fmt.Sprintf("%s [ ] %s", cursor, coffee))
			lines = append(lines, row)
		}

	}
	return lipgloss.JoinVertical(lipgloss.Left, lines...)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "k", "up":
			if m.cursor > 0 {
				m.cursor--
			}
		case "j", "down":
			if m.cursor < len(m.coffeeModels)-1 {
				m.cursor++
			}
		case "enter", " ":
			_, ok := m.selected[m.cursor]
			if ok {
				delete(m.selected, m.cursor)

			} else {
				m.selected[m.cursor] = struct{}{}
			}

		}
	}

	return m, nil
}

func main() {
	_, err := tea.NewProgram(initialModel(), tea.WithAltScreen()).Run()

	if err != nil {
		fmt.Printf("error while running the TUI")

		os.Exit(1)
	}
}
