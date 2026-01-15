package ui

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func AskUser(prompt string) string {
	fmt.Print(prompt)
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}

func AskUserInt(prompt string) (int, error) {
	input := AskUser(prompt)
	return strconv.Atoi(input)
}
func AskUserString(prompt string) string {
	input := AskUser(prompt)
	return strings.TrimSpace(input)
}

// func ConfirmAction(message string, args ...interface{}) bool {
// 	fmt.Printf("\n"+message+"\n", args...)

// 	reader := bufio.NewReader(os.Stdin)
// 	fmt.Print("Proceed? (y/N): ")
// 	input, _ := reader.ReadString('\n')
// 	input = strings.TrimSpace(input)

// 	return strings.ToLower(input) == "y"
// }

// func AskUserConfirmation(msgs ...string) bool {
// 	var (
// 		question string = "Are you sure you want to continue? (y/n): " // default question
// 		response string
// 	)
// 	if len(msgs) > 0 {
// 		question = fmt.Sprintf("%s - confirmation required (y/n): ", msgs[0])
// 	}
// 	fmt.Print(question)
// 	fmt.Scanln(&response)
// 	return response == "y" || response == "Y"
// }

// Name: AskUserConfirmation
//
// Description: prompts the user for confirmation.
//
// Usage examples:
//
//	ui.AskUserConfirmation()
//	ui.AskUserConfirmation("Delete all resources")
//	ui.AskUserConfirmation("Install %s/%s (%s) into namespace %s as release %s",
//	    helmRepo.Name, helmChart.FullName, helmChart.Version, k8sNsName, helmReleaseName)
//
// if !ui.AskUserConfirmation(
//
//	"You are about to install %s/%s (%s) into namespace %s as release %s",
//	helmRepo.Name, helmChart.FullName, helmChart.Version, k8sNsName, helmReleaseName,
//
//	) {
//		logger.Info("aborted by user")
//		return
//	}
//
//	if !ui.AskUserConfirmation("Delete all Helm releases in namespace %s", k8sNsName) {
//		logger.Info("operation cancelled")
//		return
//	}
//
// if !ui.AskUserConfirmation("") { // falls back to default message
//
//		return
//	}
//
// Returns true if the user confirms with "y" or "Y".
func AskUserConfirmation(format string, args ...interface{}) bool {
	var question string

	// If a formatted message is provided, build the custom question
	if format != "" {
		msg := fmt.Sprintf(format, args...)
		question = fmt.Sprintf("\n%s - confirmation required (y/N): ", msg)
	} else {
		// Default prompt
		question = "Are you sure you want to continue? (y/N): "
	}

	// Prompt user
	fmt.Print(question)
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	return strings.ToLower(input) == "y"
}
