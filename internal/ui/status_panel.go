package ui

import (
	"github.com/charmbracelet/lipgloss"
)

var (
	statusPanelStyle = lipgloss.NewStyle().
				Width(120).
				Height(10).
				Align(lipgloss.Left, lipgloss.Top).
				PaddingLeft(2)
	statusTitleStyle  = lipgloss.NewStyle().MarginLeft(1).Foreground(lipgloss.Color(gnomeGreen))
	bottomStatusStyle = lipgloss.NewStyle().Margin(2)
	promptInputStyle  = lipgloss.NewStyle().
				Width(80).
				MarginLeft(2).
				BorderStyle(lipgloss.NormalBorder()).
				BorderForeground(lipgloss.Color(gnomeGreen))
)

const (
	norbot = `
    _  __  ____    ___    ___   ____  ______
   / |/ / / __ \  / _ \  / _ ) / __ \/_  __/
  /    / / /_/ / / , _/ / _  |/ /_/ / / /   
 /_/|_/  \____/ /_/|_| /____/ \____/ /_/      
`
	gnomeGreen = "#39FF14"
	darkGreen  = "#243407"
)

func (m model) welcomePanelView() string {
	s := statusTitleStyle.Render(norbot)
	s += bottomStatusStyle.Render("Press enter to unleash the gnomes...")
	return s
}

func (m model) inputPanelView() string {
	s := statusTitleStyle.Render(norbot)
	s += "\n"
	s += promptInputStyle.Render(m.textInput.View())
	s += "\n"
	return s
}

func (m model) loadingPanelView() string {
	s := statusTitleStyle.Render(norbot)
	s += bottomStatusStyle.Render(m.progress.View())
	return s
}

func (m model) readyPanelView() string {
	s := statusTitleStyle.Render(norbot)
	s += bottomStatusStyle.Render("Press y to apply Norbot changes. Press space to reject selected file.")
	return s
}

func (m model) finishPanelView() string {
	s := statusTitleStyle.Render(norbot)
	s += bottomStatusStyle.Render("Norbot finished. Bowing. More bowing")
	return s
}

func (m model) errorPanelView() string {
	s := norbot
	s += bottomStatusStyle.Render("Norbot encountered an error! Geez...")
	return s
}
