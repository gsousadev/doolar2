package application

import "github.com/gsousadev/doolar2/internal/domain/entity/task_list"

// TaskManager define o contrato para gerenciamento de listas de tarefas
// Esta interface permite que a camada de apresentação não dependa diretamente
// da implementação concreta do serviço
type TaskManager interface {
	// CreateTaskList cria uma nova lista de tarefas
	CreateTaskList(dto CreateTaskListDTO) (*task_list.TaskListEntity, error)

	// GetTaskList busca uma lista de tarefas por ID
	GetTaskList(id string) (*task_list.TaskListEntity, error)

	// AddTaskToList adiciona uma nova task a uma lista existente
	AddTaskToList(listID string, dto CreateTaskDTO) (*task_list.TaskListEntity, error)

	// GetPendingTasks retorna apenas as tasks pendentes de uma lista
	GetPendingTasks(listID string) ([]task_list.ITask, error)

	// GetTasksByStatus retorna tasks filtradas por status
	GetTasksByStatus(listID string, status string) ([]task_list.ITask, error)

	// UpdateTaskStatus atualiza o status de uma task
	UpdateTaskStatus(listID, taskID string, newStatus string) error

	// DeleteTaskList remove uma lista de tarefas
	DeleteTaskList(id string) error

	// GetTaskListForStats retorna a lista completa para cálculo de estatísticas
	GetTaskListForStats(listID string) (*task_list.TaskListEntity, error)
}
