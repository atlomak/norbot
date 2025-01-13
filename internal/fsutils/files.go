package fsutils

import (
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type Inode struct {
}

type FilesMsg struct {
	Files []FileInfo
	Err   error
}

type FileInfo struct {
	Name        string
	IsDir       bool
	Size        int64
	Permissions string
	ModTime     time.Time
}

func ListFiles(dir string) tea.Cmd {
	return func() tea.Msg {
		entries, err := os.ReadDir(dir)
		if err != nil {
			return FilesMsg{Err: err}
		}

		var files []FileInfo
		for _, entry := range entries {
			info, err := entry.Info()
			if err != nil {
				continue
			}

			files = append(files, FileInfo{
				Name:        entry.Name(),
				Size:        info.Size(),
				Permissions: info.Mode().String(),
				ModTime:     info.ModTime(),
				IsDir:       info.IsDir(),
			})
		}
		return FilesMsg{Files: files, Err: nil}
	}
}
