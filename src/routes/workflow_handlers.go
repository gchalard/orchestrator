package routes

import (
	"net/http"

	workflowroutes "github.com/gchalard/orchestrator/src/routes/workflows"
)

func CreateWorkflowHandler(w http.ResponseWriter, r *http.Request) {
	workflowroutes.CreateWorkflowHandler(w, r)
}

func GetWorkflowsHandler(w http.ResponseWriter, r *http.Request) {
	workflowroutes.GetWorkflowsHandler(w, r)
}

func TriggerWorkflowsHandler(w http.ResponseWriter, r *http.Request) {
	workflowroutes.TriggerWorkflowsHandler(w, r)
}
