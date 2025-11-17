package task_list

import (
	"time"

	"github.com/gsousadev/doolar2/internal/domain/entity"
)

type TimedHomeTask struct {
	*TimedTaskEntity
	Room *entity.Room
}

func NewTimedHomeTask(room *entity.Room, title, description string, startDate, endDate time.Time) *TimedHomeTask {
	return &TimedHomeTask{
		TimedTaskEntity: NewTimedTaskEntity(title, description, startDate, endDate),
		Room:            room,
	}
}
