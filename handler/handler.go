// handler/handler.go
package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"multiple-choice-test-engine/models"
	"multiple-choice-test-engine/validator"

	"github.com/gin-gonic/gin"
)

func StartTestHandler(c *gin.Context) {
	project, exists := c.Get("project")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Project not found"})
		return
	}

	if project.(*models.Project).Test == nil || len(project.(*models.Project).Test.Questions) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No questions loaded"})
		return
	}

	// Initialize session to store user's answers and test start time
	session := make(map[string]string)
	c.Set("session", session)
	startTime := time.Now()
	c.Set("startTime", startTime)

	// Shuffle questions (optional)
	project.(*models.Project).ShuffleQuestions()

	// Set a timer for the test duration (optional, set to 300 seconds here)
	timerDone := project.(*models.Project).SetTimer(300)

	// Display the first question
	question := project.(*models.Project).Test.Questions[0]
	c.JSON(http.StatusOK, gin.H{"question": project.(*models.Project).DisplayQuestion(question)})

	// Notify the client when the test is over (optional)
	go func() {
		<-timerDone
		c.SSEvent("testCompleted", "Test completed due to time limit.")
	}()
}

func SubmitAnswerHandler(c *gin.Context) {
	session, exists := c.Get("session")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Session not found"})
		return
	}

	var requestBody map[string]string
	if err := c.BindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON format"})
		return
	}

	// Validate user input
	if !validator.IsValidOption(requestBody["answer"]) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid answer option"})
		return
	}

	// Store the user's answer in the session
	session.(map[string]string)[requestBody["question"]] = requestBody["answer"]

	// Move to the next question
	currentQuestionIndex, _ := strconv.Atoi(requestBody["questionIndex"])
	project := c.MustGet("project").(*models.Project)
	if currentQuestionIndex < len(project.Test.Questions)-1 {
		nextQuestion := project.Test.Questions[currentQuestionIndex+1]
		c.JSON(http.StatusOK, gin.H{"question": project.DisplayQuestion(nextQuestion)})
	} else {
		c.JSON(http.StatusOK, gin.H{"message": "Test completed"})
	}
}

func GetScoreHandler(c *gin.Context) {
	session, exists := c.Get("session")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Session not found"})
		return
	}

	project := c.MustGet("project").(*models.Project)

	// Extract user answers from the session
	userAnswers := make([]string, len(project.Test.Questions))
	for i, question := range project.Test.Questions {
		userAnswers[i] = session.(map[string]string)[question.Text]
	}

	// Calculate and display the final score
	score := project.CalculateScore(userAnswers)
	c.JSON(http.StatusOK, gin.H{"score": score, "totalQuestions": len(project.Test.Questions)})
}

// StartTestHandlerCLI runs the test engine in CLI mode.
func StartTestHandlerCLI(project *models.Project, userInputFunc func() string) {
	if project.Test == nil || len(project.Test.Questions) == 0 {
		fmt.Println("No questions loaded.")
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
		answer := userInputFunc()

		// Validate user input
		for !validator.IsValidOption(answer) {
			fmt.Println("Invalid answer option. Please enter a valid option.")
			fmt.Print("Your answer: ")
			answer = userInputFunc()
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
