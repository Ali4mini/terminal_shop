package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	coffeeModels []string
	cursor       int
	selected     int
	err          error
}

func initialModel() model {
	return model{
		coffeeModels: []string{"arch", "debian", "manjaro"},
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) View() string {
	s := "hello world\n"
	for _, coffee := range m.coffeeModels {
		s += fmt.Sprintf("%s\n", coffee)
	}
	return s
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q":
			return m, tea.Quit
		}
	}

	return m, nil
}

func main() {
	fmt.Printf("hello world")
	_, err := tea.NewProgram(initialModel()).Run()

	if err != nil {
		fmt.Printf("error while running the TUI")

		os.Exit(1)
	}
}
