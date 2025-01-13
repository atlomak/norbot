package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/atlomak/norbot/internal/tui"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

func main() {
	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(os.Getenv("GEMINI_API_KEY")))
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	if len(os.Getenv("DEBUG")) > 0 {
		f, err := tea.LogToFile("debug.log", "debug")
		if err != nil {
			fmt.Println("fatal:", err)
			os.Exit(1)
		}
		defer f.Close()
	}

	p := tea.NewProgram(tui.InitModel(client, ctx))
	if _, err := p.Run(); err != nil {
		fmt.Printf("Geez, there's been an error: %v", err)
		os.Exit(1)
	}
}
