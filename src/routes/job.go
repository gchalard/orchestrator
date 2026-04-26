package routes

import (
	"net/http"

	workflowroutes "github.com/gchalard/orchestrator/src/routes/workflows"
)

func GetJobDefinitionHandler(w http.ResponseWriter, r *http.Request) {
	workflowroutes.GetJobDefinitionHandler(w, r)
}
