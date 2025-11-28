package application

import (
	"errors"

	task_list "github.com/gsousadev/doolar2/internal/tasks/domain/entity"
	"github.com/gsousadev/doolar2/internal/tasks/domain/repository"
)

var (
	ErrTaskListNotFound = errors.New("task list not found")
	ErrTaskNotFound     = errors.New("task not found")
	ErrInvalidStatus    = errors.New("invalid status")
)

// TaskManagerService é o serviço de aplicação que orquestra casos de uso
// Implementa a interface TaskManager
type TaskManagerService struct {
	repo repository.TaskListRepository
}

// NewTaskManagerService cria uma nova instância do serviço
func NewTaskManagerService(repo repository.TaskListRepository) TaskManager {
	return &TaskManagerService{
		repo: repo,
	}
}

// CreateTaskListDTO - DTO para criar uma lista
type CreateTaskListDTO struct {
	Title string `json:"title" validate:"required"`
}

// CreateTaskDTO - DTO para criar uma task
type CreateTaskDTO struct {
	Title       string `json:"title" validate:"required"`
	Description string `json:"description"`
}

// CreateTaskList cria uma nova lista de tarefas
func (s *TaskManagerService) CreateTaskList(dto CreateTaskListDTO) (*task_list.TaskListEntity, error) {
	taskList := task_list.NewTaskListEntity(dto.Title)

	if err := s.repo.Add(taskList); err != nil {
		return nil, err
	}

	if err := s.repo.Flush(); err != nil {
		return nil, err
	}

	return taskList, nil
}

// GetTaskList busca uma lista de tarefas por ID
func (s *TaskManagerService) GetTaskList(id string) (*task_list.TaskListEntity, error) {
	taskList, err := s.repo.FindByID(id)
	if err != nil {
		return nil, ErrTaskListNotFound
	}

	return taskList, nil
}

// AddTaskToList adiciona uma nova task a uma lista existente
func (s *TaskManagerService) AddTaskToList(listID string, dto CreateTaskDTO) (*task_list.TaskListEntity, error) {
	taskList, err := s.repo.FindByID(listID)
	if err != nil {
		return nil, ErrTaskListNotFound
	}

	task := task_list.NewTaskEntity(dto.Title, dto.Description)
	taskList.AddTask(task)

	if err := s.repo.Update(taskList); err != nil {
		return nil, err
	}

	if err := s.repo.Flush(); err != nil {
		return nil, err
	}

	return taskList, nil
}

// GetPendingTasks retorna apenas as tasks pendentes de uma lista
func (s *TaskManagerService) GetPendingTasks(listID string) ([]task_list.ITask, error) {
	taskList, err := s.repo.FindByID(listID)
	if err != nil {
		return nil, ErrTaskListNotFound
	}

	// Filtra tasks pendentes na camada de aplicação
	pending := make([]task_list.ITask, 0)
	for _, task := range taskList.Tasks {
		if task.GetStatus() == task_list.StatusPending {
			pending = append(pending, task)
		}
	}

	return pending, nil
}

// GetTasksByStatus retorna tasks filtradas por status
func (s *TaskManagerService) GetTasksByStatus(listID string, status string) ([]task_list.ITask, error) {
	taskList, err := s.repo.FindByID(listID)
	if err != nil {
		return nil, ErrTaskListNotFound
	}

	taskStatus := task_list.Status(status)

	// Valida status
	validStatuses := []task_list.Status{
		task_list.StatusPending,
		task_list.StatusInProgress,
		task_list.StatusCompleted,
		task_list.StatusCancelled,
	}

	isValid := false
	for _, s := range validStatuses {
		if s == taskStatus {
			isValid = true
			break
		}
	}

	if !isValid {
		return nil, ErrInvalidStatus
	}

	// Filtra na aplicação
	filtered := make([]task_list.ITask, 0)
	for _, task := range taskList.Tasks {
		if task.GetStatus() == taskStatus {
			filtered = append(filtered, task)
		}
	}

	return filtered, nil
}

// UpdateTaskStatus atualiza o status de uma task
func (s *TaskManagerService) UpdateTaskStatus(listID, taskID string, newStatus string) error {
	taskList, err := s.repo.FindByID(listID)
	if err != nil {
		return ErrTaskListNotFound
	}

	// Busca a task
	var targetTask task_list.ITask
	for _, task := range taskList.Tasks {
		if task.GetID().String() == taskID {
			targetTask = task
			break
		}
	}

	if targetTask == nil {
		return ErrTaskNotFound
	}

	// Muda o status
	if err := targetTask.ChangeStatus(task_list.Status(newStatus)); err != nil {
		return err
	}

	// Persiste
	if err := s.repo.Update(taskList); err != nil {
		return err
	}

	return s.repo.Flush()
}

// DeleteTaskList remove uma lista de tarefas
func (s *TaskManagerService) DeleteTaskList(id string) error {
	if err := s.repo.Remove(id); err != nil {
		return err
	}

	return s.repo.Flush()
}

// GetTaskList retorna a lista completa para cálculo de estatísticas
func (s *TaskManagerService) GetTaskListForStats(listID string) (*task_list.TaskListEntity, error) {
	taskList, err := s.repo.FindByID(listID)
	if err != nil {
		return nil, ErrTaskListNotFound
	}

	return taskList, nil
}
