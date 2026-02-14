package main

import (
	"github.com/charmbracelet/lipgloss"
)

var (
	productStyle       = lipgloss.NewStyle().Foreground(lipgloss.Color("#00ff88"))
	activeProductStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#00ff88"))
	headerStyle        = lipgloss.NewStyle().Background(lipgloss.Color("#fcba03")).Bold(true)
	XStyle             = lipgloss.NewStyle().Foreground(lipgloss.Color("#e60707"))
)
