package entity

import "github.com/gsousadev/doolar2/internal/shared/domain/entity"

type TaskListEntity struct {
	*entity.Entity
	Title string
	Tasks []ITask
}

func NewTaskListEntity(title string) *TaskListEntity {
	return &TaskListEntity{
		Entity: entity.NewEntity(),
		Title:  title,
		Tasks:  []ITask{},
	}
}

func (tl *TaskListEntity) AddTask(task ITask) {
	tl.Tasks = append(tl.Tasks, task)
}
