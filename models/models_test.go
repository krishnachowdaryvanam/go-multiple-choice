package models

import (
	"testing"
)

func TestDisplayQuestion(t *testing.T) {
	// Test displaying a question and choices
	question := Question{
		Text:       "What is 2 + 2?",
		Options:    []string{"3", "4", "5"},
		CorrectAns: "4",
	}

	project := NewProject(nil)
	output := project.DisplayQuestion(question)

	expectedOutput := "Question: What is 2 + 2?\nA. 3\nB. 4\nC. 5\nYour answer: "
	if output != expectedOutput {
		t.Errorf("Displayed output does not match the expected content")
	}
}

func TestValidateAnswerCorrect(t *testing.T) {
	// Test validating correct answers
	question := Question{
		Text:       "What is 2 + 2?",
		Options:    []string{"3", "4", "5"},
		CorrectAns: "4",
	}

	project := NewProject(&Test{Questions: []Question{question}})
	result := project.ValidateAnswer("4", question.CorrectAns)

	if !result {
		t.Error("Expected correct answer validation to return true")
	}
}

func TestValidateAnswerIncorrect(t *testing.T) {
	// Test validating incorrect answers
	question := Question{
		Text:       "What is 2 + 2?",
		Options:    []string{"3", "4", "5"},
		CorrectAns: "4",
	}

	project := NewProject(&Test{Questions: []Question{question}})
	result := project.ValidateAnswer("3", question.CorrectAns)

	if result {
		t.Error("Expected incorrect answer validation to return false")
	}
}

func TestCalculateScore(t *testing.T) {
	// Test calculating the score
	questions := []Question{
		{Text: "Q1", Options: []string{"A", "B"}, CorrectAns: "A"},
		{Text: "Q2", Options: []string{"X", "Y"}, CorrectAns: "Y"},
	}

	project := NewProject(&Test{Questions: questions})
	userAnswers := []string{"A", "Y"}

	score := project.CalculateScore(userAnswers)

	// Assert that the score is calculated correctly
	expectedScore := 2 // Both answers are correct
	if score != expectedScore {
		t.Errorf("Calculated score does not match the expected value")
	}
}
