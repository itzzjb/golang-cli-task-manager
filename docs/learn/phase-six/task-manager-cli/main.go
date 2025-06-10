package main

import (
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/fatih/color"
)

func main() {

	// Single input
	var name string
	promptName := survey.Input{
		Message: "What is your name?",
	}
	errName := survey.AskOne(&promptName, &name)
	if errName != nil {
		panic(errName)
	}

	// Direct Color Printing
	msg := fmt.Sprintf("Hi %s, Welcome to CLI Color Tool!", name)
	fmt.Println(color.RedString(msg))

	// Create selectedColor functions
	red := color.New(color.FgRed).SprintFunc()
	green := color.New(color.FgGreen).SprintFunc()
	blue := color.New(color.FgBlue).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()

	// Multiple choice
	var selectedColor string
	promptColor := survey.Select{
		Message: "What is your favorite selectedColor?",
		Options: []string{"Red", "Green", "Blue", "Yellow"},
	}
	errColor := survey.AskOne(&promptColor, &selectedColor)
	if errColor != nil {
		panic(errColor)
	}

	switch selectedColor {
	case "Red":
		fmt.Println("You picked", red("Red"))
	case "Green":
		fmt.Println("You picked", green("Green"))
	case "Blue":
		fmt.Println("You picked", blue("Blue"))
	case "Yellow":
		fmt.Println("You picked", yellow("Yellow"))
	default:
		fmt.Println("You picked an unknown selectedColor")
	}

	// Confirmation
	var confirm bool
	promptConfirm := survey.Confirm{
		Message: "Are you sure?",
	}
	errConfirm := survey.AskOne(&promptConfirm, &confirm)
	if errConfirm != nil {
		panic(errConfirm)
	}

}
