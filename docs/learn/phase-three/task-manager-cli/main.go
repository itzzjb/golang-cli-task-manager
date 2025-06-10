package main

import (
	"fmt"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "task",
	Short: "Task manager cli",
	Long:  "A simple command-line interface for managing tasks.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Use 'task --help' for commands.")
	},
}

var addCmd = &cobra.Command{
	Use:   "add [task]",
	Short: "Add a new task",
	Long:  "Add a new task to the task list. Provide the task description as an argument.",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			fmt.Println("Please provide a task description.")
			return
		}
		task := args[0]
		fmt.Printf("Task added: %s\n", task)
	},
}

func main() {
	rootCmd.AddCommand(addCmd)
	err := rootCmd.Execute()
	if err != nil {
		return
	}
}

// go build -o task-manager-cli main.go
// ./task-manager-cli add "Learn Go with Cobra"
