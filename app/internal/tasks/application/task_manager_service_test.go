package application

import (
	"errors"
	"testing"

	task_list "github.com/gsousadev/doolar2/internal/tasks/domain/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockTaskListRepository é um mock do repositório para testes
type MockTaskListRepository struct {
	mock.Mock
}

func (m *MockTaskListRepository) Add(t *task_list.TaskListEntity) error {
	args := m.Called(t)
	return args.Error(0)
}

func (m *MockTaskListRepository) FindByID(id string) (*task_list.TaskListEntity, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*task_list.TaskListEntity), args.Error(1)
}

func (m *MockTaskListRepository) Update(t *task_list.TaskListEntity) error {
	args := m.Called(t)
	return args.Error(0)
}

func (m *MockTaskListRepository) Remove(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockTaskListRepository) Flush() error {
	args := m.Called()
	return args.Error(0)
}

func TestCreateTaskList_Success(t *testing.T) {
	// Arrange
	mockRepo := new(MockTaskListRepository)
	service := NewTaskManagerService(mockRepo)

	dto := CreateTaskListDTO{
		Title: "Test List",
	}

	mockRepo.On("Add", mock.AnythingOfType("*task_list.TaskListEntity")).Return(nil)
	mockRepo.On("Flush").Return(nil)

	// Act
	result, err := service.CreateTaskList(dto)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "Test List", result.Title)
	assert.NotEmpty(t, result.ID)
	mockRepo.AssertExpectations(t)
}

func TestCreateTaskList_AddError(t *testing.T) {
	// Arrange
	mockRepo := new(MockTaskListRepository)
	service := NewTaskManagerService(mockRepo)

	dto := CreateTaskListDTO{
		Title: "Test List",
	}

	expectedError := errors.New("database error")
	mockRepo.On("Add", mock.AnythingOfType("*task_list.TaskListEntity")).Return(expectedError)

	// Act
	result, err := service.CreateTaskList(dto)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, expectedError, err)
	mockRepo.AssertExpectations(t)
}

func TestCreateTaskList_FlushError(t *testing.T) {
	// Arrange
	mockRepo := new(MockTaskListRepository)
	service := NewTaskManagerService(mockRepo)

	dto := CreateTaskListDTO{
		Title: "Test List",
	}

	expectedError := errors.New("flush error")
	mockRepo.On("Add", mock.AnythingOfType("*task_list.TaskListEntity")).Return(nil)
	mockRepo.On("Flush").Return(expectedError)

	// Act
	result, err := service.CreateTaskList(dto)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, expectedError, err)
	mockRepo.AssertExpectations(t)
}
