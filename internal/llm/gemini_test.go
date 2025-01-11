package llm

import (
	"context"
	"log"
	"os"
	"sort"
	"testing"

	"github.com/atlomak/norbot/internal/fsutils"
	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

func generateOutput(files fsutils.FileList) []Action {
	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(os.Getenv("GEMINI_API_KEY")))
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()
	model := InitGeminiModel(client, ctx)

	output, err := model.Query(files, "")
	if err != nil {
		log.Fatal(err)
	}
	return output
}

func TestOutputSize(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}

	files, err := fsutils.ReadDir("test_dir", 0)
	if err != nil {
		t.Fatal(err)
	}

	output := generateOutput(files)

	expectedSize := len(files)
	gotSize := len(output)
	if expectedSize != gotSize {
		t.Fatalf("expected: %d, got: %d", expectedSize, gotSize)
	}
}

func TestOutputFileNames(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}

	files, err := fsutils.ReadDir("test_dir", 0)
	if err != nil {
		t.Fatal(err)
	}

	output := generateOutput(files)
	sort.Slice(output, func(i, j int) bool {
		return output[i].Name < output[j].Name
	})

	expectedNames := []string{
		"Dir/",
		"Dir2/",
		"test_file.txt",
		"test_file_2.txt",
		"test_file_3.txt",
	}

	expectedSize := len(expectedNames)
	gotSize := len(output)
	if expectedSize != gotSize {
		t.Fatalf("expected: %d, got: %d", expectedSize, gotSize)
	}

	for i, name := range expectedNames {
		t.Run(name, func(t *testing.T) {
			outputName := output[i].Name
			if name != outputName {
				t.Fatalf("expected: %s, got: %s\n", name, outputName)
			}
		})
	}

}
