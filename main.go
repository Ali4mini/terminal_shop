package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/ssh"
	"github.com/charmbracelet/wish"
	"github.com/charmbracelet/wish/bubbletea"
)

type productsMsg []string

type model struct {
	coffeeModels []string
	cursor       int
	selected     map[int]struct{}
	width        int
	height       int
	err          error
}

func fetchProducts() tea.Msg {
	resp, err := http.Get("http://localhost:9991/products")
	if err != nil {
		fmt.Printf("error from 69 chambers, %s", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("had an error reading from body")
	}

	return productsMsg(strings.Split(string(body), ";"))

}

func initialModel() model {
	return model{
		coffeeModels: []string{},
		selected:     make(map[int]struct{}),
	}
}

func (m model) Init() tea.Cmd {
	return fetchProducts
}

func (m model) View() string {
	if len(m.coffeeModels) == 0 {
		return "Loading or Server not started..."
	}
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
	s := lipgloss.JoinVertical(lipgloss.Left, lines...)
	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, s)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case productsMsg:
		m.coffeeModels = msg // Data received!
		return m, nil
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
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

func startSSHServer() {
	s, err := wish.NewServer(
		wish.WithAddress(net.JoinHostPort("localhost", "2222")),
		wish.WithMiddleware(
			bubbletea.Middleware(teaHandler),
		),
	)
	if err != nil {
		log.Fatalln("Could not start server", "error", err)
	}
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	log.Printf("Starting SSH server on localhost:2222")
	go func() {
		if err = s.ListenAndServe(); err != nil {
			log.Fatalln(err)
		}
	}()

	<-done
	log.Println("Stopping SSH server...")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := s.Shutdown(ctx); err != nil {
		log.Fatalln(err)
	}

}

// teaHandler is the bridge between Wish and Bubble Tea
func teaHandler(s ssh.Session) (tea.Model, []tea.ProgramOption) {
	// This function runs every time a new user connects via SSH
	m := initialModel()

	// We return a new Bubble Tea program for this specific session
	return m, []tea.ProgramOption{tea.WithAltScreen()}
}

func main() {

	go StartMockServer()
	startSSHServer()

}
