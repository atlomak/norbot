package llm

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sort"

	"github.com/atlomak/norbot/internal/fsutils"
	"github.com/google/generative-ai-go/genai"
)

const (
	query string = `
You are an assistant helping to clean and organize directories efficiently. 
For each file or directory in the provided list, determine an appropriate action:
- Move files to new directories to group them by type, date, name etc.
- Leave system or configuration files (e.g., hidden files, .log, .conf) in their current locations.
- For files that do not match a known category, leave them in place.

Rename files with inconsistent naming to use lowercase and replace spaces with underscores. 
Ensure all actions follow a clear, user-friendly folder structure.
When renaming or moving files, ensure there are no name conflicts in the destination directory:
- If a file with the same name already exists in the target location, add a unique suffix (e.g., "_1", "_2") to the file name
or don't move it at all to prevent overwriting or conflicts.

Examples of actions:
1. If a photo is moved to 'NewDir/NewSubDir', provide:
   { "action": "move", "name": "example.jpg", "result": "NewDir/NewSubDir/example.jpg" }
2. If a file is renamed, provide:
   { "action": "move", "name": "old-file.txt", "result": "old_file.txt" }
3. If a file is left in place, provide:
   { "action": "keep", "name": "example.conf", "result": "example.conf" }

You should also consider additional instructions provided by the user after this prompt.
These instructions might modify or extend your actions for organizing files and folders.
Always prioritize user input if given, but you have to follow rules of schema, and remember about naming conventions.
`

	actionsDesciption = `
The type of operation to be performed. Possible values are:
- 'move': File or directory is moved to a new location or changed name.
- 'keep': File or directory is left unchanged.
`

	nameDescription = `
The original name or path of the file or directory before any action. 
 - For directories, always include a trailing "/" at the end of the name.`

	resultDescription = `
The new name or path of the file or directory after the action. 
  - If the action is "keep", this field should match the "name" field.	
  - For directories, always include a trailing "/" at the end of the name.
`
)

type Action struct {
	Name   string
	Type   string
	Result string
}

type GeminiModel struct {
	model *genai.GenerativeModel
	ctx   context.Context
}

func (m GeminiModel) Query(files fsutils.FileList, prompt string) ([]Action, error) {
	var resp *genai.GenerateContentResponse
	var err error
	if prompt != "" {
		log.Printf("given prompt: %s", prompt)
		query := fmt.Sprintf("additional prompt:\n%s\ndata:\n%s", prompt, files.Details())
		resp, err = m.model.GenerateContent(m.ctx, genai.Text(query))
	} else {
		resp, err = m.model.GenerateContent(m.ctx, genai.Text(files.Details()))
	}
	if err != nil {
		return nil, err
	}

	actions := make([]Action, 0, len(files))
	for _, part := range resp.Candidates[0].Content.Parts {
		if txt, ok := part.(genai.Text); ok {
			var output []map[string]string
			if err := json.Unmarshal([]byte(txt), &output); err != nil {
				return nil, err
			}
			for _, action := range output {
				actions = append(actions, Action{Name: action["name"], Type: action["action"], Result: action["result"]})
			}
		}
	}
	sortActions(actions)
	return actions, nil
}

func sortActions(actions []Action) {
	sort.Slice(actions, func(i, j int) bool {
		return actions[i].Name < actions[j].Name
	})
}

// func safeguardMissingCreates(files fsutils.FileList, actins []Action) []Action {
// 	entries := files.String()
// 	check := map[string]int

// 	for _, v := range entries {
// 		check[v] = 1
// 	}

// 	for _, v := range actins {
// 		path := strings.Split(strings.TrimSuffix(file, "/"), "/")
// 		for
// 	}

// }

func InitGeminiModel(client *genai.Client, ctx context.Context) *GeminiModel {

	model := client.GenerativeModel("gemini-1.5-flash")
	model.ResponseMIMEType = "application/json"
	model.ResponseSchema = &genai.Schema{
		Type: genai.TypeArray,
		Items: &genai.Schema{
			Type: genai.TypeObject,
			Properties: map[string]*genai.Schema{
				"action": {
					Type:        genai.TypeString,
					Enum:        []string{"move", "keep"},
					Description: actionsDesciption,
				},
				"name": {
					Type:        genai.TypeString,
					Description: nameDescription,
				},
				"result": {
					Type:        genai.TypeString,
					Description: resultDescription,
				},
			},
		},
	}
	model.SystemInstruction = genai.NewUserContent(genai.Text(query))
	return &GeminiModel{
		model: model,
		ctx:   ctx,
	}
}
