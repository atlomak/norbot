package llm

import (
	"context"
	"encoding/json"
	"sort"

	"github.com/atlomak/norbot/internal/fsutils"
	"github.com/google/generative-ai-go/genai"
)

var query string = `
Organize files in a directory. Rename files for consistency.
For photos, group into Year/Month folders using creation dates. For other files,
organize by type in subdirectories.
If you don't know what to do with specific file, just leave it as is.
Respond with provided schema.
`

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

func InitGeminiModel(client *genai.Client, ctx context.Context) GeminiModel {

	model := client.GenerativeModel("gemini-1.5-flash")
	model.ResponseMIMEType = "application/json"
	model.ResponseSchema = &genai.Schema{
		Type: genai.TypeArray,
		Items: &genai.Schema{
			Type: genai.TypeObject,
			Properties: map[string]*genai.Schema{
				"action": {
					Type: genai.TypeString,
					Enum: []string{"move", "rename", "keep"},
				},
				"name": {
					Type:        genai.TypeString,
					Description: "Original name of the file or directory before any action. If it is dir name, keep trailing '/'",
				},
				"result": {
					Type:        genai.TypeString,
					Description: "New name after action. If file is moved to a new directory, it's name should be with additonal path. If action is keep, the field should be equal to the name.",
				},
			},
		},
	}
	model.SystemInstruction = genai.NewUserContent(genai.Text(query))
	return GeminiModel{
		model: model,
		ctx:   ctx,
	}
}
