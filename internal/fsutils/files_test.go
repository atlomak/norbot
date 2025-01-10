package fsutils

import (
	"testing"
)

func TestReadDir(t *testing.T) {

	depth := 1
	root := "test_dir"

	files, err := ReadDir(root, depth)
	if err != nil {
		t.Fatal(err)
	}

	expectedTotal := 10
	gotTotal := len(files)
	for _, f := range files {
		if f.Info.IsDir() {
			gotTotal += len(f.Children)
		}
	}

	if expectedTotal != gotTotal {
		t.Fatalf("expected: %d got: %d", expectedTotal, gotTotal)
	}

}

func TestListFiles(t *testing.T) {

	depth := 1
	root := "test_dir"

	files, err := ReadDir(root, depth)
	if err != nil {
		t.Fatal(err)
	}

	output := ListFiles("", files)
	expectedOutput :=
		`Dir
Dir/test_file_1.txt
Dir/test_file_2.txt
Dir2
Dir2/test_file_1.txt
Dir2/test_file_2.txt
Dir3
test_file.txt
test_file_2.txt
test_file_3.txt
`
	if output != expectedOutput {
		t.Fatalf("\nexpected:%s\ngot:\n%s\n", expectedOutput, output)
	}

}
