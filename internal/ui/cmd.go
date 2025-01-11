package ui

import (
	"strings"

	"github.com/atlomak/norbot/internal/fsutils"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type readDirMsg struct {
	items []list.Item
	err   error
}

func readDirCmd(root string) tea.Cmd {
	return func() tea.Msg {
		files, err := fsutils.ReadDir(root, 0)
		if err != nil {
			return readDirMsg{
				err: err,
			}
		}
		items := filesToItems(files)
		return readDirMsg{items: items, err: nil}
	}
}

func filesToItems(files fsutils.FileList) []list.Item {
	items := make([]list.Item, 0, len(files))
	s := strings.Split(files.String(), "\n")
	for _, file := range s {
		items = append(items, item(file))
	}
	return items
}
