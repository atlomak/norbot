package ui

import (
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

const (
	dirIcon  = "\U0001F4C1"
	fileIcon = "\U0001F4C4"
)

type item struct {
	name   string
	action string
	result string
}

func (i item) FilterValue() string { return "" }

type itemDelegate struct{}

func (d itemDelegate) Height() int                             { return 1 }
func (d itemDelegate) Spacing() int                            { return 0 }
func (d itemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(item)
	if !ok {
		return
	}

	const colWidthName = 30
	const colWidthAction = 10
	const colWidthResult = 50
	str := fmt.Sprintf("%-*s %-*s %-*s", colWidthName, renderItem(i.name), colWidthAction, i.action, colWidthResult, renderItem(i.result))

	fn := itemStyle.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return selectedItemStyle.Render("> " + strings.Join(s, " "))
		}
	}

	fmt.Fprint(w, fn(str))
}

func renderItem(file string) string {
	if file == "" {
		return file
	}

	var b strings.Builder
	isDir := strings.HasSuffix(file, "/")

	path := strings.Split(strings.TrimSuffix(file, "/"), "/")
	file = path[len(path)-1]

	parents := len(path) - 1

	if parents == 1 {
		b.WriteString("└─")
	} else if parents > 1 {
		for i := 0; i < parents-1; i++ {
			b.WriteString("│ ")
		}
		b.WriteString("└─")
	}

	if isDir {
		b.WriteString(fmt.Sprintf("%s %s", dirIcon, file))
		return b.String()
	}
	b.WriteString(fmt.Sprintf("%s %s", fileIcon, file))
	return b.String()
}

func (m model) View() string {
	var s string
	if m.waiting {
		s += m.spinner.View()
	}
	s += "\n" + m.list.View()
	return s
}
