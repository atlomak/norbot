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

type FileList []Node

func ReadDir(root string, depth int) (FileList, error) {
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

func (l FileList) String() string {
	return listFiles("", l)
}

func (l FileList) Details() string {
	return listFilesDetails("", l)
}

func listFiles(root string, files []Node) string {
	s := ""
	for _, f := range files {
		relPath := f.Info.Name()
		if root != "" {
			relPath = fmt.Sprintf("%s/%s", root, f.Info.Name())
		}

		name := relPath
		if f.Info.IsDir() {
			name += "/"
		}

		s += fmt.Sprintf("%s\n", name)

		if f.Children != nil {
			s += listFiles(relPath, f.Children)
		}
	}
	return s
}

func listFilesDetails(root string, files []Node) string {
	timeFormat := "Jan _2  2006"
	s := ""
	for _, f := range files {
		relPath := f.Info.Name()
		if root != "" {
			relPath = fmt.Sprintf("%s/%s", root, f.Info.Name())
		}

		name := relPath
		if f.Info.IsDir() {
			name += "/"
		}

		s += fmt.Sprintf("%8d %s %s\n",
			f.Info.Size(),
			f.Info.ModTime().Format(timeFormat),
			name,
		)

		if f.Children != nil {
			s += listFilesDetails(relPath, f.Children)
		}
	}
	return s
}
