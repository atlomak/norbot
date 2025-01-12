package ui

import (
	"log"
	"sort"
	"strings"

	"github.com/atlomak/norbot/internal/fsutils"
	"github.com/atlomak/norbot/internal/llm"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	titleStyle        = lipgloss.NewStyle().MarginLeft(2)
	itemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170"))
	paginationStyle   = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	helpStyle         = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)
)

type model struct {
	list     list.Model
	files    fsutils.FileList
	actions  map[string]llm.Action
	llm      *llm.GeminiModel
	spinner  spinner.Model
	waiting  bool
	ready    bool
	maxDepth int
}

func (m model) Init() tea.Cmd {
	return tea.Batch(readDir(".", 0), m.spinner.Tick)
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
			m.waiting = true
			queryCmd := m.queryResult(m.files)
			return m, queryCmd
		case "y":
			m.actions = nil
			return m, m.applyChanges
		}
	case readDirMsg:
		if msg.err != nil {
			log.Fatal(msg.err.Error())
		}
		m.files = msg.files
		items := filesToItems(m.files)
		cmd := m.list.SetItems(items)
		return m, cmd
	case queryResultMsg:
		if msg.err != nil {
			log.Fatal(msg.err.Error())
		}
		m.maxDepth = maxDepth(msg.actions)
		m.actions = actionsToMap(msg.actions)
		cmd := m.list.SetItems(m.resultsToItems(m.actions))
		m.waiting = false
		m.ready = true
		return m, cmd
	case applyChangesMsg:
		if msg.err != nil {
			log.Fatal(msg.err.Error())
		}
		return m, readDir(".", m.maxDepth)
	}

	var listCmd, spinCmd tea.Cmd

	m.spinner, spinCmd = m.spinner.Update(msg)
	m.list, listCmd = m.list.Update(msg)
	return m, tea.Batch(spinCmd, listCmd)
}

func (m model) resultsToItems(actionResults map[string]llm.Action) []list.Item {
	items := filesToItems(m.files)
	for index, listItem := range items {
		fileItem := listItem.(item)
		if action, exists := actionResults[fileItem.name]; exists {
			fileItem.action = action.Type
			fileItem.result = action.Result
			items[index] = fileItem
			delete(actionResults, fileItem.name)
		}
	}

	for _, remainingAction := range actionResults {
		log.Printf("add create actions: %s", remainingAction)
		newItem := item{
			action: remainingAction.Type,
			result: remainingAction.Result,
		}
		items = append(items, newItem)
	}

	sort.Slice(items, func(i, j int) bool {
		itemA := items[i].(item)
		itemB := items[j].(item)
		return itemA.result <= itemB.result
	})

	return items
}

func filesToItems(files fsutils.FileList) []list.Item {
	items := make([]list.Item, 0, len(files))

	s := strings.Split(files.String(), "\n")
	s = s[:len(s)-1] // because of newline at the end of string

	for _, file := range s {
		items = append(items, item{name: file})
	}
	return items
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
	items := []list.Item{}

	const defaultWidth = 120
	const listHeight = 30

	l := list.New(items, itemDelegate{}, defaultWidth, listHeight)
	l.Title = "Files"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = titleStyle
	l.Styles.PaginationStyle = paginationStyle
	l.Styles.HelpStyle = helpStyle

	spin := spinner.New()
	spin.Spinner = spinner.Pulse

	m := model{list: l, llm: llm, spinner: spin}

	return m
}
