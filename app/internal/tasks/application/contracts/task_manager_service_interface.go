package contracts

import (
	"github.com/gsousadev/doolar2/internal/tasks/application/dtos"
	"github.com/gsousadev/doolar2/internal/tasks/domain/entity"
)

type ITaskManagerService interface {
	CreateTaskList(dto dtos.CreateTaskListDTO) (*entity.TaskListEntity, error)
}
