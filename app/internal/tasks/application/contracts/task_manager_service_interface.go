package contracts

import (
	"github.com/gsousadev/doolar-golang/internal/tasks/application/dtos"
	"github.com/gsousadev/doolar-golang/internal/tasks/domain/entity"
)

type ITaskManagerService interface {
	CreateTaskList(dto dtos.CreateTaskListDTO) (*entity.TaskListEntity, error)
}
