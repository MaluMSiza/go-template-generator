package generator

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// CloneTemplate baixa o template de um repositório Git
func CloneTemplate(templateRepo, outputDir string) error {
	if _, err := os.Stat(outputDir); !os.IsNotExist(err) {
		return fmt.Errorf("o diretório %s já existe", outputDir)
	}

	cmd := exec.Command("git", "clone", templateRepo, outputDir)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("falha ao clonar o repositório: %v", err)
	}

	if err := os.RemoveAll(filepath.Join(outputDir, ".git")); err != nil {
		return fmt.Errorf("falha ao remover .git: %v", err)
	}

	fmt.Printf("Template clonado com sucesso em %s\n", outputDir)
	return nil
}

// ReplacePlaceholders substitui placeholders nos arquivos do template
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
			return fmt.Errorf("falha ao ler o arquivo %s: %v", path, err)
		}

		newContent := strings.ReplaceAll(string(content), placeholder, value)

		if err := os.WriteFile(path, []byte(newContent), info.Mode()); err != nil {
			return fmt.Errorf("falha ao escrever o arquivo %s: %v", path, err)
		}

		return nil
	})
}
