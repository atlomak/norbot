package fsutils

import "testing"

func TestListFilesSize(t *testing.T) {

	depth := 1
	root := "test_dir"

	files, err := listFiles(root, depth)
	if err != nil {
		t.Fatal(err)
	}

	expectedTotal := 9
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
