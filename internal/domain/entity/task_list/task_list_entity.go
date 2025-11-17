package tasklist

import "github.com/gsousadev/doolar-golang/internal/domain/entity"

type TaskList struct {
	*entity.Entity
	Title string
	Tasks []ITask
}

func NewTaskList(title string) *TaskList {
	return &TaskList{
		Entity: entity.NewEntity(),
		Title:  title,
		Tasks:  []ITask{},
	}
}

func (tl *TaskList) AddTask(task ITask) {
	tl.Tasks = append(tl.Tasks, task)
}
