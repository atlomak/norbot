package llm

import (
	"context"
	"encoding/json"
	"sort"

	"github.com/atlomak/norbot/internal/fsutils"
	"github.com/google/generative-ai-go/genai"
)

const (
	query string = `
You are an assistant helping to clean and organize directories efficiently. 
For each file or directory in the provided list, determine an appropriate action:
- If it is a photo, move it to a 'Photos/YYYY/MM' folder based on its creation date. 
  Create the 'Photos/YYYY' and 'Photos/YYYY/MM' folders if they do not already exist.
- For documents (e.g., PDFs, DOCs), move them to a 'Documents' folder, grouped by type. 
  Create the 'Documents' folder if it does not exist.
- For media files like videos or music, organize them into 'Videos' or 'Music' folders. 
  Create these folders if they do not exist.
- For archives (e.g., ZIP, RAR), move them to an 'Archives' folder. 
  Create the 'Archives' folder if it does not exist.
- Leave system or configuration files (e.g., hidden files, .log, .conf) in their current locations.
- For files that do not match a known category, leave them in place.
- Whenever a folder is created, include a record for that action.

Rename files with inconsistent naming to use lowercase and replace spaces with underscores. 
Ensure all actions follow a clear, user-friendly folder structure.
When renaming or moving files, ensure there are no name conflicts in the destination directory:
- If a file with the same name already exists in the target location, add a unique suffix (e.g., "_1", "_2") to the file name
or don't move it at all to prevent overwriting or conflicts.

Examples of actions:
1. If a photo is moved to 'Photos/2025/01', provide:
   { "action": "create", "name": "", "result": "Photos/" }
   { "action": "create", "name": "", "result": "Photos/2025/" }
   { "action": "create", "name": "", "result": "Photos/2025/01/" }
   { "action": "move", "name": "example.jpg", "result": "Photos/2025/01/example.jpg" }
2. If a file is renamed, provide:
   { "action": "move", "name": "old file.txt", "result": "old_file.txt" }
3. If a file is left in place, provide:
   { "action": "keep", "name": "example.conf", "result": "example.conf" }
`

	actionsDesciption = `
The type of operation to be performed. Possible values are:
- 'move': File or directory is moved to a new location or changed name.
- 'keep': File or directory is left unchanged.
- 'create': A new folder is created.`

	nameDescription = `
The original name or path of the file or directory before any action. 
 - If the action is "create", this field should be an empty string as there is no original path.
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
	model  *genai.GenerativeModel
	prompt string
	ctx    context.Context
}

func (m GeminiModel) Query(files fsutils.FileList, prompt string) ([]Action, error) {
	var resp *genai.GenerateContentResponse
	var err error
	if prompt != "" {
		resp, err = m.model.GenerateContent(m.ctx, genai.Text(m.prompt), genai.Text(files.Details()))
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
					Enum:        []string{"move", "keep", "create"},
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
