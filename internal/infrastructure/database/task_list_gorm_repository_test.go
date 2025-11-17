package database

import (
	"testing"

	"github.com/gsousadev/doolar2/internal/domain/entity/task_list"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func setupGormTestDB(t *testing.T) *gorm.DB {
	cfg := Config{
		Host:     "localhost",
		Port:     5432,
		User:     "postgres",
		Password: "postgres",
		DBName:   "doolar_test",
		SSLMode:  "disable",
	}

	db, err := NewGormConnection(cfg)
	if err != nil {
		t.Skip("Database not available for integration tests")
	}

	// Auto migrate
	db.AutoMigrate(&taskListGormModel{})

	// Limpar tabela
	db.Exec("TRUNCATE TABLE task_lists CASCADE")

	return db
}

func TestGormRepository_UnitOfWork_Flush(t *testing.T) {
	db := setupGormTestDB(t)
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}()

	repo := NewTaskListGormRepository(db).(*TaskListGormRepository)

	// Arrange
	taskList1 := task_list.NewTaskListEntity("Lista 1")
	taskList2 := task_list.NewTaskListEntity("Lista 2")
	taskList3 := task_list.NewTaskListEntity("Lista 3")

	// Act - Adiciona operações à pilha (NÃO executa ainda)
	err := repo.Add(taskList1)
	require.NoError(t, err)

	err = repo.Add(taskList2)
	require.NoError(t, err)

	err = repo.Add(taskList3)
	require.NoError(t, err)

	// Assert - Verifica que operações estão pendentes
	assert.Equal(t, 3, repo.PendingCount(), "Deve ter 3 operações pendentes")

	// Assert - Nada foi persistido ainda
	all, _ := repo.FindAll()
	assert.Len(t, all, 0, "Nada deve estar no banco antes do Flush")

	// Act - Executa todas as operações
	err = repo.Flush()
	require.NoError(t, err)

	// Assert - Pilha foi limpa
	assert.Equal(t, 0, repo.PendingCount(), "Pilha deve estar vazia após Flush")

	// Assert - Tudo foi persistido
	all, err = repo.FindAll()
	require.NoError(t, err)
	assert.Len(t, all, 3, "Deve ter 3 task lists no banco")
}

func TestGormRepository_UnitOfWork_Rollback(t *testing.T) {
	db := setupGormTestDB(t)
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}()

	repo := NewTaskListGormRepository(db).(*TaskListGormRepository)

	// Arrange
	taskList1 := task_list.NewTaskListEntity("Lista 1")
	taskList2 := task_list.NewTaskListEntity("Lista 2")

	// Persiste a primeira
	repo.Add(taskList1)
	repo.Flush()

	// Act - Tenta remover ID inválido (causará erro)
	repo.Remove("invalid-uuid-that-does-not-exist")
	repo.Add(taskList2)

	// Flush deve falhar e fazer rollback
	err := repo.Flush()

	// Assert - Erro deve ocorrer
	assert.Error(t, err)

	// Assert - taskList2 NÃO deve ter sido inserida (rollback)
	all, _ := repo.FindAll()
	assert.Len(t, all, 1, "Apenas 1 task list deve existir (rollback funcionou)")
}

func TestGormRepository_MixedOperations(t *testing.T) {
	db := setupGormTestDB(t)
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}()

	repo := NewTaskListGormRepository(db).(*TaskListGormRepository)

	// Arrange
	taskList1 := task_list.NewTaskListEntity("Original")
	taskList2 := task_list.NewTaskListEntity("To Delete")

	// Act - Adiciona e executa
	repo.Add(taskList1)
	repo.Add(taskList2)
	repo.Flush()

	// Act - Update e Delete na pilha
	taskList1.Title = "Updated"
	repo.Update(taskList1)
	repo.Remove(taskList2.ID.String())

	// Assert - Ainda não executou
	assert.Equal(t, 2, repo.PendingCount())

	// Act - Flush
	err := repo.Flush()
	require.NoError(t, err)

	// Assert - Verifica resultado
	all, err := repo.FindAll()
	require.NoError(t, err)
	assert.Len(t, all, 1, "Apenas 1 deve existir")
	assert.Equal(t, "Updated", all[0].Title, "Título deve estar atualizado")
}

func TestGormRepository_Clear(t *testing.T) {
	db := setupGormTestDB(t)
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}()

	repo := NewTaskListGormRepository(db).(*TaskListGormRepository)

	// Arrange
	taskList := task_list.NewTaskListEntity("Lista")
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
