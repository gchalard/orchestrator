package models

import "fmt"

type Job struct {
	Name   string   `json:"name"`
	Needs  []string `json:"needs,omitempty"`
	RunsOn string   `json:"runs_on"`
	Steps  []Step   `json:"steps"`
}

func (j *Job) Get() (map[string]interface{}, error) {
	steps := make([]map[string]interface{}, 0, len(j.Steps))
	for _, step := range j.Steps {
		stepData, err := step.Get()
		if err != nil {
			return nil, fmt.Errorf("failed to get step: %w", err)
		}
		steps = append(steps, stepData)
	}

	return map[string]interface{}{
		"name":    j.Name,
		"needs":   j.Needs,
		"runs_on": j.RunsOn,
		"steps":   steps,
	}, nil
}
