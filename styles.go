package main

import (
	"github.com/charmbracelet/lipgloss"
)

var (
	// Colors
	subtle      = lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#383838"}
	accentColor = lipgloss.Color("#D4A373")
	highlight   = lipgloss.Color("#FAEDCD") // Cream
	special     = lipgloss.Color("#73AD21") // Green

	// Text Styles
	boldAccentStyle = lipgloss.NewStyle().Foreground(accentColor).Bold(true)
	dimStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("#777777"))

	// List Styles
	titleStyle = lipgloss.NewStyle().
			Foreground(highlight).
			Background(accentColor).
			Padding(0, 1)

	// Container Styles
	listStyle = lipgloss.NewStyle().
			Margin(1, 2)

	detailStyle = lipgloss.NewStyle().
			Padding(1, 2).
			BorderStyle(lipgloss.NormalBorder()).
			BorderLeft(true).
			BorderForeground(subtle)

	priceTagStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#000000")).
			Background(lipgloss.Color("#D4A373")).
			Padding(0, 1).
			Bold(true).
			MarginTop(1)
)
