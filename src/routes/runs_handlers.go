package routes

import (
	"net/http"

	runroutes "github.com/gchalard/orchestrator/src/routes/runs"
)

func GetRunsHandler(w http.ResponseWriter, r *http.Request) {
	runroutes.ListRunsHandler(w, r)
}

func GetRunHandler(w http.ResponseWriter, r *http.Request) {
	runroutes.GetRunHandler(w, r)
}
