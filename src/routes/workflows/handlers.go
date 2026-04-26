package workflows

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gchalard/orchestrator/src/models"
	"github.com/gchalard/orchestrator/src/routes/runs"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

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
	var payload triggerPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil || payload.WorkflowName == "" {
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

	if err := validateJobDependencies(workflow.Jobs); err != nil {
		http.Error(w, fmt.Sprintf("invalid workflow dependencies: %v", err), http.StatusBadRequest)
		return
	}

	runID := uuid.New().String()
	runs.Start(runID, payload.WorkflowName)

	go func() {
		if err := RunWorkflow(payload.WorkflowName, workflow, runID); err != nil {
			runs.Fail(runID)
			return
		}
		runs.Succeed(runID)
	}()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	_ = json.NewEncoder(w).Encode(map[string]string{
		"run_id":        runID,
		"workflow_name": payload.WorkflowName,
		"status":        runs.StatusRunning,
	})
}

func GetJobDefinitionHandler(w http.ResponseWriter, r *http.Request) {
	workflowName := chi.URLParam(r, "workflowName")
	jobID := chi.URLParam(r, "jobID")

	workflowStore.mu.RLock()
	workflow, workflowOK := workflowStore.data[workflowName]
	workflowStore.mu.RUnlock()
	if !workflowOK {
		http.Error(w, "workflow not found", http.StatusNotFound)
		return
	}

	job, jobOK := workflow.Jobs[jobID]
	if !jobOK {
		http.Error(w, "job not found", http.StatusNotFound)
		return
	}

	jobData, err := job.Get()
	if err != nil {
		http.Error(w, "failed to serialize job", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(jobData)
}
