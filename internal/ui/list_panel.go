package ui

import (
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	dirIcon        = "\U0001F4C1"
	fileIcon       = "\U0001F4C4"
	newFile        = "\U00002728"
	colWidthName   = 40
	colWidthAction = 10
)

var (
	titleStyle        = lipgloss.NewStyle().MarginLeft(2)
	itemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170"))
	rejectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("#FF6347"))
	paginationStyle   = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	helpStyle         = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)
)

type item struct {
	rejected bool
	name     string
	action   string
	result   string
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

	var str string
	name := i.name
	if name == "" {
		str = fmt.Sprintf("%-*s %-*s %s", colWidthName+15, newFile, colWidthAction, i.action, renderItem(i.result))
	} else {
		if len(name) > colWidthName {
			name = trimName(name)
		}
		str = fmt.Sprintf("%-*s %-*s %s", colWidthName+15, renderItem(name), colWidthAction, i.action, renderItem(i.result))
	}

	fn := itemStyle.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return selectedItemStyle.Render("> " + strings.Join(s, " "))
		}
	} else if i.rejected {
		fn = func(s ...string) string {
			return rejectedItemStyle.Render("x " + strings.Join(s, " "))
		}
	}

	fmt.Fprint(w, fn(str))
}

func trimName(name string) string {
	r := []rune(name)
	trunc := r[:colWidthName-4]
	return string(trunc) + "..."
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

func initList() list.Model {
	items := []list.Item{}

	const defaultWidth = 120
	const listHeight = 30

	l := list.New(items, itemDelegate{}, defaultWidth, listHeight)
	// l.Title = "Files"
	l.SetShowTitle(false)
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.SetShowHelp(true)
	l.Styles.Title = titleStyle
	l.Styles.PaginationStyle = paginationStyle
	l.Styles.HelpStyle = helpStyle

	return l
}
