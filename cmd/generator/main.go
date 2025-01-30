package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/MaluMSiza/go-template-generator/internal/generator"
	"github.com/spf13/cobra"
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "generator",
		Short: "Templates generator Go",
	}

	var (
		moduleName   string
		outputDir    string
		templateRepo string
	)

	var grpcCmd = &cobra.Command{
		Use:   "grpc-template",
		Short: "Template for microservice gRPC",
		Run: func(cmd *cobra.Command, args []string) {
			// Check if dir already exists
			if _, err := os.Stat(outputDir); !os.IsNotExist(err) {
				log.Fatalf("O diretório %s já existe", outputDir)
			}

			// Clone template
			if err := generator.CloneTemplate(templateRepo, outputDir); err != nil {
				log.Fatalf("Error to generate template: %v", err)
			}

			// Init new Go module on output dir
			if err := os.Chdir(outputDir); err != nil {
				log.Fatalf("Erro to change output dir: %v", err)
			}
			cmdInit := exec.Command("go", "mod", "init", moduleName)
			if err := cmdInit.Run(); err != nil {
				log.Fatalf("Erro to inialize Go module: %v", err)
			}

			// Change placeholders
			if err := generator.ReplacePlaceholders(outputDir, "github.com/MaluMSiza/gRPC-microservice-template", moduleName); err != nil {
				log.Fatalf("Error to change placeholders: %v", err)
			}

			fmt.Println("Success to generate template!")
		},
	}

	// Flags to grpc-template command
	grpcCmd.Flags().StringVarP(&moduleName, "module", "m", "", "Name of go module (ex: github.com/my-user/my-project)")
	grpcCmd.Flags().StringVarP(&outputDir, "output", "o", "./my-project-grpc", "Output directory")
	grpcCmd.Flags().StringVarP(&templateRepo, "repo", "r", "https://github.com/MaluMSiza/gRPC-microservice-template.git", "Template repository")

	// Mark flag --module required
	grpcCmd.MarkFlagRequired("module")

	rootCmd.AddCommand(grpcCmd)

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
