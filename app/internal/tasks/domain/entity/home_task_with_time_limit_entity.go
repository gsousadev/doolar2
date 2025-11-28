package task_list

import (
	"time"

	house_entity "github.com/gsousadev/doolar2/internal/house/domain/entity"
)

type TimedHomeTask struct {
	*TimedTaskEntity
	Room *house_entity.Room
}

func NewTimedHomeTask(room *house_entity.Room, title, description string, startDate, endDate time.Time) *TimedHomeTask {
	return &TimedHomeTask{
		TimedTaskEntity: NewTimedTaskEntity(title, description, startDate, endDate),
		Room:            room,
	}
}
