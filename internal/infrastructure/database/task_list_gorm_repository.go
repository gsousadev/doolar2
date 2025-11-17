package database

import (
	"errors"

	"github.com/google/uuid"
	"github.com/gsousadev/doolar2/internal/domain/entity"
	"github.com/gsousadev/doolar2/internal/domain/entity/task_list"
	"github.com/gsousadev/doolar2/internal/domain/repository"
	"gorm.io/gorm"
)

// TaskListGormRepository implementa repository.TaskListRepository com Unit of Work
type TaskListGormRepository struct {
	db                *gorm.DB
	pendingOperations []func(*gorm.DB) error
}

// NewTaskListGormRepository cria um novo repositório GORM
func NewTaskListGormRepository(db *gorm.DB) repository.TaskListRepository {
	return &TaskListGormRepository{
		db:                db,
		pendingOperations: make([]func(*gorm.DB) error, 0),
	}
}

// taskListGormModel é o modelo GORM (Data Mapper)
type taskListGormModel struct {
	ID    string `gorm:"primaryKey;type:uuid"`
	Title string `gorm:"type:varchar(255);not null"`
	Tasks string `gorm:"type:jsonb;default:'[]'"` // JSON array de task IDs
}

// TaskListGormModel é exportado para auto-migration
type TaskListGormModel = taskListGormModel

// TableName especifica o nome da tabela
func (taskListGormModel) TableName() string {
	return "task_lists"
}

// domainToModel converte domain entity → database model
func domainToGormModel(entity *task_list.TaskListEntity) *taskListGormModel {
	taskIDs := make([]string, len(entity.Tasks))
	for i, task := range entity.Tasks {
		taskIDs[i] = task.GetID().String()
	}

	// Serializar tasks manualmente se necessário
	// Por simplicidade, usando string vazio se não houver tasks
	tasksJSON := "[]"
	if len(taskIDs) > 0 {
		// Você pode usar json.Marshal aqui
		tasksJSON = "[" + taskIDs[0] + "]" // simplificado
	}

	return &taskListGormModel{
		ID:    entity.ID.String(),
		Title: entity.Title,
		Tasks: tasksJSON,
	}
}

// modelToDomain converte database model → domain entity
func gormModelToDomain(model *taskListGormModel) (*task_list.TaskListEntity, error) {
	entityID, err := uuid.Parse(model.ID)
	if err != nil {
		return nil, err
	}

	return &task_list.TaskListEntity{
		Entity: &entity.Entity{ID: entityID},
		Title:  model.Title,
		Tasks:  []task_list.ITask{},
	}, nil
}

// Add adiciona operação à pilha de execução
func (r *TaskListGormRepository) Add(t *task_list.TaskListEntity) error {
	model := domainToGormModel(t)

	// Adiciona operação à pilha (não executa ainda!)
	operation := func(tx *gorm.DB) error {
		return tx.Create(model).Error
	}

	r.pendingOperations = append(r.pendingOperations, operation)
	return nil
}

// FindByID busca imediatamente (não usa pilha)
func (r *TaskListGormRepository) FindByID(id string) (*task_list.TaskListEntity, error) {
	var model taskListGormModel

	result := r.db.First(&model, "id = ?", id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("task list not found")
		}
		return nil, result.Error
	}

	return gormModelToDomain(&model)
}

// Remove adiciona operação de remoção à pilha
func (r *TaskListGormRepository) Remove(id string) error {
	operation := func(tx *gorm.DB) error {
		result := tx.Delete(&taskListGormModel{}, "id = ?", id)
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return errors.New("task list not found")
		}
		return nil
	}

	r.pendingOperations = append(r.pendingOperations, operation)
	return nil
}

// Update adiciona operação de update à pilha
func (r *TaskListGormRepository) Update(t *task_list.TaskListEntity) error {
	model := domainToGormModel(t)

	operation := func(tx *gorm.DB) error {
		result := tx.Model(&taskListGormModel{}).
			Where("id = ?", model.ID).
			Updates(map[string]interface{}{
				"title": model.Title,
				"tasks": model.Tasks,
			})

		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return errors.New("task list not found")
		}
		return nil
	}

	r.pendingOperations = append(r.pendingOperations, operation)
	return nil
}

// Flush executa todas as operações pendentes em uma transação
func (r *TaskListGormRepository) Flush() error {
	if len(r.pendingOperations) == 0 {
		return nil // Nada para fazer
	}

	// Inicia transação
	return r.db.Transaction(func(tx *gorm.DB) error {
		// Executa todas as operações na ordem
		for _, operation := range r.pendingOperations {
			if err := operation(tx); err != nil {
				return err // Rollback automático
			}
		}

		// Limpa a pilha após sucesso
		r.pendingOperations = make([]func(*gorm.DB) error, 0)
		return nil // Commit automático
	})
}

// FindAll busca todas as task lists (operação imediata)
func (r *TaskListGormRepository) FindAll() ([]*task_list.TaskListEntity, error) {
	var models []taskListGormModel

	result := r.db.Find(&models)
	if result.Error != nil {
		return nil, result.Error
	}

	entities := make([]*task_list.TaskListEntity, 0, len(models))
	for _, model := range models {
		entity, err := gormModelToDomain(&model)
		if err != nil {
			return nil, err
		}
		entities = append(entities, entity)
	}

	return entities, nil
}

// Clear limpa a pilha de operações pendentes (útil para testes)
func (r *TaskListGormRepository) Clear() {
	r.pendingOperations = make([]func(*gorm.DB) error, 0)
}

// PendingCount retorna o número de operações pendentes
func (r *TaskListGormRepository) PendingCount() int {
	return len(r.pendingOperations)
}
