package ui

import (
	"log"
	"strings"

	"github.com/atlomak/norbot/internal/fsutils"
	"github.com/atlomak/norbot/internal/llm"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const listHeight = 14

var (
	titleStyle        = lipgloss.NewStyle().MarginLeft(2)
	itemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170"))
	paginationStyle   = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	helpStyle         = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)
)

type model struct {
	list    list.Model
	files   fsutils.FileList
	llm     *llm.GeminiModel
	spinner spinner.Model
	waiting bool
}

func (m model) Init() tea.Cmd {
	return tea.Batch(readDir("internal/test_dir"), m.spinner.Tick)
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
		for _, v := range msg.actions {
			log.Println(v)
		}
		m.waiting = false
		return m, nil
	}

	var listCmd, spinCmd tea.Cmd

	m.spinner, spinCmd = m.spinner.Update(msg)
	m.list, listCmd = m.list.Update(msg)
	return m, tea.Batch(spinCmd, listCmd)
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

func InitModel(llm *llm.GeminiModel) model {
	items := []list.Item{}

	const defaultWidth = 30

	l := list.New(items, itemDelegate{}, defaultWidth, listHeight)
	l.Title = "Files"
	l.SetShowStatusBar(true)
	l.SetFilteringEnabled(false)
	l.Styles.Title = titleStyle
	l.Styles.PaginationStyle = paginationStyle
	l.Styles.HelpStyle = helpStyle

	spin := spinner.New()
	spin.Spinner = spinner.Pulse

	m := model{list: l, llm: llm, spinner: spin}

	return m
}
