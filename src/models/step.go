package models

import "fmt"

type Step struct {
	Name             string            `json:"name"`
	Run              string            `json:"run"`
	WorkingDirectory string            `json:"working_directory,omitempty"`
	Env              map[string]string `json:"env,omitempty"`
	Shell            string            `json:"shell,omitempty"`
}

func (s *Step) Get() (map[string]interface{}, error) {
	env, err := s.getEnv()
	if err != nil {
		return nil, fmt.Errorf("failed to get env: %w", err)
	}
	return map[string]interface{}{
		"name":              s.Name,
		"run":               s.Run,
		"working_directory": s.WorkingDirectory,
		"env":               env,
		"shell":             s.Shell,
	}, nil
}

func (s *Step) getEnv() (map[string]string, error) {
	env := make(map[string]string)
	for key, value := range s.Env {
		env[key] = value
	}
	return env, nil
}
