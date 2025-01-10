package fsutils

import (
	"fmt"
	"io/fs"
	"os"
)

type DirMsg struct {
	Files []Node
	Err   error
}

type Node struct {
	Info     fs.FileInfo
	Children []Node
}

func ReadDir(root string, depth int) ([]Node, error) {
	entries, err := os.ReadDir(root)
	if err != nil {
		return nil, err
	}

	var files []Node
	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			continue
		}

		var node Node
		if depth != 0 && info.IsDir() {
			children, err := ReadDir(fmt.Sprintf("%s/%s", root, info.Name()), depth-1)
			if err != nil {
				return nil, err
			}
			node = Node{
				Info:     info,
				Children: children,
			}
		} else {
			node = Node{
				Info:     info,
				Children: nil,
			}
		}
		files = append(files, node)
	}
	return files, err
}

func ListFiles(root string, files []Node) string {
	s := ""
	for _, f := range files {
		relPath := f.Info.Name()
		if root != "" {
			relPath = fmt.Sprintf("%s/%s", root, f.Info.Name())
		}
		s += fmt.Sprintf("%s\n", relPath)

		if f.Children != nil {
			s += ListFiles(relPath, f.Children)
		}
	}
	return s
}
