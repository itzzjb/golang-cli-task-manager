package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: task-manager-cli <command>")
		return
	}
	command := os.Args[1]
	switch command {
	case "hello":
		fmt.Println("Hello, Task Manager CLI!")
	default:
		fmt.Printf("Unknown command: %s\n", command)
	}
}

// go build -o task-manager-cli main.go
// ./task-manager-cli hello
