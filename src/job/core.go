package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"sync"
)

func getJobDefinition(workflowName, jobID, orchestratorEndpoint string) ([]map[string]interface{}, error) {
	url := fmt.Sprintf("%s/workflows/%s/jobs/%s", orchestratorEndpoint, workflowName, jobID)
	response, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to get job definition: %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get job definition: status code %d", response.StatusCode)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read job definition: %w", err)
	}

	var payload struct {
		Steps []map[string]interface{} `json:"steps"`
	}
	if err := json.Unmarshal(body, &payload); err != nil {
		return nil, fmt.Errorf("failed to decode job definition: %w", err)
	}

	return payload.Steps, nil
}

func executeStep(step map[string]interface{}) error {
	command, ok := step["run"].(string)
	if !ok || command == "" {
		return fmt.Errorf("step is missing a non-empty run command")
	}

	shell := "bash"
	if value, ok := step["shell"].(string); ok && value != "" {
		shell = value
	}

	workingDirectory := "."
	if value, ok := step["working_directory"].(string); ok && value != "" {
		workingDirectory = value
	}

	envVars := os.Environ()
	if rawEnv, ok := step["env"].(map[string]interface{}); ok {
		for key, value := range rawEnv {
			stringValue, ok := value.(string)
			if !ok {
				return fmt.Errorf("env var %q must be a string", key)
			}
			envVars = append(envVars, fmt.Sprintf("%s=%s", key, stringValue))
		}
	}

	cmd := exec.Command(shell, "-c", command)
	cmd.Dir = workingDirectory
	cmd.Env = envVars

	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("failed to open stdout pipe: %w", err)
	}

	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("failed to open stderr pipe: %w", err)
	}

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start step command: %w", err)
	}

	var wg sync.WaitGroup
	var outputMu sync.Mutex
	var outputBuilder strings.Builder
	var streamErrMu sync.Mutex
	var streamErr error

	streamOutput := func(reader io.Reader) {
		defer wg.Done()
		scanner := bufio.NewScanner(reader)
		for scanner.Scan() {
			line := scanner.Text()
			log.Printf("step output:\t%s", line)

			outputMu.Lock()
			outputBuilder.WriteString(line)
			outputBuilder.WriteByte('\n')
			outputMu.Unlock()
		}

		if err := scanner.Err(); err != nil {
			streamErrMu.Lock()
			if streamErr == nil {
				streamErr = err
			}
			streamErrMu.Unlock()
		}
	}

	wg.Add(2)
	go streamOutput(stdoutPipe)
	go streamOutput(stderrPipe)

	cmdErr := cmd.Wait()
	wg.Wait()

	if streamErr != nil {
		return fmt.Errorf("failed to stream command output: %w", streamErr)
	}

	if cmdErr != nil {
		outputMu.Lock()
		output := outputBuilder.String()
		outputMu.Unlock()
		return fmt.Errorf("failed to execute step: %w (output: %s)", cmdErr, output)
	}

	return nil
}

func Main() error {
	orchestratorEndpoint := os.Getenv("ORCHESTRATOR_ENDPOINT")
	workflowName := os.Getenv("WORKFLOW_NAME")
	jobID := os.Getenv("JOB_ID")

	steps, err := getJobDefinition(workflowName, jobID, orchestratorEndpoint)
	if err != nil {
		return fmt.Errorf("failed to get job definition: %w", err)
	}

	for _, step := range steps {
		err := executeStep(step)
		if err != nil {
			return fmt.Errorf("failed to execute step: %w", err)
		}
	}

	return nil
}

func main() {
	err := Main()
	if err != nil {
		log.Fatalf("failed to run job: %v", err)
	}
	os.Exit(0)
}
