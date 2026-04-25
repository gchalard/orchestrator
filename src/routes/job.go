package routes

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

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
