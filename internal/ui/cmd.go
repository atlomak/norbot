package ui

import (
	"strings"

	"github.com/atlomak/norbot/internal/fsutils"
	"github.com/atlomak/norbot/internal/llm"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type readDirMsg struct {
	files fsutils.FileList
	err   error
}

type queryResultMsg struct {
	actions []llm.Action
	err     error
}

func readDir(root string) tea.Cmd {
	return func() tea.Msg {
		files, err := fsutils.ReadDir(root, 1)
		if err != nil {
			return readDirMsg{
				err: err,
			}
		}
		return readDirMsg{files: files, err: nil}
	}
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

func (m model) queryResult(files fsutils.FileList) tea.Cmd {
	return func() tea.Msg {
		actions, err := m.llm.Query(files, "")
		if err != nil {
			return queryResultMsg{err: err}
		}
		return queryResultMsg{actions: actions, err: nil}
	}
}
