package tui

import (
	"context"

	"github.com/atlomak/norbot/internal/fsutils"
	"github.com/atlomak/norbot/internal/llm"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/google/generative-ai-go/genai"
)

type Model struct {
	waiting bool
	started bool
	files   []fsutils.FileInfo
	result  map[string]llm.Action
	err     error
	gemini  llm.GeminiModel
	spinner spinner.Model
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(m.spinner.Tick, fsutils.ListFiles("."))
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:

		switch msg.String() {

		case "ctrl+c", "q":
			return m, tea.Quit

		case " ":
			m.result = nil
			return m, fsutils.ListFiles(".")
		case "enter":
			m.waiting = true
			m.started = true
			return m, m.gemini.Query(m.files)
		}
	case fsutils.FilesMsg:
		if msg.Err != nil {
			m.err = msg.Err
			return m, nil
		}
		m.files = msg.Files
		return m, nil
	case llm.GeminiMsg:
		if msg.Err != nil {
			m.err = msg.Err
			return m, nil
		}
		m.waiting = false
		m.result = msg.Actions
		return m, nil
	default:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}
	return m, nil
}

type Entry struct {
	left  string
	arrow string
	right string
}

func InitModel(client *genai.Client, ctx context.Context) Model {
	s := spinner.New()
	s.Spinner = spinner.Points
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	genModel := llm.InitGeminiModel(client, ctx)
	return Model{spinner: s, gemini: genModel, waiting: false}
}
