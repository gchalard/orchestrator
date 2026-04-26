package workflows

import (
	"fmt"

	"github.com/gchalard/orchestrator/src/models"
)

func validateJobDependencies(jobs map[string]models.Job) error {
	for jobID, job := range jobs {
		for _, dep := range job.Needs {
			if dep == jobID {
				return fmt.Errorf("job %s cannot depend on itself", jobID)
			}
			if _, exists := jobs[dep]; !exists {
				return fmt.Errorf("job %s depends on unknown job %s", jobID, dep)
			}
		}
	}
	return nil
}

func dependenciesSuccessful(deps []string, completed map[string]string) bool {
	for _, dep := range deps {
		if completed[dep] != "succeeded" {
			return false
		}
	}
	return true
}

func failedDependency(deps []string, completed map[string]string) string {
	for _, dep := range deps {
		if completed[dep] == "failed" {
			return dep
		}
	}
	return ""
}
