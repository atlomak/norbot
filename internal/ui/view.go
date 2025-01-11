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

type item string

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

	str := renderItem(i)

	fn := itemStyle.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return selectedItemStyle.Render("> " + strings.Join(s, " "))
		}
	}

	fmt.Fprint(w, fn(str))
}

func renderItem(it item) string {
	var b strings.Builder
	s := string(it)
	isDir := strings.HasSuffix(s, "/")

	path := strings.Split(strings.TrimSuffix(s, "/"), "/")
	s = path[len(path)-1]

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
		b.WriteString(fmt.Sprintf("%s %s", dirIcon, s))
		return b.String()
	}
	b.WriteString(fmt.Sprintf("%s %s", fileIcon, s))
	return b.String()
}

func (m model) View() string {
	return "\n" + m.list.View()
}
