package ui

import (
	"sort"
	"time"

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

type tickMsg time.Time

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

func (m model) queryResult(files fsutils.FileList, prompt string) tea.Cmd {
	return func() tea.Msg {
		actions, err := m.llm.Query(files, prompt)
		if err != nil {
			return queryResultMsg{err: err}
		}
		return queryResultMsg{actions: actions, err: nil}
	}
}

func (m *model) startQuery(files fsutils.FileList, prompt string) tea.Cmd {
	progressMsg := m.progress.SetPercent(0)
	tickCmd := tickCmd()
	queryCmd := m.queryResult(files, prompt)
	return tea.Sequence(progressMsg, tickCmd, queryCmd)
}

func (m *model) toggleItem() tea.Msg {
	selected := m.list.SelectedItem().(item)
	toggled := m.toggleItemAction(selected)
	idx := m.list.Index()
	return m.list.SetItem(idx, toggled)
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

func (m model) sortItems() tea.Msg {
	items := m.list.Items()
	sort.Slice(items, func(i, j int) bool {
		itemA := items[i].(item)
		itemB := items[j].(item)
		return itemA.result <= itemB.result
	})
	return m.list.SetItems(items)
}

func tickCmd() tea.Cmd {
	return tea.Tick(time.Millisecond*400, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}
