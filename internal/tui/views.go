package tui

import "fmt"

func (m Model) welcomeView() string {
	s := "Hello, I'm Norbot, your nifty odd-jobbing robot. Please ensure all important data is backed up.\n\n"
	s += "Press Enter to unleash gnome...\n"
	return s
}

func (m Model) waitingView() string {
	entries := make([]Entry, 0, 5)
	maxLeft := 0
	for _, file := range m.files {
		entries = append(entries, Entry{left: file.Name, arrow: m.spinner.View()})
		if len(file.Name) > maxLeft {
			maxLeft = len(file.Name)
		}
	}
	s := "Pointlessly blowing files around!\n\n"
	for _, e := range entries {
		s += fmt.Sprintf("%-*s %s %s\n", maxLeft, e.left, e.arrow, e.right)
	}
	return s
}

func (m Model) changesView() string {
	entries := make([]Entry, 0, 5)
	maxLeft := 0
	for _, file := range m.files {
		if m.result[file.Name].Type == "move" {
			entries = append(entries, Entry{file.Name, ">", m.result[file.Name].Result})
		}
		if m.result[file.Name].Type == "rename" {
			entries = append(entries, Entry{file.Name, "~", m.result[file.Name].Result})
		}
		if m.result[file.Name].Type == "keep" {
			entries = append(entries, Entry{file.Name, "=", file.Name})
		}
		if len(file.Name) > maxLeft {
			maxLeft = len(file.Name)
		}
	}

	s := "Your files need mowing! More mowing!\n\n"
	for _, e := range entries {
		s += fmt.Sprintf("%-*s %s %s\n", maxLeft, e.left, e.arrow, e.right)
	}
	s += "\nPress y to apply changes..."
	return s
}

func (m Model) filesView() string {
	entries := make([]Entry, 0, 5)
	maxLeft := 0
	for _, file := range m.files {
		entries = append(entries, Entry{left: file.Name})
		if len(file.Name) > maxLeft {
			maxLeft = len(file.Name)
		}
	}
	s := "Current state of your garden mr.\n\n"
	for _, e := range entries {
		s += fmt.Sprintf("%-*s %s %s\n", maxLeft, e.left, e.arrow, e.right)
	}
	s += "\nPress Enter to rethink it a bit..."
	return s
}

func (m Model) View() string {

	s := ""
	if !m.started {
		s += m.welcomeView()
		return s
	}

	if m.err != nil {
		s += m.err.Error()
	}

	if m.waiting {
		s += m.waitingView()
	} else if m.result != nil {
		s += m.changesView()
	} else {
		s += m.filesView()
	}

	s += "\nPress space to refresh directory.\n"
	s += "Press q to quit.\n"

	return s
}
