package ui

import (
	"github.com/charmbracelet/lipgloss"
)

var (
	focusedModelStyle = lipgloss.NewStyle().
				Width(120).
				Height(8).
				Align(lipgloss.Left, lipgloss.Top).
				PaddingLeft(2)
	// BorderStyle(lipgloss.NormalBorder()).
	// BorderForeground(lipgloss.Color("69"))
	// focused = lipgloss.NewStyle().
	// 	Width(80).
	// 	Align(lipgloss.Left, lipgloss.Top)
	// 	// BorderStyle(lipgloss.NormalBorder()).
	// 	// BorderForeground(lipgloss.Color("69"))
	bottomStatusStyle = lipgloss.NewStyle().MarginLeft(4)
)

const (
	norbot = `
    _  __  ____    ___    ___   ____  ______
   / |/ / / __ \  / _ \  / _ ) / __ \/_  __/
  /    / / /_/ / / , _/ / _  |/ /_/ / / /   
 /_/|_/  \____/ /_/|_| /____/ \____/ /_/      
`
)

func (m model) welcomePanelView() string {
	s := norbot
	s += "\n"
	s += "\n"
	return s
}

func (m model) loadingPanelView() string {
	s := norbot
	s += "\n"
	s += bottomStatusStyle.Render(m.progress.View())
	s += "\n"
	return s
}

func (m model) readyPanelView() string {
	s := norbot
	s += "\n"
	s += bottomStatusStyle.Render("Press y to apply Norbot changes")
	s += "\n"
	return s
}

func (m model) finishPanelView() string {
	s := norbot
	s += "\n"
	s += bottomStatusStyle.Render("Norbot finished. No job is too little!")
	s += "\n"
	return s
}

func (m model) errorPanelView() string {
	s := norbot
	s += "\n"
	s += bottomStatusStyle.Render("Norbot encountered an error! Geez...")
	s += "\n"
	return s
}
