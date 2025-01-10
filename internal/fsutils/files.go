package fsutils

import (
	"fmt"
	"io/fs"
	"log"
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

func listFiles(root string, depth int) ([]Node, error) {
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
		log.Println(info.Name())

		var node Node
		if depth != 0 && info.IsDir() {
			children, err := listFiles(fmt.Sprintf("%s/%s", root, info.Name()), depth-1)
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

// func filesToJson(files []FileInfo) ([]byte, error) {
// 	req, err := json.Marshal(files)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return req, nil
// }
