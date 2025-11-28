package presentation

import (
	"encoding/json"
	"net/http"

	"github.com/gsousadev/doolar2/internal/tasks/application"
	task_list "github.com/gsousadev/doolar2/internal/tasks/domain/entity"
)

// TaskManagerHandler é o handler HTTP para gerenciamento de tasks
// Depende da interface TaskManager, não da implementação concreta
type TaskManagerHandler struct {
	service application.TaskManager
}

// NewTaskManagerHandler cria uma nova instância do handler
func NewTaskManagerHandler(service application.TaskManager) *TaskManagerHandler {
	return &TaskManagerHandler{
		service: service,
	}
}

// CreateTaskListRequest representa a requisição de criação
type CreateTaskListRequest struct {
	Title string `json:"title"`
}

// CreateTaskRequest representa a requisição para adicionar task
type CreateTaskRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

// UpdateTaskStatusRequest representa a requisição de atualização de status
type UpdateTaskStatusRequest struct {
	Status string `json:"status"`
}

// TaskListResponse - DTO de resposta da lista
type TaskListResponse struct {
	ID    string         `json:"id"`
	Title string         `json:"title"`
	Tasks []TaskResponse `json:"tasks"`
	Stats StatsResponse  `json:"stats"`
}

// TaskResponse - DTO de task individual
type TaskResponse struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      string `json:"status"`
}

// StatsResponse - Estatísticas da lista
type StatsResponse struct {
	Total      int `json:"total"`
	Pending    int `json:"pending"`
	InProgress int `json:"in_progress"`
	Completed  int `json:"completed"`
	Cancelled  int `json:"cancelled"`
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

// CreateTaskList godoc
// @Summary Criar uma nova lista de tarefas
// @Description Cria uma nova lista de tarefas vazia
// @Tags task-lists
// @Accept json
// @Produce json
// @Param request body CreateTaskListRequest true "Dados da lista"
// @Success 201 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /task-lists [post]
func (h *TaskManagerHandler) CreateTaskList(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var req CreateTaskListRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.Title == "" {
		respondError(w, http.StatusBadRequest, "Title is required")
		return
	}

	dto := application.CreateTaskListDTO{
		Title: req.Title,
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

// GetTaskList godoc
// @Summary Buscar lista de tarefas
// @Description Retorna uma lista de tarefas completa com todas as tasks e estatísticas
// @Tags task-lists
// @Produce json
// @Param id path string true "Task List ID"
// @Success 200 {object} SuccessResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /task-lists/{id} [get]
func (h *TaskManagerHandler) GetTaskList(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	// Extrai ID da URL (assumindo pattern: /task-lists/{id})
	id := extractIDFromPath(r.URL.Path, "/task-lists/")
	if id == "" {
		respondError(w, http.StatusBadRequest, "Invalid task list ID")
		return
	}

	taskList, err := h.service.GetTaskList(id)
	if err != nil {
		if err == application.ErrTaskListNotFound {
			respondError(w, http.StatusNotFound, "Task list not found")
			return
		}
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Transforma entidade em DTO na camada de apresentação
	response := mapTaskListToResponse(taskList)
	respondSuccess(w, http.StatusOK, "Task list retrieved successfully", response)
}

// AddTaskToList godoc
// @Summary Adicionar task a uma lista
// @Description Adiciona uma nova task a uma lista existente
// @Tags tasks
// @Accept json
// @Produce json
// @Param id path string true "Task List ID"
// @Param request body CreateTaskRequest true "Dados da task"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /task-lists/{id}/tasks [post]
func (h *TaskManagerHandler) AddTaskToList(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	id := extractIDFromPath(r.URL.Path, "/task-lists/")
	if id == "" {
		respondError(w, http.StatusBadRequest, "Invalid task list ID")
		return
	}

	var req CreateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.Title == "" {
		respondError(w, http.StatusBadRequest, "Title is required")
		return
	}

	dto := application.CreateTaskDTO{
		Title:       req.Title,
		Description: req.Description,
	}

	taskList, err := h.service.AddTaskToList(id, dto)
	if err != nil {
		if err == application.ErrTaskListNotFound {
			respondError(w, http.StatusNotFound, "Task list not found")
			return
		}
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Transforma entidade em DTO na camada de apresentação
	response := mapTaskListToResponse(taskList)
	respondSuccess(w, http.StatusOK, "Task added successfully", response)
}

// GetPendingTasks godoc
// @Summary Listar tasks pendentes
// @Description Retorna todas as tasks com status pendente de uma lista
// @Tags tasks
// @Produce json
// @Param id path string true "Task List ID"
// @Success 200 {object} SuccessResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /task-lists/{id}/tasks/pending [get]
func (h *TaskManagerHandler) GetPendingTasks(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	id := extractIDFromPath(r.URL.Path, "/task-lists/")
	if id == "" {
		respondError(w, http.StatusBadRequest, "Invalid task list ID")
		return
	}

	tasks, err := h.service.GetPendingTasks(id)
	if err != nil {
		if err == application.ErrTaskListNotFound {
			respondError(w, http.StatusNotFound, "Task list not found")
			return
		}
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Transforma entidades em DTOs na camada de apresentação
	response := mapTasksToResponse(tasks)
	respondSuccess(w, http.StatusOK, "Pending tasks retrieved successfully", response)
}

// GetStatistics godoc
// @Summary Obter estatísticas da lista
// @Description Retorna estatísticas de uma lista de tarefas
// @Tags task-lists
// @Produce json
// @Param id path string true "Task List ID"
// @Success 200 {object} SuccessResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /task-lists/{id}/statistics [get]
func (h *TaskManagerHandler) GetStatistics(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	id := extractIDFromPath(r.URL.Path, "/task-lists/")
	if id == "" {
		respondError(w, http.StatusBadRequest, "Invalid task list ID")
		return
	}

	taskList, err := h.service.GetTaskListForStats(id)
	if err != nil {
		if err == application.ErrTaskListNotFound {
			respondError(w, http.StatusNotFound, "Task list not found")
			return
		}
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Calcula estatísticas na camada de apresentação
	stats := calculateStats(taskList)
	respondSuccess(w, http.StatusOK, "Statistics retrieved successfully", stats)
}

// UpdateTaskStatus godoc
// @Summary Atualizar status de uma task
// @Description Atualiza o status de uma task específica
// @Tags tasks
// @Accept json
// @Produce json
// @Param listId path string true "Task List ID"
// @Param taskId path string true "Task ID"
// @Param request body UpdateTaskStatusRequest true "Novo status"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /task-lists/{listId}/tasks/{taskId}/status [patch]
func (h *TaskManagerHandler) UpdateTaskStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPatch {
		respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	// Extrair IDs da URL
	listID := extractIDFromPath(r.URL.Path, "/task-lists/")
	taskID := extractTaskIDFromPath(r.URL.Path)

	if listID == "" || taskID == "" {
		respondError(w, http.StatusBadRequest, "Invalid IDs")
		return
	}

	var req UpdateTaskStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.Status == "" {
		respondError(w, http.StatusBadRequest, "Status is required")
		return
	}

	err := h.service.UpdateTaskStatus(listID, taskID, req.Status)
	if err != nil {
		if err == application.ErrTaskListNotFound {
			respondError(w, http.StatusNotFound, "Task list not found")
			return
		}
		if err == application.ErrTaskNotFound {
			respondError(w, http.StatusNotFound, "Task not found")
			return
		}
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondSuccess(w, http.StatusOK, "Task status updated successfully", nil)
}

// DeleteTaskList godoc
// @Summary Deletar lista de tarefas
// @Description Remove uma lista de tarefas e todas as suas tasks
// @Tags task-lists
// @Produce json
// @Param id path string true "Task List ID"
// @Success 200 {object} SuccessResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /task-lists/{id} [delete]
func (h *TaskManagerHandler) DeleteTaskList(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	id := extractIDFromPath(r.URL.Path, "/task-lists/")
	if id == "" {
		respondError(w, http.StatusBadRequest, "Invalid task list ID")
		return
	}

	err := h.service.DeleteTaskList(id)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondSuccess(w, http.StatusOK, "Task list deleted successfully", nil)
}

// Helper functions
func respondError(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(ErrorResponse{
		Error:   http.StatusText(statusCode),
		Message: message,
	})
}

func respondSuccess(w http.ResponseWriter, statusCode int, message string, data interface{}) {
	json.NewEncoder(w).Encode(SuccessResponse{
		Message: message,
		Data:    data,
	})
}

func extractIDFromPath(path string, prefix string) string {
	// Remove o prefixo e extrai o ID
	// Exemplo: /task-lists/uuid-here/tasks → uuid-here
	if len(path) <= len(prefix) {
		return ""
	}

	remaining := path[len(prefix):]

	// Pega até a próxima barra
	for i, c := range remaining {
		if c == '/' {
			return remaining[:i]
		}
	}

	return remaining
}

func extractTaskIDFromPath(path string) string {
	// Exemplo: /task-lists/list-uuid/tasks/task-uuid/status
	parts := splitPath(path)

	for i, part := range parts {
		if part == "tasks" && i+1 < len(parts) {
			return parts[i+1]
		}
	}

	return ""
}

func splitPath(path string) []string {
	var parts []string
	current := ""

	for _, c := range path {
		if c == '/' {
			if current != "" {
				parts = append(parts, current)
				current = ""
			}
		} else {
			current += string(c)
		}
	}

	if current != "" {
		parts = append(parts, current)
	}

	return parts
}

// Mapper functions - transformam entidades em DTOs
func mapTaskListToResponse(taskList *task_list.TaskListEntity) *TaskListResponse {
	tasks := make([]TaskResponse, len(taskList.Tasks))
	stats := StatsResponse{Total: len(taskList.Tasks)}

	for i, task := range taskList.Tasks {
		taskEntity := task.(*task_list.TaskEntity)
		tasks[i] = TaskResponse{
			ID:          task.GetID().String(),
			Title:       taskEntity.Title,
			Description: taskEntity.Description,
			Status:      string(task.GetStatus()),
		}

		// Calcula stats
		switch task.GetStatus() {
		case task_list.StatusPending:
			stats.Pending++
		case task_list.StatusInProgress:
			stats.InProgress++
		case task_list.StatusCompleted:
			stats.Completed++
		case task_list.StatusCancelled:
			stats.Cancelled++
		}
	}

	return &TaskListResponse{
		ID:    taskList.ID.String(),
		Title: taskList.Title,
		Tasks: tasks,
		Stats: stats,
	}
}

func mapTasksToResponse(tasks []task_list.ITask) []TaskResponse {
	response := make([]TaskResponse, len(tasks))

	for i, task := range tasks {
		taskEntity := task.(*task_list.TaskEntity)
		response[i] = TaskResponse{
			ID:          task.GetID().String(),
			Title:       taskEntity.Title,
			Description: taskEntity.Description,
			Status:      string(task.GetStatus()),
		}
	}

	return response
}

func calculateStats(taskList *task_list.TaskListEntity) *StatsResponse {
	stats := &StatsResponse{
		Total: len(taskList.Tasks),
	}

	for _, task := range taskList.Tasks {
		switch task.GetStatus() {
		case task_list.StatusPending:
			stats.Pending++
		case task_list.StatusInProgress:
			stats.InProgress++
		case task_list.StatusCompleted:
			stats.Completed++
		case task_list.StatusCancelled:
			stats.Cancelled++
		}
	}

	return stats
}
