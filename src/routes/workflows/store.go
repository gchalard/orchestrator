package workflows

import (
	"sync"

	"github.com/gchalard/orchestrator/src/models"
)

var workflowStore = struct {
	mu   sync.RWMutex
	data map[string]models.Workflow
}{
	data: make(map[string]models.Workflow),
}
