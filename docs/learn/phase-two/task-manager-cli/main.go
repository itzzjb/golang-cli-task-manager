package main

import (
	"flag"
	"fmt"
)

func main() {

	var name = flag.String("name", "Default Name", "name of the task")
	var verbose = flag.Bool("verbose", false, "enable verbose output")
	var count = flag.Int("count", 1, "number of times to greet")
	var timeout = flag.Duration("timeout", 0, "timeout duration")

	flag.Parse()

	if *timeout > 0 {
		fmt.Println("Timeout set to", *timeout)
	} else {
		fmt.Println("No timeout set")
	}

	if *count < 1 {
		fmt.Println("Count must be at least 1")
		return
	} else if *count > 1 {
		for i := 0; i < *count; i++ {
			if *verbose {
				fmt.Printf("Greeting %d: Hello, %s!\n", i+1, *name)
			} else {
				fmt.Printf("Hello, %s!\n", *name)
			}
		}
		return
	} else {
		if *verbose {
			fmt.Println("Greeting Hello once to", *name, "verbosely")
		} else {
			fmt.Println("Greeting Hello once to", *name)
		}
	}
}

// go build -o task-manager-cli main.go
// ./task-manager-cli
// ./task-manager-cli -name "John Doe"

// ./task-manager-cli -name "John Doe" -verbose true ( No need to specify true for boolean flags, just use the flag name )
// ./task-manager-cli -name "John Doe" -verbose
// ./task-manager-cli -verbose
