# Step-by-Step CLI Building Tutorial

## 🚀 Building Your First Go CLI Tool

This hands-on tutorial will guide you through building a complete task management CLI tool. We'll build it step by step, explaining each concept as we go.

## Step 1: Project Setup (15 minutes)

### 1.1 Create the Project Structure
```bash
mkdir task-manager-cli
cd task-manager-cli

# Initialize Go module
go mod init task-manager-cli

# Create directory structure
mkdir -p cmd internal/task internal/config tests testdata
```

### 1.2 Install Dependencies
```bash
# CLI framework
go get github.com/spf13/cobra@latest

# Configuration management
go get github.com/spf13/viper@latest

# Interactive prompts
go get github.com/AlecAivazis/survey/v2@latest

# Colored output
go get github.com/fatih/color@latest
```

**🎯 What you learned:**
- Go module initialization
- Project structure best practices
- Dependency management with `go get`

## Step 2: Create the Main Entry Point (10 minutes)

### 2.1 Create `main.go`
```go
package main

import (
	"os"
	"task-manager-cli/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
```

**🎯 What you learned:**
- Main function as application entry point
- Error handling with exit codes
- Package imports and structure

## Step 3: Setup the Root Command (20 minutes)

### 3.1 Create `cmd/root.go`
```go
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

var rootCmd = &cobra.Command{
	Use:   "task-manager",
	Short: "A simple task management CLI",
	Long: `Task Manager helps you organize your daily tasks.
	
You can add, list, complete, and delete tasks right from your terminal.`,
	Version: "1.0.0",
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)
	
	// Global flags
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file")
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "verbose output")
	
	// Bind flags to viper
	viper.BindPFlag("verbose", rootCmd.PersistentFlags().Lookup("verbose"))
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)
		
		viper.AddConfigPath(home)
		viper.AddConfigPath(".")
		viper.SetConfigType("yaml")
		viper.SetConfigName(".task-manager")
	}
	
	viper.AutomaticEnv()
	
	if err := viper.ReadInConfig(); err == nil {
		if viper.GetBool("verbose") {
			fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
		}
	}
}
```

### 3.2 Test Your Root Command
```bash
go run . --help
```

**🎯 What you learned:**
- Cobra command structure
- Global flags vs command-specific flags
- Configuration initialization
- Help text generation

## Step 4: Define Task Data Structures (25 minutes)

### 4.1 Create `internal/task/task.go`
```go
package task

import (
	"time"
)

type Status string
type Priority string

const (
	StatusPending   Status = "pending"
	StatusCompleted Status = "completed"
)

const (
	PriorityLow    Priority = "low"
	PriorityMedium Priority = "medium"
	PriorityHigh   Priority = "high"
)

type Task struct {
	ID          int        `json:"id"`
	Description string     `json:"description"`
	Priority    Priority   `json:"priority"`
	Status      Status     `json:"status"`
	DueDate     *time.Time `json:"due_date,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	CompletedAt *time.Time `json:"completed_at,omitempty"`
	Tags        []string   `json:"tags,omitempty"`
}

// IsOverdue returns true if the task is past its due date
func (t *Task) IsOverdue() bool {
	if t.DueDate == nil || t.Status == StatusCompleted {
		return false
	}
	return time.Now().After(*t.DueDate)
}

// Complete marks the task as completed
func (t *Task) Complete() {
	t.Status = StatusCompleted
	now := time.Now()
	t.CompletedAt = &now
}
```

**🎯 What you learned:**
- Go struct definitions and tags
- Constants and custom types
- JSON serialization tags
- Pointer usage for optional fields
- Method receivers

## Step 5: Create Storage Layer (30 minutes)

### 5.1 Create `internal/task/storage.go`
```go
package task

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

type Storage interface {
	GetTasks() ([]Task, error)
	SaveTasks(tasks []Task) error
	AddTask(task Task) error
}

type FileStorage struct {
	filePath string
}

func NewFileStorage() *FileStorage {
	storagePath := viper.GetString("storage.path")
	if storagePath == "" {
		storagePath = "tasks.json"
	}
	
	return &FileStorage{
		filePath: storagePath,
	}
}

func (fs *FileStorage) GetTasks() ([]Task, error) {
	// Check if file exists
	if _, err := os.Stat(fs.filePath); os.IsNotExist(err) {
		return []Task{}, nil
	}
	
	// Read file
	data, err := os.ReadFile(fs.filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read tasks file: %w", err)
	}
	
	// Handle empty file
	if len(data) == 0 {
		return []Task{}, nil
	}
	
	// Parse JSON
	var tasks []Task
	if err := json.Unmarshal(data, &tasks); err != nil {
		return nil, fmt.Errorf("failed to parse tasks file: %w", err)
	}
	
	return tasks, nil
}

func (fs *FileStorage) SaveTasks(tasks []Task) error {
	// Ensure directory exists
	dir := filepath.Dir(fs.filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}
	
	// Marshal to JSON with indentation
	data, err := json.MarshalIndent(tasks, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal tasks: %w", err)
	}
	
	// Write to file
	if err := os.WriteFile(fs.filePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write tasks file: %w", err)
	}
	
	return nil
}

func (fs *FileStorage) AddTask(task Task) error {
	tasks, err := fs.GetTasks()
	if err != nil {
		return fmt.Errorf("failed to load existing tasks: %w", err)
	}
	
	// Assign new ID
	task.ID = len(tasks) + 1
	
	// Add to tasks slice
	tasks = append(tasks, task)
	
	// Save updated tasks
	return fs.SaveTasks(tasks)
}
```

**🎯 What you learned:**
- Interface design in Go
- File I/O operations
- JSON marshaling and unmarshaling
- Error wrapping with `fmt.Errorf`
- Directory creation

## Step 6: Build the Add Command (35 minutes)

### 6.1 Create `cmd/add.go`
```go
package cmd

import (
	"fmt"
	"strings"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"task-manager-cli/internal/task"
)

var addCmd = &cobra.Command{
	Use:   "add [task description]",
	Short: "Add a new task",
	Long:  `Add a new task to your task list.`,
	Args:  cobra.MaximumNArgs(1),
	RunE:  addTask,
}

var (
	priority    string
	dueDate     string
	interactive bool
)

func init() {
	rootCmd.AddCommand(addCmd)
	
	addCmd.Flags().StringVarP(&priority, "priority", "p", "medium", "Task priority")
	addCmd.Flags().StringVarP(&dueDate, "due", "d", "", "Due date (YYYY-MM-DD)")
	addCmd.Flags().BoolVarP(&interactive, "interactive", "i", false, "Interactive mode")
}

func addTask(cmd *cobra.Command, args []string) error {
	var taskDescription string
	var err error
	
	// Get task description
	if interactive || len(args) == 0 {
		prompt := &survey.Input{
			Message: "Enter task description:",
		}
		if err := survey.AskOne(prompt, &taskDescription); err != nil {
			return err
		}
	} else {
		taskDescription = strings.Join(args, " ")
	}
	
	if taskDescription == "" {
		return fmt.Errorf("task description cannot be empty")
	}
	
	// Parse due date if provided
	var due *time.Time
	if dueDate != "" {
		parsedDate, err := time.Parse("2006-01-02", dueDate)
		if err != nil {
			return fmt.Errorf("invalid due date format: %s", dueDate)
		}
		due = &parsedDate
	}
	
	// Create task
	newTask := task.Task{
		Description: taskDescription,
		Priority:    task.Priority(priority),
		DueDate:     due,
		CreatedAt:   time.Now(),
		Status:      task.StatusPending,
	}
	
	// Save task
	storage := task.NewFileStorage()
	if err := storage.AddTask(newTask); err != nil {
		return fmt.Errorf("failed to add task: %w", err)
	}
	
	// Success message
	successColor := color.New(color.FgGreen, color.Bold)
	successColor.Printf("✅ Task added successfully!\n")
	fmt.Printf("   Description: %s\n", newTask.Description)
	fmt.Printf("   Priority: %s\n", newTask.Priority)
	
	return nil
}
```

### 6.2 Test the Add Command
```bash
# Try different ways to add tasks
go run . add "Learn Go CLI programming"
go run . add "Build awesome CLI" --priority high
go run . add --interactive
```

**🎯 What you learned:**
- Command arguments and flags
- Interactive prompts with Survey
- Colored output with color package
- Time parsing in Go
- Command validation

## Step 7: Build the List Command (30 minutes)

### 7.1 Create `cmd/list.go`
```go
package cmd

import (
	"fmt"
	"os"
	"strconv"
	"text/tabwriter"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"task-manager-cli/internal/task"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List tasks",
	Long:  `List all tasks or filter by status and priority.`,
	RunE:  listTasks,
}

var (
	listStatus   string
	listPriority string
	listAll      bool
)

func init() {
	rootCmd.AddCommand(listCmd)
	
	listCmd.Flags().StringVarP(&listStatus, "status", "s", "", "Filter by status")
	listCmd.Flags().StringVarP(&listPriority, "priority", "p", "", "Filter by priority")
	listCmd.Flags().BoolVarP(&listAll, "all", "a", false, "Show all tasks")
}

func listTasks(cmd *cobra.Command, args []string) error {
	storage := task.NewFileStorage()
	tasks, err := storage.GetTasks()
	if err != nil {
		return fmt.Errorf("failed to load tasks: %w", err)
	}
	
	if len(tasks) == 0 {
		infoColor := color.New(color.FgYellow)
		infoColor.Println("📝 No tasks found. Use 'task-manager add' to create your first task!")
		return nil
	}
	
	// Filter tasks
	filteredTasks := filterTasks(tasks, listStatus, listPriority, listAll)
	
	if len(filteredTasks) == 0 {
		infoColor := color.New(color.FgYellow)
		infoColor.Println("🔍 No tasks match the specified filters.")
		return nil
	}
	
	// Display tasks
	displayTasks(filteredTasks)
	
	return nil
}

func filterTasks(tasks []task.Task, status, priority string, showAll bool) []task.Task {
	var filtered []task.Task
	
	for _, t := range tasks {
		// Filter by status
		if status != "" && string(t.Status) != status {
			continue
		}
		
		// Show only pending by default
		if !showAll && status == "" && t.Status == task.StatusCompleted {
			continue
		}
		
		// Filter by priority
		if priority != "" && string(t.Priority) != priority {
			continue
		}
		
		filtered = append(filtered, t)
	}
	
	return filtered
}

func displayTasks(tasks []task.Task) {
	headerColor := color.New(color.FgCyan, color.Bold)
	headerColor.Println("\n📋 Your Tasks:")
	
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	
	fmt.Fprintln(w, "ID\tDescription\tPriority\tStatus\tCreated")
	fmt.Fprintln(w, "--\t-----------\t--------\t------\t-------")
	
	for i, t := range tasks {
		id := strconv.Itoa(i + 1)
		description := truncateString(t.Description, 30)
		priority := formatPriority(string(t.Priority))
		status := formatStatus(string(t.Status))
		created := t.CreatedAt.Format("2006-01-02")
		
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n", 
			id, description, priority, status, created)
	}
	
	w.Flush()
	fmt.Println()
}

func formatPriority(priority string) string {
	switch priority {
	case "high":
		return color.RedString("🔴 " + priority)
	case "medium":
		return color.YellowString("🟡 " + priority)
	case "low":
		return color.GreenString("🟢 " + priority)
	default:
		return priority
	}
}

func formatStatus(status string) string {
	switch status {
	case "completed":
		return color.GreenString("✅ done")
	case "pending":
		return color.YellowString("⏳ pending")
	default:
		return status
	}
}

func truncateString(str string, maxLen int) string {
	if len(str) <= maxLen {
		return str
	}
	return str[:maxLen-3] + "..."
}
```

### 7.2 Test the List Command
```bash
go run . list
go run . list --status pending
go run . list --priority high
go run . list --all
```

**🎯 What you learned:**
- Slice filtering and iteration
- Table formatting with tabwriter
- String manipulation and truncation
- Conditional colored output
- Advanced flag usage

## Step 8: Build Complete and Delete Commands (25 minutes)

### 8.1 Create `cmd/complete.go`
```go
package cmd

import (
	"fmt"
	"strconv"
	"time"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"task-manager-cli/internal/task"
)

var completeCmd = &cobra.Command{
	Use:   "complete <task-id>",
	Short: "Mark a task as completed",
	Args:  cobra.ExactArgs(1),
	RunE:  completeTask,
}

func init() {
	rootCmd.AddCommand(completeCmd)
}

func completeTask(cmd *cobra.Command, args []string) error {
	taskID, err := strconv.Atoi(args[0])
	if err != nil {
		return fmt.Errorf("invalid task ID: %s", args[0])
	}
	
	storage := task.NewFileStorage()
	tasks, err := storage.GetTasks()
	if err != nil {
		return fmt.Errorf("failed to load tasks: %w", err)
	}
	
	if taskID < 1 || taskID > len(tasks) {
		return fmt.Errorf("task ID %d not found", taskID)
	}
	
	targetTask := &tasks[taskID-1]
	
	if targetTask.Status == task.StatusCompleted {
		warningColor := color.New(color.FgYellow)
		warningColor.Printf("⚠️  Task #%d is already completed!\n", taskID)
		return nil
	}
	
	// Mark as completed
	targetTask.Complete()
	
	// Save updated tasks
	if err := storage.SaveTasks(tasks); err != nil {
		return fmt.Errorf("failed to save task completion: %w", err)
	}
	
	successColor := color.New(color.FgGreen, color.Bold)
	successColor.Printf("🎉 Task #%d completed successfully!\n", taskID)
	fmt.Printf("   Description: %s\n", targetTask.Description)
	
	return nil
}
```

### 8.2 Test All Commands Together
```bash
# Add some tasks
go run . add "Learn Go basics" --priority high
go run . add "Build CLI tool" --priority medium
go run . add "Write tests" --priority low

# List tasks
go run . list

# Complete a task
go run . complete 1

# List again to see the change
go run . list --all
```

**🎯 What you learned:**
- String to integer conversion
- Array/slice bounds checking
- Method calls on struct pointers
- State modification and persistence

## Step 9: Testing Your CLI (20 minutes)

### 9.1 Create `tests/task_test.go`
```go
package tests

import (
	"testing"
	"time"
	"task-manager-cli/internal/task"
)

func TestTaskComplete(t *testing.T) {
	testTask := task.Task{
		Description: "Test complete method",
		Status:      task.StatusPending,
	}
	
	// Test complete method
	testTask.Complete()
	
	if testTask.Status != task.StatusCompleted {
		t.Error("Task status should be completed")
	}
	
	if testTask.CompletedAt == nil {
		t.Error("CompletedAt should be set")
	}
}

func TestTaskIsOverdue(t *testing.T) {
	now := time.Now()
	yesterday := now.Add(-24 * time.Hour)
	
	testTask := task.Task{
		Description: "Overdue task",
		Status:      task.StatusPending,
		DueDate:     &yesterday,
	}
	
	if !testTask.IsOverdue() {
		t.Error("Task with past due date should be overdue")
	}
}
```

### 9.2 Run Tests
```bash
go test ./tests/
go test -v ./tests/  # verbose output
```

**🎯 What you learned:**
- Unit testing in Go
- Table-driven tests
- Test assertions
- Testing time-dependent code

## Step 10: Building and Distribution (15 minutes)

### 10.1 Create a Makefile
```makefile
BINARY_NAME=task-manager

build:
	go build -o $(BINARY_NAME) .

test:
	go test ./...

clean:
	rm -f $(BINARY_NAME)

install:
	go install .

# Cross-platform builds
build-all:
	GOOS=linux GOARCH=amd64 go build -o bin/$(BINARY_NAME)-linux-amd64
	GOOS=darwin GOARCH=amd64 go build -o bin/$(BINARY_NAME)-darwin-amd64
	GOOS=windows GOARCH=amd64 go build -o bin/$(BINARY_NAME)-windows-amd64.exe

.PHONY: build test clean install build-all
```

### 10.2 Build and Install
```bash
# Build for current platform
make build

# Install globally
make install

# Cross-platform builds
make build-all
```

**🎯 What you learned:**
- Build automation with Make
- Cross-platform compilation
- Go install process
- Binary distribution

## 🎯 Congratulations! You've Built a Complete CLI Tool!

### What You've Accomplished:
✅ **Project Structure**: Organized Go project with proper package structure
✅ **CLI Framework**: Used Cobra for professional command-line interface
✅ **Data Management**: JSON-based persistence with proper error handling
✅ **User Experience**: Interactive prompts and colored output
✅ **Testing**: Unit tests for core functionality
✅ **Distribution**: Cross-platform builds and installation

### Key Go Concepts You've Mastered:
- Package organization and imports
- Struct definitions and methods
- Interface design and implementation
- Error handling and wrapping
- JSON marshaling/unmarshaling
- File I/O operations
- String manipulation
- Time handling
- Testing patterns

### Next Challenges:
1. **Add more features**: Search, tags, import/export
2. **Improve UX**: Better error messages, progress bars
3. **Add configuration**: YAML config files, environment variables
4. **Enhance testing**: Integration tests, mocking
5. **Create plugins**: Extensible architecture

## 🚀 Ready for More?

Now that you have a solid foundation, explore:
- **TUI Development**: Build terminal user interfaces with Bubble Tea
- **API Integration**: Connect to external services
- **Database Integration**: Use SQLite or PostgreSQL
- **Plugin Architecture**: Make your CLI extensible
- **CI/CD**: Automate builds and releases with GitHub Actions

Happy coding! 🎉