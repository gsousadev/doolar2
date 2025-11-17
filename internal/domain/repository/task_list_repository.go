package repository

import "github.com/gsousadev/doolar2/internal/domain/entity/task_list"

type TaskListRepository interface {
	Add(t *task_list.TaskListEntity) error
	FindByID(id string) (*task_list.TaskListEntity, error)
	Remove(id string) error
	Flush() error
}
