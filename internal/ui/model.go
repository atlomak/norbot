package ui

import (
	"log"

	"github.com/atlomak/norbot/internal/fsutils"
	"github.com/atlomak/norbot/internal/llm"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	list        list.Model
	files       fsutils.FileList
	actions     map[string]llm.Action
	llm         *llm.GeminiModel
	maxDepth    int
	textInput   textinput.Model
	progress    progress.Model
	progessDone bool
	status      status
	err         error
}

type status int

const (
	Started status = iota
	Input
	Waiting
	Ready
	Finished
	Error
)

func (m model) Init() tea.Cmd {
	return readDir(".", 0)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetHeight(msg.Height - statusPanelStyle.GetHeight())
		m.list.SetWidth(msg.Width)
		return m, nil
	case readDirMsg:
		if msg.err != nil {
			m.handleError(msg.err, msg)
			return m, nil
		}
		return m, m.setItems(msg.files)
	case queryResultMsg:
		if msg.err != nil {
			m.handleError(msg.err, msg)
			return m, nil
		}
		m.progessDone = true
		updateResults := m.updateResults(msg.actions)
		return m, tea.Batch(updateResults, m.sortItems)
	case applyChangesMsg:
		if msg.err != nil {
			m.handleError(msg.err, msg)
			return m, nil
		}
		return m, readDir(".", m.maxDepth)
	case tickMsg:
		if m.progessDone && m.progress.Percent() < 1.0 {
			cmd := m.progress.SetPercent(1.0)
			return m, tea.Sequence(cmd, tickCmd())
		} else if m.progessDone {
			m.status = Ready
			return m, nil
		}

		cmd := m.progress.IncrPercent(0.03)
		return m, tea.Sequence(cmd, tickCmd())
	case progress.FrameMsg:
		progressModel, cmd := m.progress.Update(msg)
		m.progress = progressModel.(progress.Model)
		return m, cmd
	case tea.KeyMsg:
		if m.status == Input {
			switch keypress := msg.String(); keypress {
			case "q", "ctrl+c":
				return m, tea.Quit
			case "enter":
				if m.status == Finished {
					return m, tea.Quit
				}
				m.progessDone = false
				m.status = Waiting
				m.textInput.Blur()
				return m, m.startQuery(m.files, m.textInput.Value())
			}
			var promptCmd tea.Cmd
			m.textInput, promptCmd = m.textInput.Update(msg)
			return m, promptCmd
		}
		switch keypress := msg.String(); keypress {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "enter":
			if m.status == Finished {
				return m, tea.Quit
			}
			m.progessDone = false
			m.status = Waiting
			return m, m.startQuery(m.files, "")
		case "y":
			if m.status == Finished {
				return m, tea.Quit
			}
			m.status = Finished
			return m, m.applyChanges
		case " ":
			if m.status != Ready {
				return m, nil
			}
			return m, tea.Sequence(m.toggleItem, m.sortItems)
		case "p":
			m.status = Input
			m.textInput.Focus()
			return m, nil
		}
	}

	var listCmd, promptCmd tea.Cmd

	m.textInput, promptCmd = m.textInput.Update(msg)
	m.list, listCmd = m.list.Update(msg)
	return m, tea.Batch(promptCmd, listCmd)
}

func (m *model) handleError(err error, msg tea.Msg) {
	m.status = Error
	m.err = err
	log.Printf("msg: %T returned error: %s", msg, err.Error())
}

func (m model) View() string {
	var statusPanel string
	switch m.status {
	case Started:
		statusPanel = m.welcomePanelView()
	case Input:
		statusPanel = m.inputPanelView()
	case Waiting:
		statusPanel = m.loadingPanelView()
	case Ready:
		statusPanel = m.readyPanelView()
	case Finished:
		statusPanel = m.finishPanelView()
	case Error:
		statusPanel = m.errorPanelView()
		return lipgloss.JoinVertical(lipgloss.Top, statusPanelStyle.Render(statusPanel), m.err.Error())
	}
	s := lipgloss.JoinVertical(lipgloss.Top, statusPanelStyle.Render(statusPanel), m.list.View())
	return s
}

func InitModel(llm *llm.GeminiModel) model {

	progess := progress.New(progress.WithDefaultScaledGradient())
	l := initList()
	textInput := textinput.New()
	m := model{list: l, llm: llm, progress: progess, status: Started, textInput: textInput}

	return m
}
