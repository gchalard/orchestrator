package workflows

import "github.com/gchalard/orchestrator/src/models"

type triggerPayload struct {
	WorkflowName string `json:"workflow_name"`
}

type runState struct {
	pendingJobs   map[string]models.Job
	spawnedJobs   map[string]string
	completedJobs map[string]string
}

func newRunState(jobs map[string]models.Job) runState {
	pending := make(map[string]models.Job, len(jobs))
	for id, job := range jobs {
		pending[id] = job
	}

	return runState{
		pendingJobs:   pending,
		spawnedJobs:   make(map[string]string, len(jobs)),
		completedJobs: make(map[string]string, len(jobs)),
	}
}
