package ui

import (
	"log"
	"strings"

	"github.com/atlomak/norbot/internal/fsutils"
	"github.com/atlomak/norbot/internal/llm"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	list     list.Model
	files    fsutils.FileList
	actions  map[string]llm.Action
	llm      *llm.GeminiModel
	maxDepth int
	progress progress.Model
	status   status
	err      error
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
			m.status = Waiting
			tickCmd := tickCmd()
			queryCmd := m.queryResult(m.files)
			return m, tea.Batch(tickCmd, queryCmd)
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
			selected := m.list.SelectedItem().(item)
			toggled := m.toggleItemAction(selected)
			idx := m.list.Index()
			itemCmd := m.list.SetItem(idx, toggled)
			return m, tea.Sequence(itemCmd, m.sortItems)
		}
	case readDirMsg:
		if msg.err != nil {
			m.status = Error
			m.err = msg.err
			log.Println(msg.err.Error())
			return m, nil
		}
		m.files = msg.files
		items := m.filesToItems(m.files)
		cmd := m.list.SetItems(items)
		return m, cmd
	case queryResultMsg:
		if msg.err != nil {
			m.status = Error
			m.err = msg.err
			log.Println(msg.err.Error())
			return m, nil
		}
		m.maxDepth = maxDepth(msg.actions)
		m.actions = actionsToMap(msg.actions)
		cmd := m.list.SetItems(m.resultsToItems(m.actions))

		m.status = Ready
		return m, tea.Batch(cmd, m.sortItems)
	case applyChangesMsg:
		if msg.err != nil {
			m.status = Error
			m.err = msg.err
			log.Println(msg.err.Error())
			return m, nil
		}
		return m, readDir(".", m.maxDepth)
	case tickMsg:
		if m.progress.Percent() == 1.0 && m.status == Ready {
			cmd := m.progress.SetPercent(0)
			return m, cmd
		}

		cmd := m.progress.IncrPercent(0.08)
		return m, tea.Batch(tickCmd(), cmd)
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
		return lipgloss.JoinVertical(lipgloss.Top, focusedModelStyle.Render(statusPanel), m.err.Error())
	}
	s := lipgloss.JoinVertical(lipgloss.Top, focusedModelStyle.Render(statusPanel), m.list.View())
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

func (m model) filesToItems(files fsutils.FileList) []list.Item {
	items := make([]list.Item, 0, len(files))

	s := strings.Split(files.String(), "\n")
	s = s[:len(s)-1] // because of newline at the end of string

	for _, file := range s {
		items = append(items, item{name: file})
	}
	return items
}

func (m model) toggleItemAction(it item) list.Item {
	if it.rejected {
		if action, exists := m.actions[it.name]; exists {
			it.action = action.Type
			it.result = action.Result
			it.rejected = false
		}
		return it
	}
	it.action = "keep"
	it.result = it.name
	it.rejected = true
	return it
}

func actionsToMap(actions []llm.Action) map[string]llm.Action {
	result := make(map[string]llm.Action)
	for _, action := range actions {
		log.Println(action)
		if action.Name == "" {
			result[action.Result] = action
		} else {
			result[action.Name] = action
		}
	}
	return result
}

func maxDepth(actions []llm.Action) int {
	maxDepth := 0
	for _, action := range actions {
		path := strings.Split(strings.TrimSuffix(action.Result, "/"), "/")
		parents := len(path) - 1
		if parents > maxDepth {
			maxDepth = parents
		}
	}
	return maxDepth
}

func InitModel(llm *llm.GeminiModel) model {

	progess := progress.New(progress.WithDefaultScaledGradient())
	l := initList()
	m := model{list: l, llm: llm, progress: progess, status: Started}

	return m
}
