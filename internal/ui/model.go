package ui

import (
	"log"

	"github.com/atlomak/norbot/internal/fsutils"
	"github.com/atlomak/norbot/internal/llm"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	list        list.Model
	files       fsutils.FileList
	actions     map[string]llm.Action
	llm         *llm.GeminiModel
	maxDepth    int
	progress    progress.Model
	progessDone bool
	status      status
	err         error
}

type status int

const (
	Started status = iota
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
	case tea.KeyMsg:
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
		}
	case readDirMsg:
		if msg.err != nil {
			m.status = Error
			m.err = msg.err
			log.Println(msg.err.Error())
			return m, nil
		}
		return m, m.setItems(msg.files)
	case queryResultMsg:
		if msg.err != nil {
			m.status = Error
			m.err = msg.err
			log.Println(msg.err.Error())
			return m, nil
		}
		m.progessDone = true
		updateResults := m.updateResults(msg.actions)
		return m, tea.Batch(updateResults, m.sortItems)
	case applyChangesMsg:
		if msg.err != nil {
			m.status = Error
			m.err = msg.err
			log.Println(msg.err.Error())
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
	}

	var listCmd, spinCmd tea.Cmd

	m.list, listCmd = m.list.Update(msg)
	return m, tea.Batch(spinCmd, listCmd)
}

func (m model) View() string {
	var statusPanel string
	switch m.status {
	case Started:
		statusPanel = m.welcomePanelView()
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

func (m model) resultsToItems(actions map[string]llm.Action) []list.Item {
	items := m.filesToItems(m.files)
	remaining := make(map[string]llm.Action)
	for k, v := range actions {
		remaining[k] = v
	}

	for index, listItem := range items {
		fileItem := listItem.(item)
		if action, exists := remaining[fileItem.name]; exists {
			fileItem.action = action.Type
			fileItem.result = action.Result
			items[index] = fileItem
			delete(remaining, fileItem.name)
		}
	}

	for _, remainingAction := range remaining {
		log.Printf("add create actions: %s", remainingAction)
		newItem := item{
			action: remainingAction.Type,
			result: remainingAction.Result,
		}
		items = append(items, newItem)
	}

	return items
}

func InitModel(llm *llm.GeminiModel) model {

	progess := progress.New(progress.WithDefaultScaledGradient())
	l := initList()
	m := model{list: l, llm: llm, progress: progess, status: Started}

	return m
}
