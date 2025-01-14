package ui

import (
	"log"
	"sort"
	"time"

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

func (m model) toggleItemAction(it item) list.Item {
	log.Printf("toggleItem: %v\n", it)
	if it.rejected {
		if it.name == "" {
			it.action = "create"
			it.rejected = false
			return it
		}
		if action, exists := m.actions[it.name]; exists {
			it.action = action.Type
			it.result = action.Result
			it.rejected = false
		}
		return it
	}
	if it.action == "keep" {
		return it
	}
	if it.name != "" {
		it.result = it.name
		it.action = "keep"
	} else {
		it.action = "!create"
	}
	it.rejected = true
	return it
}

func (m *model) setItems(files fsutils.FileList) tea.Cmd {
	m.files = files
	items := filesToItems(m.files)
	return m.list.SetItems(items)
}

func (m *model) updateResults(actions []llm.Action) tea.Cmd {
	m.maxDepth = maxDepth(actions)
	m.actions = actionsToMap(actions)
	return m.list.SetItems(m.resultsToItems(m.actions))
}

func (m model) resultsToItems(actions map[string]llm.Action) []list.Item {
	items := filesToItems(m.files)
	remaining := make(map[string]llm.Action)
	for k, v := range actions {
		remaining[k] = v
	}

	for i, listItem := range items {
		fileItem := listItem.(item)
		if action, exists := remaining[fileItem.name]; exists {
			fileItem.action = action.Type
			fileItem.result = action.Result
			items[i] = fileItem
			delete(remaining, fileItem.name)
		} else {
			fileItem.action = "keep"
			fileItem.result = fileItem.name
		}
	}

	for _, remainingAction := range remaining {
		if remainingAction.Type != "create" {
			continue
		}
		log.Printf("add create actions: %s", remainingAction)
		newItem := item{
			action: remainingAction.Type,
			result: remainingAction.Result,
		}
		items = append(items, newItem)
	}

	return items
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
	return tea.Tick(time.Millisecond*200, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}
