package ui

import (
	"reflect"
	"testing"

	"github.com/atlomak/norbot/internal/llm"
)

func TestActionsToMap(t *testing.T) {
	tests := []struct {
		name     string
		actions  []llm.Action
		expected map[string]llm.Action
	}{
		{
			name: "Single action, no hierarchy",
			actions: []llm.Action{
				{Type: "move", Result: "file.txt", Name: "file.txt"},
			},
			expected: map[string]llm.Action{
				"file.txt": {Type: "move", Result: "file.txt", Name: "file.txt"},
			},
		},
		{
			name: "Multiple actions with directory hierarchy",
			actions: []llm.Action{
				{Type: "move", Result: "dir/file.txt", Name: "file.txt"},
			},
			expected: map[string]llm.Action{
				"dir/":     {Type: "create", Result: "dir/"},
				"file.txt": {Type: "move", Result: "dir/file.txt", Name: "file.txt"},
			},
		},
		{
			name: "Prevent overriding existing actions",
			actions: []llm.Action{
				{Type: "keep", Result: "dir/", Name: "dir/"},
				{Type: "move", Result: "dir/file.txt", Name: "file.txt"},
			},
			expected: map[string]llm.Action{
				"dir/":     {Type: "keep", Result: "dir/", Name: "dir/"},
				"file.txt": {Type: "move", Result: "dir/file.txt", Name: "file.txt"},
			},
		},
		{
			name: "Complex hierarchy with multiple actions",
			actions: []llm.Action{
				{Type: "move", Result: "root/dir1/file1.txt", Name: "file1.txt"},
				{Type: "move", Result: "root/dir2/file2.txt", Name: "file2.txt"},
			},
			expected: map[string]llm.Action{
				"root/":      {Type: "create", Result: "root/"},
				"root/dir1/": {Type: "create", Result: "root/dir1/"},
				"file1.txt":  {Type: "move", Result: "root/dir1/file1.txt", Name: "file1.txt"},
				"root/dir2/": {Type: "create", Result: "root/dir2/"},
				"file2.txt":  {Type: "move", Result: "root/dir2/file2.txt", Name: "file2.txt"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := generateActionMapWithDirs(tt.actions)
			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("actionsToMap() = %v, want %v", got, tt.expected)
			}
		})
	}
}
