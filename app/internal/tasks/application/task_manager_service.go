package application

import (
	"errors"

	"github.com/gsousadev/doolar2/internal/tasks/application/dtos"
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
	repo repository.ITaskListRepository
}

// NewTaskManagerService cria uma nova instância do serviço
func NewTaskManagerService(repo repository.ITaskListRepository) *TaskManagerService {
	return &TaskManagerService{repo: repo}
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
func (s *TaskManagerService) CreateTaskList(dto dtos.CreateTaskListDTO) (*task_list.TaskListEntity, error) {
	taskList := task_list.NewTaskListEntity(dto.Title)

	if err := s.repo.Add(taskList); err != nil {
		return nil, err
	}

	if err := s.repo.Flush(); err != nil {
		return nil, err
	}

	return taskList, nil
}
