package main

import (
	"github.com/charmbracelet/lipgloss"
)

var (
	// Colors
	subtle      = lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#383838"}
	accentColor = lipgloss.Color("#D4A373")

	// Define a style that uses the accent color
	accentTextStyle = lipgloss.NewStyle().Foreground(accentColor)

	// Add a bold version for headers
	boldAccentStyle = lipgloss.NewStyle().Foreground(accentColor).Bold(true)
	highlight       = lipgloss.Color("#FAEDCD") // Cream
	special         = lipgloss.Color("#73AD21") // Green for "Selected"

	// Styles

	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(highlight).
			Background(accentColor).
			Padding(0, 1).
			MarginBottom(1)

	checkStyle = lipgloss.NewStyle().
			Foreground(special)

	helpStyle = lipgloss.NewStyle().
			Foreground(subtle).
			MarginTop(1)

	// Sidebar/Detail Art Style
	detailStyle = lipgloss.NewStyle().
			PaddingLeft(4).
			BorderStyle(lipgloss.NormalBorder()).
			BorderLeft(true).
			BorderForeground(subtle)
	dimStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#777777"))

	priceTagStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#000000")).
			Background(lipgloss.Color("#D4A373")).
			Padding(0, 1).
			Bold(true)

	// Make sure these are also defined
	headerStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("#FAEDCD")).Bold(true).Underline(true)
	itemStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("#D9DCCF"))
	selectedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#D4A373")).Bold(true)
	mainStyle     = lipgloss.NewStyle().BorderStyle(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("#383838")).Padding(1)
)
