package generator

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// CloneTemplate download a template of Git repository
func CloneTemplate(templateRepo, outputDir string) error {
	if _, err := os.Stat(outputDir); !os.IsNotExist(err) {
		return fmt.Errorf("dir %s already exists", outputDir)
	}

	cmd := exec.Command("git", "clone", templateRepo, outputDir)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to clone the repository: %v", err)
	}

	if err := os.RemoveAll(filepath.Join(outputDir, ".git")); err != nil {
		return fmt.Errorf("failed to remove .git: %v", err)
	}

	fmt.Printf("Template cloned with success on %s\n", outputDir)
	return nil
}

// ReplacePlaceholders change placeholders
func ReplacePlaceholders(outputDir, placeholder, value string) error {
	return filepath.Walk(outputDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		content, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("error to read the file %s: %v", path, err)
		}

		newContent := strings.ReplaceAll(string(content), placeholder, value)

		if err := os.WriteFile(path, []byte(newContent), info.Mode()); err != nil {
			return fmt.Errorf("error to write the file %s: %v", path, err)
		}

		return nil
	})
}
