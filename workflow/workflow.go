// workflow/workflow.go
package workflow

import (
	"fmt"
	"log"

	"github.com/kaytu-io/cerberus/execution"
	"github.com/rhosocial/go-dag/pkg/dag"
)

// Define the workflow names
const (
	DefaultWorkflow     = "default"
	EgressCheckWorkflow = "egress_check"
)

// WorkflowManager manages the workflows
type WorkflowManager struct {
	Workflows map[string]*dag.DAG
}

// NewWorkflowManager creates a new WorkflowManager
func NewWorkflowManager() *WorkflowManager {
	wm := &WorkflowManager{
		Workflows: make(map[string]*dag.DAG),
	}
	wm.initializeWorkflows()
	return wm
}

// initializeWorkflows sets up the workflows
func (wm *WorkflowManager) initializeWorkflows() {
	// Initialize the default workflow
	wm.Workflows[DefaultWorkflow] = wm.createDefaultWorkflow()

	// Initialize the egress check workflow
	wm.Workflows[EgressCheckWorkflow] = wm.createEgressCheckWorkflow()
}

// createDefaultWorkflow creates the default workflow with both tasks
func (wm *WorkflowManager) createDefaultWorkflow() *dag.DAG {
	d := dag.NewDAG()

	// Define tasks
	checkEgressTask := &dag.Task{
		Name: "CheckInternetEgress",
		Func: func() error {
			return execution.CheckInternetEgress()
		},
	}

	fetchReleaseTask := &dag.Task{
		Name: "FetchLatestRelease",
		Func: func() error {
			owner := "kaytu-io"
			repo := "managed-platform-config"
			token := "" // Add your GitHub token if needed
			return execution.FetchLatestRelease(owner, repo, token)
		},
	}

	// Add tasks to the DAG
	d.AddTask(checkEgressTask)
	d.AddTask(fetchReleaseTask)

	// Define dependencies
	// FetchLatestRelease depends on CheckInternetEgress
	d.AddDependency(fetchReleaseTask, checkEgressTask)

	return d
}

// createEgressCheckWorkflow creates a workflow with only the CheckInternetEgress task
func (wm *WorkflowManager) createEgressCheckWorkflow() *dag.DAG {
	d := dag.NewDAG()

	// Define the task
	checkEgressTask := &dag.Task{
		Name: "CheckInternetEgress",
		Func: func() error {
			return execution.CheckInternetEgress()
		},
	}

	// Add the task to the DAG
	d.AddTask(checkEgressTask)

	return d
}

// RunWorkflow runs the specified workflow
func (wm *WorkflowManager) RunWorkflow(name string) error {
	if workflow, exists := wm.Workflows[name]; exists {
		results := workflow.Run()
		for _, res := range results {
			if res.Error != nil {
				log.Printf("Task %s failed: %v\n", res.Task.Name, res.Error)
				return res.Error
			}
			log.Printf("Task %s completed successfully.\n", res.Task.Name)
		}
		return nil
	}
	return fmt.Errorf("workflow %s not found", name)
}
