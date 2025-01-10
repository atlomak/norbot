package llm

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/atlomak/norbot/internal/fsutils"
	tea "github.com/charmbracelet/bubbletea"
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
	Type   string
	Result string
}

type GeminiMsg struct {
	Actions map[string]Action
	Err     error
}

type GeminiModel struct {
	model  *genai.GenerativeModel
	prompt string
	ctx    context.Context
}

func (m GeminiModel) Query(files []fsutils.Node) tea.Cmd {
	return func() tea.Msg {
		s := ""
		for _, file := range files {
			s += fmt.Sprintf(
				"%-20s %-10d %-10s isDir %t %s\n",
				file.Info.Name(), file.Info.Size(), file.Info.Mode().String(), file.Info.IsDir(), file.Info.ModTime().Format(time.RFC1123),
			)
		}
		resp, err := m.model.GenerateContent(m.ctx, genai.Text(m.prompt), genai.Text(s))
		if err != nil {
			return GeminiMsg{Err: err}
		}

		msg := GeminiMsg{
			Actions: map[string]Action{},
		}
		for _, part := range resp.Candidates[0].Content.Parts {
			if txt, ok := part.(genai.Text); ok {
				var actions []map[string]string
				if err := json.Unmarshal([]byte(txt), &actions); err != nil {
					return GeminiMsg{Err: err}
				}
				for _, action := range actions {
					// log.Println("Hello world")
					msg.Actions[action["name"]] = Action{Type: action["action"], Result: action["output"]}
				}
				return msg
			}
		}
		return GeminiMsg{}
	}
}

func InitGeminiModel(client *genai.Client, ctx context.Context) GeminiModel {

	genAI := client.GenerativeModel("gemini-1.5-flash")
	genAI.ResponseMIMEType = "application/json"
	genAI.ResponseSchema = &genai.Schema{
		Type: genai.TypeArray,
		Items: &genai.Schema{
			Type: genai.TypeObject,
			Properties: map[string]*genai.Schema{
				"action": {
					Type: genai.TypeString,
					Enum: []string{"move", "rename", "keep", "create"},
				},
				"name": {
					Type: genai.TypeString,
				},
				"output": {
					Type:        genai.TypeString,
					Description: "Output of an action. Eg. action: rename, output: new name",
				},
			},
		},
	}
	return GeminiModel{
		model:  genAI,
		prompt: query,
		ctx:    ctx,
	}
}
