// main.go
package main

import (
	"flag"
	"multiple-choice-test-engine/cli"
	"multiple-choice-test-engine/models"
	"multiple-choice-test-engine/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	// Command-line flag to determine the interface to run
	useCLI := flag.Bool("cli", false, "Run the CLI interface")

	flag.Parse()

	if *useCLI {
		// Run the CLI interface
		runCLI()
	} else {
		// Run the HTTP server
		runHTTPServer()
	}
}

func runCLI() {
	project := models.NewProject(nil)

	// Run the CLI interface
	cli.RunCLI(project)
}

func runHTTPServer() {
	router := gin.Default()

	// Initialize a new project
	project := models.NewProject(nil)

	// Set up routes
	routes.SetupRoutes(router, project)

	router.Run(":8080")
}
