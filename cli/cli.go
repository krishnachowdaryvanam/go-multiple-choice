// cli.go

package cli

import (
	"fmt"
	"multiple-choice-test-engine/models"
	"multiple-choice-test-engine/validator"
	"strings"
	"time"
)

// getUserInput reads user input from the command line.
func GetUserInput() string {
	var input string
	fmt.Scanln(&input)
	return strings.TrimSpace(input)
}

// RunCLI runs the command-line interface for the test engine.
func RunCLI(project *models.Project) {
	fmt.Println("Welcome to the Multiple Choice Test!")

	// Load the test from a JSON file
	err := project.LoadTestFromJSON("test.json")
	if err != nil {
		fmt.Printf("Error loading test: %v\n", err)
		return
	}

	// Initialize session to store user's answers and test start time
	session := make(map[string]string)
	startTime := time.Now()

	// Shuffle questions (optional)
	project.ShuffleQuestions()

	// Set a timer for the test duration (optional, set to 300 seconds here)
	timerDone := project.SetTimer(300)

	fmt.Println("Answer the following questions:")

	// Display questions one by one
	for i, question := range project.Test.Questions {
		fmt.Printf("Question %d: %s\n", i+1, question.Text)
		for j, option := range question.Options {
			fmt.Printf("%c. %s\n", 'A'+j, option)
		}

		// Prompt for user's answer
		fmt.Print("Your answer: ")
		answer := GetUserInput()

		// Validate user input
		for !validator.IsValidOption(answer) {
			fmt.Println("Invalid answer option. Please enter a valid option.")
			fmt.Print("Your answer: ")
			answer = GetUserInput()
		}

		// Store the user's answer in the session
		session[question.Text] = answer

		// Notify the user when the test is over (optional)
		select {
		case <-timerDone:
			fmt.Println("\nTest completed due to time limit.")
			displayFinalScore(project, session, startTime)
			return
		default:
		}
	}

	// Display the final score
	displayFinalScore(project, session, startTime)
}

// displayFinalScore displays the final score at the end of the test.
func displayFinalScore(project *models.Project, session map[string]string, startTime time.Time) {
	// Extract user answers from the session
	userAnswers := make([]string, len(project.Test.Questions))
	for i, question := range project.Test.Questions {
		userAnswers[i] = session[question.Text]
	}

	// Calculate and display the final score
	score := project.CalculateScore(userAnswers)
	totalQuestions := len(project.Test.Questions)
	fmt.Printf("\nYou answered %d out of %d questions correctly.\n", score, totalQuestions)

	// Calculate and display the time taken for the test
	duration := time.Since(startTime).Round(time.Second)
	fmt.Printf("Time taken: %s\n", duration)
}
