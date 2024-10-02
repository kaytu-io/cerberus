// scheduler/job.go
package scheduler

import (
	"log"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/kaytu-io/cerberus/workflow"
)

var scheduler *gocron.Scheduler

// StartScheduler initializes and starts the scheduler
func StartScheduler(wm *workflow.WorkflowManager) {
	scheduler = gocron.NewScheduler(time.UTC)

	// Schedule the default workflow to run every 24 hours
	_, err := scheduler.Every(24).Hours().Do(func() {
		err := wm.RunWorkflow(workflow.DefaultWorkflow)
		if err != nil {
			log.Printf("Default Workflow execution failed: %v", err)
		} else {
			log.Printf("Default Workflow executed successfully.")
		}
	})
	if err != nil {
		log.Fatalf("Failed to schedule default workflow: %v", err)
	}

	// Schedule the egress check workflow to run every 24 hours
	_, err = scheduler.Every(24).Hours().Do(func() {
		err := wm.RunWorkflow(workflow.EgressCheckWorkflow)
		if err != nil {
			log.Printf("Egress Check Workflow execution failed: %v", err)
		} else {
			log.Printf("Egress Check Workflow executed successfully.")
		}
	})
	if err != nil {
		log.Fatalf("Failed to schedule egress check workflow: %v", err)
	}

	// Start the scheduler in a separate goroutine
	scheduler.StartAsync()
}
