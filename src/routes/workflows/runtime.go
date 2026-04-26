package workflows

import (
	"fmt"
	"time"

	"github.com/gchalard/orchestrator/src/models"
	"github.com/gchalard/orchestrator/src/service"
)

// RunWorkflow executes all jobs for a workflow run using the given runID (Kubernetes job name suffix).
func RunWorkflow(workflowName string, workflow models.Workflow, runID string) error {
	if err := validateJobDependencies(workflow.Jobs); err != nil {
		return fmt.Errorf("invalid workflow dependencies: %w", err)
	}

	state := newRunState(workflow.Jobs)
	for len(state.completedJobs) < len(workflow.Jobs) {
		progressed, err := spawnRunnableJobs(workflowName, runID, &state)
		if err != nil {
			return err
		}

		statusProgressed, err := refreshSpawnedJobs(&state)
		if err != nil {
			return err
		}
		progressed = progressed || statusProgressed

		if len(state.completedJobs) == len(workflow.Jobs) {
			break
		}

		if !progressed {
			if len(state.spawnedJobs) == 0 {
				return fmt.Errorf("workflow deadlock detected: no runnable jobs found (possible cyclic dependencies)")
			}
			time.Sleep(2 * time.Second)
		}
	}

	return nil
}

func spawnRunnableJobs(workflowName, runID string, state *runState) (bool, error) {
	progressed := false
	for jobID, job := range state.pendingJobs {
		if dep := failedDependency(job.Needs, state.completedJobs); dep != "" {
			return false, fmt.Errorf("job %s cannot be spawned because dependency %s failed", jobID, dep)
		}
		if !dependenciesSuccessful(job.Needs, state.completedJobs) {
			continue
		}

		if err := service.SpawnJob(workflowName, jobID, runID, job.RunsOn); err != nil {
			return false, fmt.Errorf("failed to spawn job %s: %w", jobID, err)
		}

		state.spawnedJobs[jobID] = fmt.Sprintf("%s-%s-%s", workflowName, runID, jobID)
		delete(state.pendingJobs, jobID)
		progressed = true
	}
	return progressed, nil
}

func refreshSpawnedJobs(state *runState) (bool, error) {
	progressed := false
	for jobID, spawnedJobName := range state.spawnedJobs {
		status, err := service.GetJobStatus(spawnedJobName)
		if err != nil {
			return false, fmt.Errorf("failed to get status for job %s: %w", spawnedJobName, err)
		}
		switch status {
		case "succeeded":
			state.completedJobs[jobID] = status
			delete(state.spawnedJobs, jobID)
			progressed = true
		case "failed":
			state.completedJobs[jobID] = status
			return false, fmt.Errorf("job %s failed", jobID)
		}
	}
	return progressed, nil
}
