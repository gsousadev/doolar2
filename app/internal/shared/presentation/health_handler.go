package presentation

import "net/http"

type HealthHandler struct{}

func (h HealthHandler) HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}
