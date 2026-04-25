package routes

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"

	"github.com/gchalard/orchestrator/src/models"
	"github.com/gchalard/orchestrator/src/service"
)

var workflowStore = struct {
	mu   sync.RWMutex
	data map[string]models.Workflow
}{
	data: make(map[string]models.Workflow),
}

func CreateWorkflowHandler(w http.ResponseWriter, r *http.Request) {
	var workflow models.Workflow
	if err := json.NewDecoder(r.Body).Decode(&workflow); err != nil {
		http.Error(w, "invalid workflow payload", http.StatusBadRequest)
		return
	}
	if workflow.Name == "" {
		http.Error(w, "workflow name is required", http.StatusBadRequest)
		return
	}

	workflowStore.mu.Lock()
	workflowStore.data[workflow.Name] = workflow
	workflowStore.mu.Unlock()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(workflow)
}

func GetWorkflowsHandler(w http.ResponseWriter, r *http.Request) {
	workflowStore.mu.RLock()
	workflows := make([]models.Workflow, 0, len(workflowStore.data))
	for _, workflow := range workflowStore.data {
		workflows = append(workflows, workflow)
	}
	workflowStore.mu.RUnlock()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(workflows)
}

func TriggerWorkflowsHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var payload struct {
		WorkflowName string `json:"workflow_name"`
	}
	if err := json.Unmarshal(body, &payload); err != nil || payload.WorkflowName == "" {
		http.Error(w, "invalid workflow payload", http.StatusBadRequest)
		return
	}

	workflowStore.mu.RLock()
	workflow, exists := workflowStore.data[payload.WorkflowName]
	workflowStore.mu.RUnlock()
	if !exists {
		http.Error(w, "workflow not found", http.StatusNotFound)
		return
	}

	type WorkflowJob struct {
		ID  string     `json:"id"`
		Job models.Job `json:"job"`
	}

	jobs := make([]WorkflowJob, 0, len(workflow.Jobs))
	for id, job := range workflow.Jobs {
		jobs = append(jobs, WorkflowJob{ID: id, Job: job})
	}

	for _, workflowJob := range jobs {
		if err = service.SpawnJob(fmt.Sprintf("%s-%s", payload.WorkflowName, workflowJob.ID), workflowJob.Job.RunsOn); err != nil {
			http.Error(w, fmt.Sprintf("failed to spawn job %s: %v", workflowJob.ID, err), http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
}
