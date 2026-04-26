package main

import (
	"net/http"

	"github.com/gchalard/orchestrator/src/routes"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	router := chi.NewRouter()

	router.Use(middleware.Logger)

	router.Get("/", routes.RootHandler)
	router.Post("/echo", routes.EchoHandler)

	router.Get("/workflows", routes.GetWorkflowsHandler)
	router.Post("/workflows", routes.CreateWorkflowHandler)
	router.Post("/workflows/trigger", routes.TriggerWorkflowsHandler)

	router.Get("/runs", routes.GetRunsHandler)
	router.Get("/runs/{runId}", routes.GetRunHandler)

	router.Get("/workflows/{workflowName}/jobs/{jobID}", routes.GetJobDefinitionHandler)

	http.ListenAndServe(":8080", router)
}
