package models

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
)

// Question represents each question in the test.
type Question struct {
	Text       string   `json:"question"`
	Options    []string `json:"choices"`
	CorrectAns string   `json:"correctAnswer"`
}

// Test represents the collection of questions.
type Test struct {
	Questions []Question `json:"questions"`
}

// Project encapsulates the test engine functionalities.
type Project struct {
	Test *Test
}

// NewProject creates a new Project instance with a given test.
func NewProject(test *Test) *Project {
	return &Project{Test: test}
}

// LoadTestFromJSON loads a test from a JSON file.
func (p *Project) LoadTestFromJSON(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&p.Test); err != nil {
		return fmt.Errorf("error decoding JSON: %v", err)
	}

	return nil
}

// DisplayQuestion displays a question and its choices.
func (p *Project) DisplayQuestion(question Question) string {
	output := fmt.Sprintf("Question: %s\n", question.Text)
	for i, option := range question.Options {
		output += fmt.Sprintf("%c. %s\n", 'A'+i, option)
	}
	output += "Your answer: "
	return output
}

// ValidateAnswer checks if the user's answer is correct.
func (p *Project) ValidateAnswer(userAnswer, correctAnswer string) bool {
	return strings.EqualFold(userAnswer, correctAnswer)
}

// CalculateScore calculates the user's final score.
func (p *Project) CalculateScore(userAnswers []string) int {
	score := 0
	for i, question := range p.Test.Questions {
		if p.ValidateAnswer(userAnswers[i], question.CorrectAns) {
			score++
		}
	}
	return score
}

// ShuffleQuestions shuffles the order of questions in the test.
func (p *Project) ShuffleQuestions() {
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(p.Test.Questions), func(i, j int) {
		p.Test.Questions[i], p.Test.Questions[j] = p.Test.Questions[j], p.Test.Questions[i]
	})
}

// SetTimer sets a timer for the test duration (in seconds).
func (p *Project) SetTimer(duration int) <-chan struct{} {
	timer := time.NewTimer(time.Duration(duration) * time.Second)
	done := make(chan struct{})
	go func() {
		<-timer.C
		close(done)
	}()
	return done
}
