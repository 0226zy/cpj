package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

// buildCmd represents the build command
var buildCmd = &cobra.Command{
	Use:   "build [binary name]",
	Short: "Build the C++ project",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		binaryName := args[0]
		buildProject(binaryName)
	},
}

func init() {
	rootCmd.AddCommand(buildCmd)
}

func buildProject(binaryName string) {
	cmd := exec.Command("bazel", "build", fmt.Sprintf("//:%s", binaryName))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		fmt.Println("Error building project:", err)
		return
	}

	fmt.Println("Project built successfully")
}
