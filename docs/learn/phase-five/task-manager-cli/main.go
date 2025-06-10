package main

import (
	"encoding/json"
	"os"
	"time"
)

type Task struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Completed   bool      `json:"completed"`
	CreatedAt   time.Time `json:"created_at"`
}

// MarshalIndent - JSON string in a pretty-printed (indented) format
func saveTasks(tasks []Task, filename string) error {
	// prefix "" - The prefix for each line in the output (here, empty, so no prefix).
	// indent "  " - The indentation to use for nested values (here, two spaces).
	data, err := json.MarshalIndent(tasks, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(filename, data, 0644) // 0644 - The file permission bits (readable and writable by the owner, readable by groups and others).
}

func loadTasks(filename string) ([]Task, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var tasks []Task
	err = json.Unmarshal(data, &tasks)
	if err != nil {
		return nil, err
	}

	return tasks, nil
}

func main() {

	taskOne := Task{
		ID:          1,
		Title:       "Learn Go",
		Description: "Complete the Go programming course",
		Completed:   false,
	}
	taskTwo := Task{
		ID:          2,
		Title:       "Build a CLI",
		Description: "Create a command-line interface for task management",
		Completed:   false,
	}

	tasks := []Task{taskOne, taskTwo}

	err := saveTasks(tasks, "tasks.json")
	if err != nil {
		panic(err)
	}

	loadedTasks, err := loadTasks("tasks.json")
	if err != nil {
		panic(err)
	}
	
	for _, task := range loadedTasks {
		println("Task ID:", task.ID)
		println("Title:", task.Title)
		println("Description:", task.Description)
		println("Completed:", task.Completed)
		println("Created At:", task.CreatedAt.String())
		println()
	}

}
