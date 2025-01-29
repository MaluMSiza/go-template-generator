package main

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"os/exec"

	"github.com/MaluMSiza/go-template-generator/internal/generator"
	"github.com/spf13/cobra"
	"io/ioutil"
)

// Config struct para ler as configurações do YAML
type Config struct {
	ModuleName   string `yaml:"module_name"`
	OutputDir    string `yaml:"output_dir"`
	TemplateRepo string `yaml:"template_repo"`
}

// Função para ler o arquivo YAML
func readConfig(filename string) (*Config, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}

func main() {
	var rootCmd = &cobra.Command{
		Use:   "generator",
		Short: "Gerador de templates Go",
	}

	var grpcCmd = &cobra.Command{
		Use:   "grpc-template [config file]",
		Short: "Gera um template de microserviço gRPC",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			configFile := args[0] // Arquivo de configuração passado como argumento
			config, err := readConfig(configFile)
			if err != nil {
				log.Fatalf("Erro ao ler o arquivo de configuração: %v", err)
			}

			// Clona o template no diretório especificado
			if err := generator.CloneTemplate(config.TemplateRepo, config.OutputDir); err != nil {
				log.Fatalf("Erro ao gerar o template: %v", err)
			}

			// Inicializa um novo módulo Go no diretório de saída
			if err := os.Chdir(config.OutputDir); err != nil {
				log.Fatalf("Erro ao mudar para o diretório de saída: %v", err)
			}
			cmdInit := exec.Command("go", "mod", "init", config.ModuleName)
			if err := cmdInit.Run(); err != nil {
				log.Fatalf("Erro ao inicializar o módulo Go: %v", err)
			}

			// Substitui placeholders
			if err := generator.ReplacePlaceholders(config.OutputDir, "github.com/MaluMSiza/gRPC-microservice-template", config.ModuleName); err != nil {
				log.Fatalf("Erro ao substituir placeholders: %v", err)
			}

			fmt.Println("Template gerado com sucesso!")
		},
	}

	rootCmd.AddCommand(grpcCmd)

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
