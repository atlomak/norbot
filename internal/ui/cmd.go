package ui

import (
	"github.com/atlomak/norbot/internal/fsutils"
	"github.com/atlomak/norbot/internal/llm"
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

func (m model) queryResult(files fsutils.FileList) tea.Cmd {
	return func() tea.Msg {
		actions, err := m.llm.Query(files, "")
		if err != nil {
			return queryResultMsg{err: err}
		}
		return queryResultMsg{actions: actions, err: nil}
	}
}
