package mongo

import (
	"context"
	"testing"
	"time"

	database "github.com/gsousadev/doolar2/internal/shared/infrastructure/database"
	task_list "github.com/gsousadev/doolar2/internal/tasks/domain/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
)

func setupMongoTestDB(t *testing.T) *TaskListMongoRepository {
	cfg := database.MongoConfig{
		URI:      "mongodb://localhost:27017",
		Database: "doolar_test",
		Timeout:  10 * time.Second,
	}

	client, err := database.NewMongoConnection(cfg)
	if err != nil {
		t.Skip("MongoDB not available for integration tests")
	}

	// Limpar coleção antes dos testes
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := client.Database(cfg.Database).Collection("task_lists")
	collection.DeleteMany(ctx, bson.M{})

	repo := NewTaskListMongoRepository(client, cfg.Database).(*TaskListMongoRepository)
	return repo
}

func TestMongoRepository_UnitOfWork_Flush(t *testing.T) {
	repo := setupMongoTestDB(t)
	defer func() {
		repo.client.Disconnect(context.Background())
	}()

	// Arrange
	taskList1 := task_list.NewTaskListEntity("Lista MongoDB 1")
	taskList2 := task_list.NewTaskListEntity("Lista MongoDB 2")
	taskList3 := task_list.NewTaskListEntity("Lista MongoDB 3")

	// Act - Adiciona operações à pilha (NÃO executa ainda)
	err := repo.Add(taskList1)
	require.NoError(t, err)

	err = repo.Add(taskList2)
	require.NoError(t, err)

	err = repo.Add(taskList3)
	require.NoError(t, err)

	// Assert - Verifica que operações estão pendentes
	assert.Equal(t, 3, repo.PendingCount(), "Deve ter 3 operações pendentes")
	assert.Equal(t, []string{"INSERT", "INSERT", "INSERT"}, repo.PendingOperationTypes())

	// Assert - Nada foi persistido ainda
	all, _ := repo.FindAll()
	assert.Len(t, all, 0, "Nada deve estar no banco antes do Flush")

	// Act - Executa todas as operações em transação
	err = repo.Flush()
	require.NoError(t, err)

	// Assert - Pilha foi limpa
	assert.Equal(t, 0, repo.PendingCount(), "Pilha deve estar vazia após Flush")

	// Assert - Tudo foi persistido
	all, err = repo.FindAll()
	require.NoError(t, err)
	assert.Len(t, all, 3, "Deve ter 3 task lists no MongoDB")
}

func TestMongoRepository_UnitOfWork_Rollback(t *testing.T) {
	repo := setupMongoTestDB(t)
	defer func() {
		repo.client.Disconnect(context.Background())
	}()

	// Arrange
	taskList1 := task_list.NewTaskListEntity("Lista Válida")
	taskList2 := task_list.NewTaskListEntity("Lista que será adicionada depois")

	// Persiste a primeira
	repo.Add(taskList1)
	repo.Flush()

	// Act - Tenta remover ID inválido (causará erro) e adicionar outra
	repo.Remove("id-inexistente-que-nao-existe")
	repo.Add(taskList2)

	// Flush deve falhar e fazer rollback
	err := repo.Flush()

	// Assert - Erro deve ocorrer
	assert.Error(t, err)

	// Assert - taskList2 NÃO deve ter sido inserida (rollback funcionou)
	all, _ := repo.FindAll()
	assert.Len(t, all, 1, "Apenas 1 task list deve existir (rollback funcionou)")
	assert.Equal(t, taskList1.ID.String(), all[0].ID.String())
}

func TestMongoRepository_MixedOperations(t *testing.T) {
	repo := setupMongoTestDB(t)
	defer func() {
		repo.client.Disconnect(context.Background())
	}()

	// Arrange
	taskList1 := task_list.NewTaskListEntity("Original MongoDB")
	taskList2 := task_list.NewTaskListEntity("To Delete MongoDB")

	// Act - Adiciona e executa
	repo.Add(taskList1)
	repo.Add(taskList2)
	repo.Flush()

	// Act - Update e Delete na pilha
	taskList1.Title = "Updated MongoDB"
	repo.Update(taskList1)
	repo.Remove(taskList2.ID.String())

	// Assert - Ainda não executou
	assert.Equal(t, 2, repo.PendingCount())
	assert.Equal(t, []string{"UPDATE", "DELETE"}, repo.PendingOperationTypes())

	// Act - Flush
	err := repo.Flush()
	require.NoError(t, err)

	// Assert - Verifica resultado
	all, err := repo.FindAll()
	require.NoError(t, err)
	assert.Len(t, all, 1, "Apenas 1 deve existir")
	assert.Equal(t, "Updated MongoDB", all[0].Title, "Título deve estar atualizado")
}

func TestMongoRepository_FindByID(t *testing.T) {
	repo := setupMongoTestDB(t)
	defer func() {
		repo.client.Disconnect(context.Background())
	}()

	// Arrange
	taskList := task_list.NewTaskListEntity("Find Me MongoDB")
	repo.Add(taskList)
	repo.Flush()

	// Act
	found, err := repo.FindByID(taskList.ID.String())

	// Assert
	require.NoError(t, err)
	assert.Equal(t, taskList.ID, found.ID)
	assert.Equal(t, taskList.Title, found.Title)
}

func TestMongoRepository_FindByID_NotFound(t *testing.T) {
	repo := setupMongoTestDB(t)
	defer func() {
		repo.client.Disconnect(context.Background())
	}()

	// Act
	found, err := repo.FindByID("non-existent-id")

	// Assert
	assert.Error(t, err)
	assert.Nil(t, found)
	assert.Contains(t, err.Error(), "not found")
}

func TestMongoRepository_Clear(t *testing.T) {
	repo := setupMongoTestDB(t)
	defer func() {
		repo.client.Disconnect(context.Background())
	}()

	// Arrange
	taskList := task_list.NewTaskListEntity("Lista MongoDB")
	repo.Add(taskList)

	// Assert - Tem 1 operação pendente
	assert.Equal(t, 1, repo.PendingCount())

	// Act - Limpa a pilha
	repo.Clear()

	// Assert - Pilha vazia
	assert.Equal(t, 0, repo.PendingCount())

	// Assert - Nada persistido
	all, _ := repo.FindAll()
	assert.Len(t, all, 0)
}
