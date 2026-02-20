package main

import (
	"context"
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
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/ssh"
	"github.com/charmbracelet/wish"
	"github.com/charmbracelet/wish/bubbletea"
)

// --- ADAPTER ---

// We need a wrapper because Product has a field named 'Description',
// but the list interface needs a method named 'Description()'.
type item struct {
	Product
}

// Implement list.Item interface on the wrapper
func (i item) Title() string       { return i.Name }
func (i item) FilterValue() string { return i.Name }
func (i item) Description() string {
	// In the list view, we show Price and Origin
	return fmt.Sprintf("$%.2f | %s", i.Price, i.Origin)
}

// --- MESSAGES ---

type productsMsg []Product
type errMsg error

// --- MODEL ---

type model struct {
	list     list.Model
	spinner  spinner.Model
	selected map[string]struct{} // Keyed by Product ID
	detail   Product             // Currently selected product for detail view
	loading  bool
	width    int
	height   int
	err      error
}

func initialModel() model {
	// Initialize Spinner
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(accentColor)

	// Initialize List
	delegate := list.NewDefaultDelegate()
	delegate.Styles.SelectedTitle = delegate.Styles.SelectedTitle.Foreground(accentColor).BorderLeftForeground(accentColor)
	delegate.Styles.SelectedDesc = delegate.Styles.SelectedDesc.Foreground(highlight)

	l := list.New([]list.Item{}, delegate, 0, 0)
	l.Title = "Coffee Menu"
	l.Styles.Title = titleStyle
	l.SetShowHelp(false)

	return model{
		list:     l,
		spinner:  s,
		loading:  true,
		selected: make(map[string]struct{}),
	}
}

// --- COMMANDS ---

func fetchProducts() tea.Msg {
	// Artificial delay to show the spinner (optional)
	time.Sleep(500 * time.Millisecond)

	c := &http.Client{Timeout: 5 * time.Second}
	// Note: Using the port 9991 defined in mockserver.go
	resp, err := c.Get("http://localhost:9991/products")
	if err != nil {
		return errMsg(err)
	}
	defer resp.Body.Close()

	var products []Product
	if err := json.NewDecoder(resp.Body).Decode(&products); err != nil {
		return errMsg(err)
	}

	return productsMsg(products)
}

// --- HELPER: PERSIAN TEXT ---
func fixPersian(input string) string {
	return garabic.Shape(input)
}

// --- UPDATE ---

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		// Responsive layout: Sidebar is 35% of width
		sidebarWidth := int(float64(m.width) * 0.35)
		if sidebarWidth < 35 {
			sidebarWidth = 35
		}
		m.list.SetSize(sidebarWidth, m.height-4)
		return m, nil

	case productsMsg:
		m.loading = false
		// Convert Products to list.Items via the 'item' wrapper
		items := make([]list.Item, len(msg))
		for i, p := range msg {
			items[i] = item{Product: p}
		}
		cmd = m.list.SetItems(items)
		cmds = append(cmds, cmd)

		// Set initial selection logic
		if len(msg) > 0 {
			m.detail = msg[0]
		}

	case errMsg:
		m.err = msg
		m.loading = false
		return m, nil

	case tea.KeyMsg:
		// If filtering, let list handle keys exclusively
		if m.list.FilterState() == list.Filtering {
			m.list, cmd = m.list.Update(msg)
			cmds = append(cmds, cmd)
			return m, tea.Batch(cmds...)
		}

		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "enter", " ":
			if i, ok := m.list.SelectedItem().(item); ok {
				if _, exists := m.selected[i.ID]; exists {
					delete(m.selected, i.ID)
				} else {
					m.selected[i.ID] = struct{}{}
				}
			}
		}
	}

	// Update Spinner
	if m.loading {
		m.spinner, cmd = m.spinner.Update(msg)
		cmds = append(cmds, cmd)
	}

	// Update List and catch selection changes
	prevItem := m.list.SelectedItem()
	m.list, cmd = m.list.Update(msg)
	cmds = append(cmds, cmd)

	newItem := m.list.SelectedItem()
	// Update the Detail view if the selection changed
	if newItem != nil && newItem != prevItem {
		m.detail = newItem.(item).Product
	}

	return m, tea.Batch(cmds...)
}

// --- VIEW ---

func (m model) View() string {
	if m.err != nil {
		return fmt.Sprintf("\n  Error: %v\n\n  Press q to quit", m.err)
	}

	if m.loading {
		return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center,
			fmt.Sprintf("%s Brewing connection...", m.spinner.View()),
		)
	}

	// 1. Render Sidebar
	sidebar := listStyle.Render(m.list.View())

	// 2. Render Details
	var detailView string

	// Check if the current detail item is selected/checked
	status := " "
	if _, ok := m.selected[m.detail.ID]; ok {
		status = "✔ ORDER ADDED"
	}

	// Shape Persian text (Description field from mockserver)
	// We wrap it to prevent overflow in the detail view
	descWidth := m.width - m.list.Width() - 10
	wrappedDesc := lipgloss.NewStyle().Width(descWidth).Render(fixPersian(m.detail.Description))

	detailView = lipgloss.JoinVertical(lipgloss.Left,
		lipgloss.JoinHorizontal(lipgloss.Center,
			boldAccentStyle.Copy().Render(m.detail.Name),
			lipgloss.NewStyle().MarginLeft(2).Foreground(special).Render(status),
		),
		dimStyle.Render(fmt.Sprintf("%s | %s Roast", m.detail.Origin, m.detail.Roast)),
		strings.Repeat("-", descWidth),
		"",
		wrappedDesc,
		"",
		priceTagStyle.Render(fmt.Sprintf("PRICE: $%.2f", m.detail.Price)),
	)

	detailBox := detailStyle.
		Width(m.width - m.list.Width() - 6).
		Height(m.height - 2).
		Render(detailView)

	return lipgloss.JoinHorizontal(lipgloss.Top, sidebar, detailBox)
}

func (m model) Init() tea.Cmd {
	return tea.Batch(m.spinner.Tick, fetchProducts)
}

// --- SERVER SETUP ---

func startSSHServer() {
	s, err := wish.NewServer(
		wish.WithAddress(net.JoinHostPort("0.0.0.0", "2222")),
		wish.WithMiddleware(
			bubbletea.Middleware(teaHandler),
		),
	)
	if err != nil {
		log.Fatalln("Could not start server", "error", err)
	}
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	log.Printf("Starting SSH server on :2222")
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
	return initialModel(), []tea.ProgramOption{tea.WithAltScreen()}
}

func main() {
	go StartMockServer()

	// Decide between Local TUI or SSH Server
	if os.Getenv("TUI_NO_SSH_SERVER") == "true" {
		p := tea.NewProgram(initialModel(), tea.WithAltScreen())
		if _, err := p.Run(); err != nil {
			fmt.Printf("Alas, there's been an error: %v", err)
			os.Exit(1)
		}
	} else {
		startSSHServer()
	}
}
