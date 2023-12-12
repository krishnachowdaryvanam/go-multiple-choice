package routes

import (
	"multiple-choice-test-engine/handler"
	"multiple-choice-test-engine/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine, project *models.Project) {
	router.GET("/loadTest", LoadTestHandler(project)) // Fix function name
	router.GET("/startTest", handler.StartTestHandler)
	router.POST("/submitAnswer", handler.SubmitAnswerHandler)
	router.GET("/getScore", handler.GetScoreHandler)
}

func LoadTestHandler(project *models.Project) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Load test from JSON file (replace with your correct file path)
		err := project.LoadTestFromJSON("C:/Users/Dell/OneDrive/Desktop/GO/multiple-choice-test-engine/test.json")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Test loaded successfully"})
	}
}
