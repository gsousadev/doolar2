package task_list

import (
	"errors"
	"slices"

	"github.com/gsousadev/doolar2/internal/domain/entity"
)

type Status string

const (
	StatusPending    Status = "pending"
	StatusInProgress Status = "in_progress"
	StatusCompleted  Status = "completed"
	StatusCancelled  Status = "cancelled"
)

var finalStatuses = []Status{
	StatusCompleted,
	StatusCancelled,
}

var ErrorChangingFinalStatus = errors.New("cannot change task status in a final state")

type ITask interface {
	entity.IEntity
	ChangeStatus(newStatus Status) error
	GetStatus() Status
}

type TaskEntity struct {
	*entity.Entity
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      Status `json:"status"`
}

func NewTaskEntity(title, description string) *TaskEntity {
	return &TaskEntity{
		Entity:      entity.NewEntity(),
		Title:       title,
		Description: description,
		Status:      StatusPending,
	}
}

func (t *TaskEntity) ChangeStatus(newStatus Status) error {

	if slices.Contains(finalStatuses, t.Status) {
		return ErrorChangingFinalStatus
	}

	t.Status = newStatus

	return nil
}

func (t *TaskEntity) GetStatus() Status {
	return t.Status
}
