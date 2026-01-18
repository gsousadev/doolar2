package entity

import (
	"github.com/gsousadev/doolar-golang/internal/shared/domain/entity"
)

type TaskListEntity struct {
	*entity.Entity
	Title       string
	Description string
	Tasks       []ITask
}

func NewTaskListEntity(title string, description string) *TaskListEntity {
	return &TaskListEntity{
		Entity:      entity.NewEntity(),
		Title:       title,
		Description: description,
		Tasks:       []ITask{},
	}
}

func (tl *TaskListEntity) AddTask(task ITask) {
	tl.Tasks = append(tl.Tasks, task)
}
