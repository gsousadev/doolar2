package repository

import task_list "github.com/gsousadev/doolar2/internal/tasks/domain/entity"

type ITaskListRepository interface {
	Add(t *task_list.TaskListEntity) error
	FindByID(id string) (*task_list.TaskListEntity, error)
	Update(t *task_list.TaskListEntity) error
	Remove(id string) error
	Flush() error
}
