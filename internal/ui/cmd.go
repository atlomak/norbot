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

type applyChangesMsg struct {
	err error
}

func readDir(root string, depth int) tea.Cmd {
	return func() tea.Msg {
		files, err := fsutils.ReadDir(root, depth)
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

func (m model) applyChanges() tea.Msg {
	var err error
	for _, v := range m.list.Items() {
		i := v.(item)
		switch i.action {
		case "create":
			err = fsutils.CreateDir(i.result)
		case "move":
			err = fsutils.MoveFile(i.name, i.result)
		case "keep":
		}
	}
	return applyChangesMsg{err: err}
}
