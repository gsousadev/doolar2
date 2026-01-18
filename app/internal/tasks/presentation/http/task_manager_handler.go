package http

import (
	"encoding/json"
	"net/http"

	"github.com/gsousadev/doolar-golang/internal/tasks/application/contracts"
	"github.com/gsousadev/doolar-golang/internal/tasks/application/dtos"
	"github.com/gsousadev/doolar-golang/internal/tasks/domain/entity"
)

type TaskManagerHandler struct {
	service contracts.ITaskManagerService
}

// NewTaskManagerHandler cria uma nova instância do handler
func NewTaskManagerHandler(service contracts.ITaskManagerService) *TaskManagerHandler {
	return &TaskManagerHandler{
		service: service,
	}
}

// CreateTaskListRequest representa a requisição de criação
type CreateTaskListRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

// ErrorResponse representa uma resposta de erro
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

// SuccessResponse representa uma resposta de sucesso
type SuccessResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func (h *TaskManagerHandler) CreateTaskList(w http.ResponseWriter, r *http.Request) {

	var req CreateTaskListRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.Title == "" {
		respondError(w, http.StatusBadRequest, "Title is required")
		return
	}

	dto := dtos.CreateTaskListDTO{
		Title:       req.Title,
		Description: req.Description,
	}

	taskList, err := h.service.CreateTaskList(dto)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Transforma entidade em DTO na camada de apresentação
	response := mapTaskListToResponse(taskList)
	respondSuccess(w, http.StatusCreated, "Task list created successfully", response)
}

func respondError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(ErrorResponse{
		Error:   http.StatusText(status),
		Message: message,
	})
}

func respondSuccess(w http.ResponseWriter, status int, message string, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(SuccessResponse{
		Message: message,
		Data:    data,
	})
}

func mapTaskListToResponse(taskList *entity.TaskListEntity) map[string]interface{} {
	tasks := []map[string]interface{}{}
	for _, task := range taskList.Tasks {
		tasks = append(tasks, map[string]interface{}{
			"id":          task.GetID(),
			"title":       task.GetTitle(),
			"description": task.GetDescription(),
			"status":      task.GetStatus(),
		})
	}

	return map[string]interface{}{
		"id":    taskList.GetID(),
		"title": taskList.Title,
		"tasks": tasks,
	}
}
