package main

import (
	"flag"
	"fmt"
)

func main() {
	var name = flag.String("name", "Default Name", "name of the task")
	var verbose = flag.Bool("verbose", false, "enable verbose output")
	flag.Parse()
	if *verbose {
		fmt.Println("Greeting to", *name, "verbosely")
	} else {
		fmt.Println("Greeting to", *name)
	}
}

// go build -o task-manager-cli main.go
// ./task-manager-cli
// ./task-manager-cli -name "John Doe"

// ./task-manager-cli -name "John Doe" -verbose true ( No need to specify true for boolean flags, just use the flag name )
// ./task-manager-cli -name "John Doe" -verbose
// ./task-manager-cli -verbose
