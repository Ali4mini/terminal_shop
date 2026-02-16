package main

import (
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/abdullahdiaa/garabic"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/ssh"
	"github.com/charmbracelet/wish"
	"github.com/charmbracelet/wish/bubbletea"
)

type productsMsg []Product

type model struct {
	coffeeModels []Product
	cursor       int
	selected     map[int]struct{}
	width        int
	height       int
	err          error
}

func fetchProducts() tea.Msg {
	resp, err := http.Get("http://localhost:9991/products")
	if err != nil {
		return productsMsg{} // In a real app, return an error message
	}
	defer resp.Body.Close()

	var products []Product
	if err := json.NewDecoder(resp.Body).Decode(&products); err != nil {
		return productsMsg{}
	}

	return productsMsg(products)
}

func initialModel() model {
	return model{
		coffeeModels: []Product{},
		selected:     make(map[int]struct{}),
	}
}

func FixPersian(input string) string {
	text := garabic.Shape(input)

	return text

}

//go:embed persian_welcome.txt
var welcomeText string

func (m model) headerView() string {

	return titleStyle.Render(welcomeText)
}
func (m model) Init() tea.Cmd {
	return fetchProducts
}

func (m model) View() string {
	if len(m.coffeeModels) == 0 {
		return "Waking up the beans..."
	}

	// SIDEBAR: Product List
	var listBuilder strings.Builder
	listBuilder.WriteString(headerStyle.Render(welcomeText) + "\n\n")

	for i, item := range m.coffeeModels {
		cursor := "  "
		style := itemStyle // Define this in your styles

		if i == m.cursor {
			cursor = "» "
			style = selectedStyle // Define this in your styles
		}

		checked := " "
		if _, ok := m.selected[i]; ok {
			checked = "✔"
		}

		// Show Name and Price in the list
		line := fmt.Sprintf("%s[%s] %-18s $%2.2f", cursor, checked, item.Name, item.Price)
		listBuilder.WriteString(style.Render(line) + "\n")
	}

	// MAIN VIEW: Product Details
	current := m.coffeeModels[m.cursor]

	// Create a "Card" for the details
	detailView := lipgloss.JoinVertical(lipgloss.Left,
		boldAccentStyle.Render(current.Name),
		dimStyle.Render(current.Origin+" | "+current.Roast+" Roast"),
		"",
		current.Description,
		"",
		priceTagStyle.Render(fmt.Sprintf("PRICE: $%.2f", current.Price)),
	)

	// Combine Sidebar and Details
	mainContent := lipgloss.JoinHorizontal(lipgloss.Top,
		lipgloss.NewStyle().Width(35).Render(listBuilder.String()),
		lipgloss.NewStyle().Padding(1, 4).Render(detailView),
	)

	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center,
		mainStyle.Render(mainContent),
	)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case productsMsg:
		m.coffeeModels = msg
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
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

func teaHandler(s ssh.Session) (tea.Model, []tea.ProgramOption) {
	m := initialModel()

	// We return a new Bubble Tea program for this specific session
	return m, []tea.ProgramOption{tea.WithAltScreen()}
}

func main() {

	go StartMockServer()
	startSSHServer()

}
