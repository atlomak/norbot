package ui

import (
	"log"
	"strings"

	"github.com/atlomak/norbot/internal/fsutils"
	"github.com/atlomak/norbot/internal/llm"
	"github.com/charmbracelet/bubbles/list"
)

func generateActionMapWithDirs(actions []llm.Action) map[string]llm.Action {
	checkMap := make(map[string]llm.Action)

	// Prevent conflicts
	for _, action := range actions {
		log.Printf("add action to map: %s", action)
		checkMap[action.Result] = action
	}

	// Add dirs if not exist, but prevent overrides
	for result := range checkMap {
		path := strings.Split(strings.TrimSuffix(result, "/"), "/")
		if len(path) > 1 {
			parenFolders := path[0 : len(path)-1]
			log.Printf("parents: %v", parenFolders)
			name := ""
			for _, parent := range parenFolders {
				name += parent + "/"
				if v, ok := checkMap[name]; !ok {
					log.Printf("create dir: %s", name)
					checkMap[name] = llm.Action{Type: "create", Result: name}
				} else {
					log.Printf("exists dir: %v", v)
				}
			}
		}
	}

	results := make(map[string]llm.Action)
	for _, v := range checkMap {
		if v.Name == "" {
			results[v.Result] = v
		} else {
			results[v.Name] = v
		}
	}
	return results
}

func maxDepth(actions []llm.Action) int {
	maxDepth := 0
	for _, action := range actions {
		path := strings.Split(strings.TrimSuffix(action.Result, "/"), "/")
		parents := len(path) - 1
		if parents > maxDepth {
			maxDepth = parents
		}
	}
	return maxDepth
}

func filesToItems(files fsutils.FileList) []list.Item {
	items := make([]list.Item, 0, len(files))

	s := strings.Split(files.String(), "\n")
	s = s[:len(s)-1] // because of newline at the end of string

	for _, file := range s {
		items = append(items, item{name: file})
	}
	return items
}
