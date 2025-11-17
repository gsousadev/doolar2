package task_list

import "github.com/gsousadev/doolar2/internal/domain/entity"

type HomeTask struct {
	*TaskEntity
	Room entity.Room
}

func NewHomeTask(title, description string, room entity.Room) HomeTask {
	return HomeTask{
		TaskEntity: NewTaskEntity(title, description),
		Room:       room,
	}
}
