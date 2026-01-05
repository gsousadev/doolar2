package cmd

import (
	"net/http"

	taskPresentation "github.com/gsousadev/doolar2/internal/tasks/presentation"
)

func setupRoutes() {
	taskRoute = taskPresentation.NewTaskManagerHandler()
}

func routes() map[string]http.Handler {

	return map[string]http.Handler{
		"POST /tasks": taskPresentation.CreateTaskList,
	}
}
