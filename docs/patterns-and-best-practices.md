# Common Go CLI Patterns and Best Practices

This guide covers common patterns and best practices for building command-line interfaces (CLIs) in Go. Use these tips to build professional, maintainable CLI applications.

## Command Structure Patterns

### 1. Git-style Subcommands

```go
rootCmd
├── add
│   ├── file
│   └── directory
├── remove
└── list
```

**Example:**
```go
// Root command
var rootCmd = &cobra.Command{
    Use:   "mycli",
    Short: "A sample CLI",
}

// Add command
var addCmd = &cobra.Command{
    Use:   "add",
    Short: "Add resources",
}

// Add subcommands
var addFileCmd = &cobra.Command{
    Use:   "file [filename]",
    Short: "Add a file",
    Run: func(cmd *cobra.Command, args []string) {
        // Implementation
    },
}

func init() {
    rootCmd.AddCommand(addCmd)
    addCmd.AddCommand(addFileCmd)
}
```

### 2. Flat Command Structure

```go
rootCmd
├── add-file
├── remove-file
├── list-files
└── version
```

**Example:**
```go
var rootCmd = &cobra.Command{
    Use:   "mycli",
    Short: "A sample CLI",
}

var addFileCmd = &cobra.Command{
    Use:   "add-file [filename]",
    Short: "Add a file",
    Run: func(cmd *cobra.Command, args []string) {
        // Implementation
    },
}

func init() {
    rootCmd.AddCommand(addFileCmd)
}
```

## Flag Patterns

### 1. Global Flags vs. Command Flags

```go
// Global flags (available to all commands)
rootCmd.PersistentFlags().StringP("output", "o", "", "output format")
rootCmd.PersistentFlags().BoolP("verbose", "v", false, "verbose output")

// Command-specific flags
addCmd.Flags().StringP("type", "t", "default", "resource type")
```

### 2. Flag Sets

```go
// Reusable flag set for database commands
func addDatabaseFlags(cmd *cobra.Command) {
    cmd.Flags().StringP("host", "", "localhost", "database host")
    cmd.Flags().IntP("port", "", 5432, "database port")
    cmd.Flags().StringP("user", "u", "", "database user")
    cmd.Flags().StringP("password", "p", "", "database password")
}

// Apply to multiple commands
addDatabaseFlags(createCmd)
addDatabaseFlags(migrateCmd)
```

### 3. Required Flags

```go
cmd.Flags().StringP("filename", "f", "", "input filename")
cmd.MarkFlagRequired("filename")
```

## Configuration Patterns

### 1. Configuration Hierarchy

1. Command-line flags (highest priority)
2. Environment variables
3. Configuration file
4. Default values (lowest priority)

```go
// Example with viper
func initConfig() {
    // 1. Load defaults
    viper.SetDefault("port", 8080)
    
    // 2. Load from file
    viper.SetConfigName("config")
    viper.AddConfigPath(".")
    viper.ReadInConfig()
    
    // 3. Load from environment
    viper.AutomaticEnv()
    
    // 4. Command-line flags override everything
    // (handled by viper.BindPFlag in flag setup)
}
```

### 2. Environment Variables Mapping

```go
// Automatically map environment variables
viper.SetEnvPrefix("MYAPP")
viper.AutomaticEnv()

// Custom mappings
viper.BindEnv("database.host", "MYAPP_DB_HOST")
```

## Error Handling Patterns

### 1. Consistent Error Messaging

```go
// Helper function for consistent error formatting
func formatError(operation string, err error) error {
    return fmt.Errorf("failed to %s: %w", operation, err)
}

// Usage
if err := doSomething(); err != nil {
    return formatError("process file", err)
}
```

### 2. Exit Codes

```go
const (
    ExitSuccess        = 0
    ExitError          = 1
    ExitConfigError    = 2
    ExitPermissionError = 3
)

func handleError(err error) int {
    if err == nil {
        return ExitSuccess
    }
    
    fmt.Fprintf(os.Stderr, "Error: %s\n", err)
    
    switch {
    case errors.Is(err, ErrConfigNotFound):
        return ExitConfigError
    case errors.Is(err, os.ErrPermission):
        return ExitPermissionError
    default:
        return ExitError
    }
}

// In main function
os.Exit(handleError(cmd.Execute()))
```

## Input/Output Patterns

### 1. Table Output

```go
func displayTable(items []Item) {
    w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
    
    fmt.Fprintln(w, "ID\tName\tCreated\tStatus")
    fmt.Fprintln(w, "--\t----\t-------\t------")
    
    for _, item := range items {
        fmt.Fprintf(w, "%d\t%s\t%s\t%s\n", 
            item.ID, 
            item.Name, 
            item.CreatedAt.Format("2006-01-02"),
            item.Status)
    }
    
    w.Flush()
}
```

### 2. Progressive Output

```go
func longRunningTask() {
    steps := 10
    
    for i := 1; i <= steps; i++ {
        fmt.Printf("\rProcessing: [")
        for j := 1; j <= steps; j++ {
            if j <= i {
                fmt.Print("=")
            } else {
                fmt.Print(" ")
            }
        }
        fmt.Printf("] %d%%", i*10)
        
        // Do work
        time.Sleep(500 * time.Millisecond)
    }
    fmt.Println("\rCompleted!          ")
}
```

### 3. Output Formats

```go
var outputFormat string

func init() {
    rootCmd.PersistentFlags().StringVarP(&outputFormat, "output", "o", "text", "Output format (text|json|yaml)")
}

func displayResult(result interface{}) error {
    switch outputFormat {
    case "json":
        data, err := json.MarshalIndent(result, "", "  ")
        if err != nil {
            return err
        }
        fmt.Println(string(data))
        
    case "yaml":
        data, err := yaml.Marshal(result)
        if err != nil {
            return err
        }
        fmt.Println(string(data))
        
    default: // text
        // Human-readable format
        fmt.Printf("Name: %s\nCount: %d\n", result.Name, result.Count)
    }
    
    return nil
}
```

## Testing Patterns

### 1. Command Testing

```go
func TestAddCommand(t *testing.T) {
    // Redirect stdout
    oldStdout := os.Stdout
    r, w, _ := os.Pipe()
    os.Stdout = w
    
    // Run command
    rootCmd.SetArgs([]string{"add", "test task"})
    rootCmd.Execute()
    
    // Restore stdout
    w.Close()
    os.Stdout = oldStdout
    
    // Read captured output
    var buf bytes.Buffer
    io.Copy(&buf, r)
    output := buf.String()
    
    // Assertions
    if !strings.Contains(output, "Task added successfully") {
        t.Errorf("Expected success message, got: %s", output)
    }
}
```

### 2. Integration Testing

```go
func TestCLIEndToEnd(t *testing.T) {
    // Build test binary
    cmd := exec.Command("go", "build", "-o", "test-cli")
    err := cmd.Run()
    if err != nil {
        t.Fatal(err)
    }
    defer os.Remove("test-cli")
    
    // Run CLI commands
    addCmd := exec.Command("./test-cli", "add", "Test task")
    addOutput, err := addCmd.CombinedOutput()
    if err != nil {
        t.Fatalf("Add command failed: %v\n%s", err, addOutput)
    }
    
    // Check output
    if !strings.Contains(string(addOutput), "Task added successfully") {
        t.Errorf("Expected success message, got: %s", string(addOutput))
    }
    
    // Check data persistence
    listCmd := exec.Command("./test-cli", "list")
    listOutput, err := listCmd.CombinedOutput()
    if err != nil {
        t.Fatalf("List command failed: %v\n%s", err, listOutput)
    }
    
    if !strings.Contains(string(listOutput), "Test task") {
        t.Errorf("Task not found in list output: %s", string(listOutput))
    }
}
```

## User Experience Patterns

### 1. Progress Indication

```go
func processFiles(files []string) {
    fmt.Printf("Processing %d files...\n", len(files))
    
    for i, file := range files {
        fmt.Printf("[%d/%d] Processing %s\n", i+1, len(files), file)
        // Process file
        time.Sleep(500 * time.Millisecond)
        fmt.Println("✓ Done")
    }
    
    fmt.Println("All files processed successfully!")
}
```

### 2. Confirmation Prompts

```go
func confirmDeletion() (bool, error) {
    confirm := false
    prompt := &survey.Confirm{
        Message: "Are you sure you want to delete this item?",
        Default: false,
    }
    
    err := survey.AskOne(prompt, &confirm)
    return confirm, err
}

func deleteItem(id string) error {
    confirmed, err := confirmDeletion()
    if err != nil {
        return err
    }
    
    if !confirmed {
        fmt.Println("Operation cancelled")
        return nil
    }
    
    // Proceed with deletion
    fmt.Printf("Deleting item %s...\n", id)
    return nil
}
```

### 3. Consistent Color Scheme

```go
// Define standard colors
var (
    success = color.New(color.FgGreen, color.Bold)
    warning = color.New(color.FgYellow)
    error   = color.New(color.FgRed, color.Bold)
    info    = color.New(color.FgCyan)
    debug   = color.New(color.FgBlue)
)

// Use consistently
success.Println("✅ Operation successful!")
warning.Println("⚠️  Some items could not be processed")
error.Println("❌ Failed to complete operation")
info.Println("ℹ️  Processing items...")
debug.Println("🔍 Debug info: connection established")
```

## Distributing CLI Applications

### 1. Cross-Platform Builds

```bash
# Build for all target platforms
GOOS=linux GOARCH=amd64 go build -o myapp-linux-amd64
GOOS=darwin GOARCH=amd64 go build -o myapp-darwin-amd64
GOOS=windows GOARCH=amd64 go build -o myapp-windows-amd64.exe
```

### 2. Version Information

```go
var (
    version   = "dev"
    commit    = "none"
    buildTime = "unknown"
)

var versionCmd = &cobra.Command{
    Use:   "version",
    Short: "Print version information",
    Run: func(cmd *cobra.Command, args []string) {
        fmt.Printf("Version: %s\n", version)
        fmt.Printf("Commit: %s\n", commit)
        fmt.Printf("Build Time: %s\n", buildTime)
    },
}

// During build:
// go build -ldflags "-X main.version=1.0.0 -X main.commit=abc123 -X main.buildTime=2023-01-01"
```

## 🎯 Checklist for Professional CLIs

✅ **Well-defined command structure**
✅ **Consistent help text and examples**
✅ **Proper error handling with meaningful messages**
✅ **Configuration from multiple sources**
✅ **Interactive mode when needed**
✅ **Appropriate progress indicators**
✅ **Consistent color scheme**
✅ **Comprehensive testing**
✅ **Version information**
✅ **Cross-platform support**
✅ **Documentation**

By following these patterns and best practices, you'll create professional, user-friendly CLI applications that stand out for their quality and usability.