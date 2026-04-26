package runs

import (
	"sort"
	"sync"
)

const (
	StatusRunning    = "running"
	StatusSuccessful = "successful"
	StatusFailed     = "failed"
)

type Run struct {
	ID           string `json:"id"`
	WorkflowName string `json:"workflow_name"`
	Status       string `json:"status"`
}

var (
	mu   sync.RWMutex
	byID = make(map[string]Run)
)

// Start records a new run as running. Call before spawning background work.
func Start(id, workflowName string) {
	mu.Lock()
	defer mu.Unlock()
	byID[id] = Run{ID: id, WorkflowName: workflowName, Status: StatusRunning}
}

// Succeed marks the run successful. No-op if the run id is unknown.
func Succeed(id string) {
	setStatus(id, StatusSuccessful)
}

// Fail marks the run failed. No-op if the run id is unknown.
func Fail(id string) {
	setStatus(id, StatusFailed)
}

func setStatus(id, status string) {
	mu.Lock()
	defer mu.Unlock()
	if r, ok := byID[id]; ok {
		r.Status = status
		byID[id] = r
	}
}

// Get returns a run by id.
func Get(id string) (Run, bool) {
	mu.RLock()
	defer mu.RUnlock()
	r, ok := byID[id]
	return r, ok
}

// List returns all runs sorted by id for stable responses.
func List() []Run {
	mu.RLock()
	out := make([]Run, 0, len(byID))
	for _, r := range byID {
		out = append(out, r)
	}
	mu.RUnlock()
	sort.Slice(out, func(i, j int) bool { return out[i].ID < out[j].ID })
	return out
}
