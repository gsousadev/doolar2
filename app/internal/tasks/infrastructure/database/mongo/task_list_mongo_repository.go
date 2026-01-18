package mongo

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/gsousadev/doolar-golang/internal/shared/domain/entity"
	task_list "github.com/gsousadev/doolar-golang/internal/tasks/domain/entity"
	"github.com/gsousadev/doolar-golang/internal/tasks/domain/repository"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// TaskListMongoRepository implementa repository.TaskListRepository com Unit of Work para MongoDB
type TaskListMongoRepository struct {
	client         *mongo.Client
	collection     *mongo.Collection
	operations     []func(mongo.SessionContext) error
	operationTypes []string // Para debugging
}

// NewTaskListMongoRepository cria um novo repositório MongoDB
func NewTaskListMongoRepository(client *mongo.Client, dbName string) repository.ITaskListRepository {
	return &TaskListMongoRepository{
		client:         client,
		collection:     client.Database(dbName).Collection("task_lists"),
		operations:     make([]func(mongo.SessionContext) error, 0),
		operationTypes: make([]string, 0),
	}
}

// taskListMongoModel é o modelo MongoDB (Data Mapper)
type taskListMongoModel struct {
	ID      string   `bson:"_id"`
	Title   string   `bson:"title"`
	TaskIDs []string `bson:"task_ids"`
}

// domainToMongoModel converte domain entity → MongoDB model
func domainToMongoModel(entity *task_list.TaskListEntity) *taskListMongoModel {
	taskIDs := make([]string, len(entity.Tasks))
	for i, task := range entity.Tasks {
		taskIDs[i] = task.GetID().String()
	}

	return &taskListMongoModel{
		ID:      entity.ID.String(),
		Title:   entity.Title,
		TaskIDs: taskIDs,
	}
}

// mongoModelToDomain converte MongoDB model → domain entity
func mongoModelToDomain(model *taskListMongoModel) (*task_list.TaskListEntity, error) {
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
func (r *TaskListMongoRepository) Add(t *task_list.TaskListEntity) error {
	model := domainToMongoModel(t)

	// Adiciona operação à pilha (não executa ainda!)
	operation := func(sessCtx mongo.SessionContext) error {
		_, err := r.collection.InsertOne(sessCtx, model)
		return err
	}

	r.operations = append(r.operations, operation)
	r.operationTypes = append(r.operationTypes, "INSERT")
	return nil
}

// FindByID busca imediatamente (não usa pilha)
func (r *TaskListMongoRepository) FindByID(id string) (*task_list.TaskListEntity, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var model taskListMongoModel
	filter := bson.M{"_id": id}

	err := r.collection.FindOne(ctx, filter).Decode(&model)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, errors.New("task list not found")
		}
		return nil, err
	}

	return mongoModelToDomain(&model)
}

// Remove adiciona operação de remoção à pilha
func (r *TaskListMongoRepository) Remove(id string) error {
	operation := func(sessCtx mongo.SessionContext) error {
		filter := bson.M{"_id": id}
		result, err := r.collection.DeleteOne(sessCtx, filter)
		if err != nil {
			return err
		}
		if result.DeletedCount == 0 {
			return errors.New("task list not found")
		}
		return nil
	}

	r.operations = append(r.operations, operation)
	r.operationTypes = append(r.operationTypes, "DELETE")
	return nil
}

// Update adiciona operação de update à pilha
func (r *TaskListMongoRepository) Update(t *task_list.TaskListEntity) error {
	model := domainToMongoModel(t)

	operation := func(sessCtx mongo.SessionContext) error {
		filter := bson.M{"_id": model.ID}
		update := bson.M{
			"$set": bson.M{
				"title":    model.Title,
				"task_ids": model.TaskIDs,
			},
		}

		result, err := r.collection.UpdateOne(sessCtx, filter, update)
		if err != nil {
			return err
		}
		if result.MatchedCount == 0 {
			return errors.New("task list not found")
		}
		return nil
	}

	r.operations = append(r.operations, operation)
	r.operationTypes = append(r.operationTypes, "UPDATE")
	return nil
}

// Flush executa todas as operações pendentes em uma transação MongoDB
func (r *TaskListMongoRepository) Flush() error {
	if len(r.operations) == 0 {
		return nil // Nada para fazer
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Inicia uma sessão
	session, err := r.client.StartSession()
	if err != nil {

		log.Fatal(err)
		return err
	}
	defer session.EndSession(ctx)

	// Executa todas as operações em uma transação
	_, err = session.WithTransaction(ctx, func(sessCtx mongo.SessionContext) (interface{}, error) {
		for _, operation := range r.operations {
			if err := operation(sessCtx); err != nil {

				log.Fatal(err)
				return nil, err // Rollback automático
			}
		}
		return nil, nil // Commit automático
	})

	if err != nil {
		log.Fatal(err)
		return err
	}

	// Limpa a pilha após sucesso
	r.operations = make([]func(mongo.SessionContext) error, 0)
	r.operationTypes = make([]string, 0)
	return nil
}

// FindAll busca todas as task lists (operação imediata)
func (r *TaskListMongoRepository) FindAll() ([]*task_list.TaskListEntity, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var models []taskListMongoModel
	if err := cursor.All(ctx, &models); err != nil {
		return nil, err
	}

	entities := make([]*task_list.TaskListEntity, 0, len(models))
	for _, model := range models {
		entity, err := mongoModelToDomain(&model)
		if err != nil {
			return nil, err
		}
		entities = append(entities, entity)
	}

	return entities, nil
}

// Clear limpa a pilha de operações pendentes (útil para testes)
func (r *TaskListMongoRepository) Clear() {
	r.operations = make([]func(mongo.SessionContext) error, 0)
	r.operationTypes = make([]string, 0)
}

// PendingCount retorna o número de operações pendentes
func (r *TaskListMongoRepository) PendingCount() int {
	return len(r.operations)
}

// PendingOperationTypes retorna os tipos de operações pendentes
func (r *TaskListMongoRepository) PendingOperationTypes() []string {
	return r.operationTypes
}
