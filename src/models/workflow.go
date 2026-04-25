package models

import "fmt"

type Workflow struct {
	Name string         `json:"name"`
	Jobs map[string]Job `json:"jobs"`
}

func (w *Workflow) Get() (map[string]interface{}, error) {
	jobs := make(map[string]interface{}, len(w.Jobs))
	for jobName, job := range w.Jobs {
		jobData, err := job.Get()
		if err != nil {
			return nil, fmt.Errorf("failed to get job: %w", err)
		}
		jobs[jobName] = jobData
	}

	return map[string]interface{}{
		"name": w.Name,
		"jobs": jobs,
	}, nil
}
