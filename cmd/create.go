package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"
	"time"

	"github.com/spf13/cobra"
)

var (
	authorName   string // Global variable to store the author name
	templatesDir string // Global variable to store the templates root directory path
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create [project name]",
	Short: "Create a new C++ project",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		projectName := args[0]

		// Determine the author name
		author := determineAuthor(cmd)

		createProject(projectName, author)
	},
}

func init() {
	defaultTemplateDir := os.ExpandEnv("$HOME/.cpj/templates")
	rootCmd.AddCommand(createCmd)

	// Add flag to allow specifying author name
	createCmd.Flags().StringVarP(&authorName, "author", "a", "", "Author name")

	// Add flag to allow specifying templates root directory
	createCmd.Flags().StringVarP(&templatesDir, "templates", "t", defaultTemplateDir, "Templates root directory path")

	// Set default templates directory
	cobra.OnInitialize(func() {
		if templatesDir == "" {
			templatesDir = "./cpj"
		}
	})
}

func determineAuthor(cmd *cobra.Command) string {
	// Check if author name is specified via flag
	if authorName != "" {
		return authorName
	}

	// Try to get author name from git config
	author, err := gitUserName()
	if err != nil {
		fmt.Println("Warning: Unable to retrieve Git username. Author will be empty.")
		return ""
	}

	return author
}

func gitUserName() (string, error) {
	cmd := exec.Command("git", "config", "user.name")
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

func createProject(projectName, author string) {
	projectPath := filepath.Join(".", projectName)

	// Create project directory
	err := os.Mkdir(projectPath, 0755)
	if err != nil {
		fmt.Println("Error creating project directory:", err)
		return
	}

	defer func() {
		// Handle cleanup in case of error
		if err != nil {
			fmt.Println("Cleaning up...")
			if err := os.RemoveAll(projectPath); err != nil {
				fmt.Println("Error cleaning up:", err)
			}
		}
	}()

	templates := map[string]string{
		"WORKSPACE":             "WORKSPACE.tmpl",
		"BUILD":                 "BUILD.tmpl",
		"src/example.cpp":       "src_example.cpp.tmpl",
		"src/example.h":         "src_example.h.tmpl",
		"cmd/main.cpp":          "cmd_main.cpp.tmpl",
		"test/example_test.cpp": "test_example_test.cpp.tmpl",
	}

	for path, tmplFile := range templates {
		tmplPath := filepath.Join(templatesDir, tmplFile)
		fullPath := filepath.Join(projectPath, path)
		err := os.MkdirAll(filepath.Dir(fullPath), 0755)
		if err != nil {
			fmt.Println("Error creating directory:", err)
			return
		}

		tmpl, err := template.ParseFiles(tmplPath)
		if err != nil {
			fmt.Println("Error parsing template:", err)
			return
		}

		file, err := os.Create(fullPath)
		if err != nil {
			fmt.Println("Error creating file:", err)
			return
		}
		defer file.Close()

		data := struct {
			Filename    string
			Date        string
			Author      string
			ProjectName string
		}{
			Filename:    filepath.Base(fullPath),
			Date:        time.Now().Format("2006-01-02 15:04:05 MST"),
			Author:      author,
			ProjectName: projectName,
		}

		err = tmpl.Execute(file, data)
		if err != nil {
			fmt.Println("Error executing template:", err)
			return
		}
	}

	fmt.Println("Project created successfully:", projectName)
}
