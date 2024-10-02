// main.go
package main

import (
	"github.com/kaytu-io/cerberus/scheduler"
	"github.com/kaytu-io/cerberus/server"
	"github.com/kaytu-io/cerberus/workflow"
)

func main() {
	// Create a new WorkflowManager
	wm := workflow.NewWorkflowManager()

	// Start the scheduler
	scheduler.StartScheduler(wm)

	// Start the web server
	server.StartServer(wm)
}
