package ui

import (
	"github.com/charmbracelet/lipgloss"
)

var (
	statusPanelStyle = lipgloss.NewStyle().
				Width(120).
				Height(8).
				Align(lipgloss.Left, lipgloss.Top).
				PaddingLeft(2)
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
	s += bottomStatusStyle.Render("Press enter to unleash the gnomes...")
	s += "\n"
	return s
}

func (m model) inputPanelView() string {
	s := norbot
	s += "\n"
	s += bottomStatusStyle.Render(m.textInput.View())
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
	s += bottomStatusStyle.Render("Press y to apply Norbot changes. Press space to reject selected file.")
	s += "\n"
	return s
}

func (m model) finishPanelView() string {
	s := norbot
	s += "\n"
	s += bottomStatusStyle.Render("Norbot finished. Bowing. More bowing")
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
