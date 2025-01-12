package fsutils

import (
	"fmt"
	"os"
)

func MoveFile(currentFileName, resultFileName string) error {
	if _, err := os.Stat(currentFileName); os.IsNotExist(err) {
		return fmt.Errorf("source file does not exist: %s", currentFileName)
	}

	if _, err := os.Stat(resultFileName); err == nil {
		return fmt.Errorf("destination file already exists: %s", resultFileName)
	} else if !os.IsNotExist(err) {
		return fmt.Errorf("failed to check destination file: %w", err)
	}

	err := os.Rename(currentFileName, resultFileName)
	if err != nil {
		return fmt.Errorf("failed to move file: %w", err)
	}

	return nil
}

func CreateDir(dirName string) error {
	if _, err := os.Stat(dirName); !os.IsNotExist(err) {
		if err != nil {
			return fmt.Errorf("failed to check directory: %w", err)
		}
	}

	err := os.MkdirAll(dirName, 0755)
	if err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	return nil
}
