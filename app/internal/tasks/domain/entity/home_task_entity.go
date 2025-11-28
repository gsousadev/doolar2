package task_list

import house_entity "github.com/gsousadev/doolar2/internal/house/domain/entity"

type HomeTask struct {
	*TaskEntity
	Room house_entity.Room
}

func NewHomeTask(title, description string, room house_entity.Room) HomeTask {
	return HomeTask{
		TaskEntity: NewTaskEntity(title, description),
		Room:       room,
	}
}
