# Go CLI Development - Complete Learning Guide

## 🎯 Welcome to Your CLI Learning Journey!

This comprehensive guide will walk you through creating a professional CLI application in Go from scratch. You'll learn every essential concept through practical examples.

## 📚 Learning Path

### Phase 1: Understanding the Foundation (Day 1-2)

#### 1.1 Project Setup
```bash
# Initialize your project
mkdir task-manager-cli
cd task-manager-cli
go mod init task-manager-cli

# Install dependencies
go get github.com/spf13/cobra@latest
go get github.com/spf13/viper@latest
go get github.com/AlecAivazis/survey/v2@latest
go get github.com/fatih/color@latest
```

#### 1.2 Understanding Go Modules
- `go.mod`: Defines your module and dependencies
- `go.sum`: Cryptographic checksums for dependencies
- `go mod tidy`: Cleans up dependencies

#### 1.3 Basic CLI Structure
Start with the simplest possible CLI:

```go
package main

import (
    "fmt"
    "os"
)

func main() {
    if len(os.Args) < 2 {
        fmt.Println("Usage: task-manager <command>")
        return
    }
    
    command := os.Args[1]
    switch command {
    case "hello":
        fmt.Println("Hello, World!")
    default:
        fmt.Printf("Unknown command: %s\n", command)
    }
}
```

### Phase 2: Working with Flags (Day 3-4)

#### 2.1 Standard Flag Package
Learn Go's built-in flag package first:

```go
package main

import (
    "flag"
    "fmt"
)

func main() {
    var name = flag.String("name", "World", "name to greet")
    var verbose = flag.Bool("verbose", false, "verbose output")
    flag.Parse()

    if *verbose {
        fmt.Printf("Greeting %s verbosely!\n", *name)
    }
    fmt.Printf("Hello, %s!\n", *name)
}
```

#### 2.2 Understanding Pointers in Flags
- Flag functions return pointers to values
- Use `*` to dereference the pointer
- Understand why flags use pointers (memory efficiency)

#### 2.3 Different Flag Types
```go
// String flags
name := flag.String("name", "default", "description")

// Integer flags
count := flag.Int("count", 1, "number of items")

// Boolean flags
verbose := flag.Bool("verbose", false, "enable verbose output")

// Duration flags
timeout := flag.Duration("timeout", 30*time.Second, "timeout duration")
```

### Phase 3: Introduction to Cobra (Day 5-7)

#### 3.1 Why Cobra?
- Used by Kubernetes, Docker, GitHub CLI
- Automatic help generation
- Subcommands support
- Flag inheritance
- Shell completion

#### 3.2 Basic Cobra Setup
```go
package main

import (
    "fmt"
    "github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
    Use:   "myapp",
    Short: "A brief description",
    Long:  "A longer description",
    Run: func(cmd *cobra.Command, args []string) {
        fmt.Println("Hello from Cobra!")
    },
}

func main() {
    rootCmd.Execute()
}
```

#### 3.3 Adding Subcommands
```go
var addCmd = &cobra.Command{
    Use:   "add [item]",
    Short: "Add a new item",
    Args:  cobra.ExactArgs(1),
    Run: func(cmd *cobra.Command, args []string) {
        fmt.Printf("Adding: %s\n", args[0])
    },
}

func init() {
    rootCmd.AddCommand(addCmd)
}
```

### Phase 4: Configuration Management (Day 8-9)

#### 4.1 Understanding Viper
- Configuration from files (YAML, JSON, TOML)
- Environment variables
- Command line flags
- Default values
- Configuration merging

#### 4.2 Basic Viper Setup
```go
import (
    "github.com/spf13/viper"
)

func initConfig() {
    viper.SetConfigName("config")
    viper.SetConfigType("yaml")
    viper.AddConfigPath(".")
    viper.AddConfigPath("$HOME/.myapp")
    
    viper.AutomaticEnv()
    
    if err := viper.ReadInConfig(); err != nil {
        // Handle error
    }
}
```

### Phase 5: Data Persistence (Day 10-11)

#### 5.1 JSON File Storage
```go
type Task struct {
    ID          int       `json:"id"`
    Description string    `json:"description"`
    Completed   bool      `json:"completed"`
    CreatedAt   time.Time `json:"created_at"`
}

func saveTasks(tasks []Task, filename string) error {
    data, err := json.MarshalIndent(tasks, "", "  ")
    if err != nil {
        return err
    }
    return os.WriteFile(filename, data, 0644)
}

func loadTasks(filename string) ([]Task, error) {
    data, err := os.ReadFile(filename)
    if err != nil {
        return nil, err
    }
    
    var tasks []Task
    err = json.Unmarshal(data, &tasks)
    return tasks, err
}
```

#### 5.2 Error Handling Best Practices
```go
func loadTasks(filename string) ([]Task, error) {
    // Check if file exists
    if _, err := os.Stat(filename); os.IsNotExist(err) {
        return []Task{}, nil // Return empty slice, not error
    }
    
    data, err := os.ReadFile(filename)
    if err != nil {
        return nil, fmt.Errorf("failed to read file %s: %w", filename, err)
    }
    
    var tasks []Task
    if err := json.Unmarshal(data, &tasks); err != nil {
        return nil, fmt.Errorf("failed to parse JSON: %w", err)
    }
    
    return tasks, nil
}
```

### Phase 6: User Interaction (Day 12-13)

#### 6.1 Interactive Prompts with Survey
```go
import "github.com/AlecAivazis/survey/v2"

// Simple input
var name string
prompt := &survey.Input{
    Message: "What's your name?",
}
survey.AskOne(prompt, &name)

// Multiple choice
var color string
prompt := &survey.Select{
    Message: "Choose a color:",
    Options: []string{"red", "blue", "green"},
}
survey.AskOne(prompt, &color)

// Confirmation
var confirmed bool
prompt := &survey.Confirm{
    Message: "Are you sure?",
}
survey.AskOne(prompt, &confirmed)
```

#### 6.2 Colored Output
```go
import "github.com/fatih/color"

// Create color functions
red := color.New(color.FgRed).SprintFunc()
green := color.New(color.FgGreen, color.Bold).SprintFunc()

fmt.Printf("This is %s and this is %s\n", red("red"), green("bold green"))

// Direct colored printing
color.Red("This is red text")
color.Green("This is green text")
```

### Phase 7: Testing CLI Applications (Day 14-15)

#### 7.1 Unit Testing
```go
func TestTaskComplete(t *testing.T) {
    task := Task{
        Description: "Test task",
        Completed:   false,
    }
    
    task.Complete()
    
    if !task.Completed {
        t.Error("Task should be completed")
    }
    
    if task.CompletedAt.IsZero() {
        t.Error("CompletedAt should be set")
    }
}
```

#### 7.2 Integration Testing
```go
func TestCLIAdd(t *testing.T) {
    // Build CLI binary
    cmd := exec.Command("go", "build", "-o", "test-cli", ".")
    err := cmd.Run()
    if err != nil {
        t.Fatal("Failed to build CLI")
    }
    defer os.Remove("test-cli")
    
    // Run CLI command
    cmd = exec.Command("./test-cli", "add", "test task")
    output, err := cmd.Output()
    if err != nil {
        t.Fatal("CLI command failed")
    }
    
    // Verify output
    if !strings.Contains(string(output), "Task added") {
        t.Error("Expected success message")
    }
}
```

### Phase 8: Advanced Features (Day 16-18)

#### 8.1 Shell Completion
```go
// Add to your root command
var completionCmd = &cobra.Command{
    Use:   "completion [bash|zsh|fish|powershell]",
    Short: "Generate completion script",
    Run: func(cmd *cobra.Command, args []string) {
        switch args[0] {
        case "bash":
            cmd.Root().GenBashCompletion(os.Stdout)
        case "zsh":
            cmd.Root().GenZshCompletion(os.Stdout)
        // ... other shells
        }
    },
}
```

#### 8.2 Plugin Architecture
```go
type Plugin interface {
    Name() string
    Execute(args []string) error
}

func loadPlugins(pluginDir string) []Plugin {
    // Load .so files or execute plugin binaries
    // Return slice of plugins
}
```

#### 8.3 Configuration Validation
```go
func validateConfig(config Config) error {
    if config.Storage.Path == "" {
        return errors.New("storage path cannot be empty")
    }
    
    validPriorities := []string{"low", "medium", "high"}
    if !contains(validPriorities, config.Defaults.Priority) {
        return fmt.Errorf("invalid priority: %s", config.Defaults.Priority)
    }
    
    return nil
}
```

### Phase 9: Building and Distribution (Day 19-20)

#### 9.1 Cross-Platform Builds
```bash
# Build for different platforms
GOOS=linux GOARCH=amd64 go build -o task-manager-linux-amd64
GOOS=windows GOARCH=amd64 go build -o task-manager-windows-amd64.exe
GOOS=darwin GOARCH=amd64 go build -o task-manager-darwin-amd64
```

#### 9.2 Makefile for Automation
```makefile
.PHONY: build test clean

build:
	go build -o task-manager .

test:
	go test ./...

clean:
	rm -f task-manager

install:
	go install .

release:
	GOOS=linux GOARCH=amd64 go build -o dist/task-manager-linux-amd64
	GOOS=windows GOARCH=amd64 go build -o dist/task-manager-windows-amd64.exe
	GOOS=darwin GOARCH=amd64 go build -o dist/task-manager-darwin-amd64
```

#### 9.3 GitHub Actions for CI/CD
```yaml
name: Build and Test

on: [push, pull_request]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v2
    - uses: actions/setup-go@v2
      with:
        go-version: 1.21
    
    - name: Test
      run: go test ./...
    
    - name: Build
      run: go build .
```

## 🎯 Key Learning Objectives

### By the end of this guide, you'll understand:

1. **CLI Architecture**
   - Command structure and hierarchy
   - Flag and argument handling
   - Error handling and user feedback

2. **Go-Specific Concepts**
   - Package organization
   - Interface design
   - Error wrapping and handling
   - JSON marshaling/unmarshaling
   - File I/O operations

3. **Professional Development**
   - Testing strategies (unit and integration)
   - Configuration management
   - Build automation
   - Cross-platform compatibility

4. **User Experience**
   - Interactive prompts
   - Colored output
   - Help generation
   - Shell completion

## 🔧 Development Workflow

### Daily Practice Routine:
1. **Morning** (30 min): Read concepts and theory
2. **Afternoon** (1-2 hours): Hands-on coding
3. **Evening** (30 min): Testing and refinement

### Checkpoint Exercises:
- **Day 5**: Build a simple calculator CLI
- **Day 10**: Create a file organizer tool
- **Day 15**: Build a note-taking application
- **Day 20**: Complete task manager project

## 🐛 Common Pitfalls and Solutions

### 1. Pointer Confusion with Flags
**Problem**: Forgetting to dereference flag pointers
```go
// Wrong
fmt.Println(name) // prints address

// Correct
fmt.Println(*name) // prints value
```

### 2. Error Handling
**Problem**: Ignoring errors or poor error messages
```go
// Poor
if err != nil {
    panic(err)
}

// Better
if err != nil {
    return fmt.Errorf("failed to process task: %w", err)
}
```

### 3. JSON Tag Confusion
**Problem**: Inconsistent JSON field naming
```go
// Inconsistent
type Task struct {
    ID          int    `json:"id"`
    Description string `json:"desc"`  // Should be "description"
    CreatedAt   time.Time `json:"created_at"`
}
```

### 4. File Path Issues
**Problem**: Hardcoded paths that break on different systems
```go
// Wrong
path := "/home/user/tasks.json"

// Better
home, _ := os.UserHomeDir()
path := filepath.Join(home, ".myapp", "tasks.json")
```

## 🎉 Next Steps

After completing this project:
1. **Add Features**: Tags, search, import/export
2. **Integrate APIs**: GitHub issues, Trello, etc.
3. **Build Plugins**: Extend functionality
4. **Create TUI**: Terminal User Interface with Bubble Tea
5. **Build Web Interface**: Complement your CLI

## 📖 Additional Resources

- [Cobra Documentation](https://cobra.dev/)
- [Viper Documentation](https://github.com/spf13/viper)
- [Go CLI Patterns](https://pkg.go.dev/flag)
- [Effective Go](https://golang.org/doc/effective_go.html)
- [Go Testing](https://golang.org/pkg/testing/)

Happy coding! 🚀