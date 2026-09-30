package snapshot

import (
	"os"
	"path/filepath"
	"strings"
)

func projectRoot(executablePath, workingDirectory string) string {
	if executablePath != "" {
		executableDirectory := filepath.Dir(executablePath)
		if strings.EqualFold(filepath.Base(executableDirectory), "builds") {
			return filepath.Dir(executableDirectory)
		}
		if _, err := os.Stat(filepath.Join(executableDirectory, "go.mod")); err == nil {
			return executableDirectory
		}
	}

	if strings.EqualFold(filepath.Base(workingDirectory), "builds") {
		return filepath.Dir(workingDirectory)
	}
	return workingDirectory
}

func ProjectRoot() (string, error) {
	executablePath, _ := os.Executable()
	workingDirectory, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return projectRoot(executablePath, workingDirectory), nil
}
