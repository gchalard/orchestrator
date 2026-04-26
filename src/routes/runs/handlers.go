package runs

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func ListRunsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(List())
}

func GetRunHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "runId")
	run, ok := Get(id)
	if !ok {
		http.Error(w, "run not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(run)
}
