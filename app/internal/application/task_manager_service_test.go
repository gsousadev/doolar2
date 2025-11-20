package application

import (
	"errors"
	"testing"

	"github.com/gsousadev/doolar2/internal/domain/entity/task_list"
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

func TestGetTaskList_Success(t *testing.T) {
	// Arrange
	mockRepo := new(MockTaskListRepository)
	service := NewTaskManagerService(mockRepo)

	expectedList := task_list.NewTaskListEntity("Test List")
	mockRepo.On("FindByID", expectedList.ID.String()).Return(expectedList, nil)

	// Act
	result, err := service.GetTaskList(expectedList.ID.String())

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedList.ID, result.ID)
	assert.Equal(t, "Test List", result.Title)
	mockRepo.AssertExpectations(t)
}

func TestGetTaskList_NotFound(t *testing.T) {
	// Arrange
	mockRepo := new(MockTaskListRepository)
	service := NewTaskManagerService(mockRepo)

	mockRepo.On("FindByID", "invalid-id").Return(nil, errors.New("not found"))

	// Act
	result, err := service.GetTaskList("invalid-id")

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, ErrTaskListNotFound, err)
	mockRepo.AssertExpectations(t)
}

func TestAddTaskToList_Success(t *testing.T) {
	// Arrange
	mockRepo := new(MockTaskListRepository)
	service := NewTaskManagerService(mockRepo)

	taskList := task_list.NewTaskListEntity("Test List")
	taskDTO := CreateTaskDTO{
		Title:       "Test Task",
		Description: "Test Description",
	}

	mockRepo.On("FindByID", taskList.ID.String()).Return(taskList, nil)
	mockRepo.On("Update", mock.AnythingOfType("*task_list.TaskListEntity")).Return(nil)
	mockRepo.On("Flush").Return(nil)

	// Act
	result, err := service.AddTaskToList(taskList.ID.String(), taskDTO)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result.Tasks, 1)
	assert.Equal(t, "Test Task", result.Tasks[0].(*task_list.TaskEntity).Title)
	assert.Equal(t, "Test Description", result.Tasks[0].(*task_list.TaskEntity).Description)
	mockRepo.AssertExpectations(t)
}

func TestAddTaskToList_ListNotFound(t *testing.T) {
	// Arrange
	mockRepo := new(MockTaskListRepository)
	service := NewTaskManagerService(mockRepo)

	taskDTO := CreateTaskDTO{
		Title: "Test Task",
	}

	mockRepo.On("FindByID", "invalid-id").Return(nil, errors.New("not found"))

	// Act
	result, err := service.AddTaskToList("invalid-id", taskDTO)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, ErrTaskListNotFound, err)
	mockRepo.AssertExpectations(t)
}

func TestGetPendingTasks_Success(t *testing.T) {
	// Arrange
	mockRepo := new(MockTaskListRepository)
	service := NewTaskManagerService(mockRepo)

	taskList := task_list.NewTaskListEntity("Test List")
	task1 := task_list.NewTaskEntity("Pending Task", "Description")
	task2 := task_list.NewTaskEntity("Another Task", "Description")
	task2.ChangeStatus(task_list.StatusInProgress)

	taskList.AddTask(task1)
	taskList.AddTask(task2)

	mockRepo.On("FindByID", taskList.ID.String()).Return(taskList, nil)

	// Act
	result, err := service.GetPendingTasks(taskList.ID.String())

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result, 1)
	assert.Equal(t, task_list.StatusPending, result[0].GetStatus())
	mockRepo.AssertExpectations(t)
}

func TestGetTasksByStatus_Success(t *testing.T) {
	// Arrange
	mockRepo := new(MockTaskListRepository)
	service := NewTaskManagerService(mockRepo)

	taskList := task_list.NewTaskListEntity("Test List")
	task1 := task_list.NewTaskEntity("Task 1", "Description")
	task2 := task_list.NewTaskEntity("Task 2", "Description")
	task2.ChangeStatus(task_list.StatusCompleted)

	taskList.AddTask(task1)
	taskList.AddTask(task2)

	mockRepo.On("FindByID", taskList.ID.String()).Return(taskList, nil)

	// Act
	result, err := service.GetTasksByStatus(taskList.ID.String(), string(task_list.StatusCompleted))

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result, 1)
	assert.Equal(t, task_list.StatusCompleted, result[0].GetStatus())
	mockRepo.AssertExpectations(t)
}

func TestGetTasksByStatus_InvalidStatus(t *testing.T) {
	// Arrange
	mockRepo := new(MockTaskListRepository)
	service := NewTaskManagerService(mockRepo)

	taskList := task_list.NewTaskListEntity("Test List")
	mockRepo.On("FindByID", taskList.ID.String()).Return(taskList, nil)

	// Act
	result, err := service.GetTasksByStatus(taskList.ID.String(), "invalid_status")

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, ErrInvalidStatus, err)
	mockRepo.AssertExpectations(t)
}

func TestUpdateTaskStatus_Success(t *testing.T) {
	// Arrange
	mockRepo := new(MockTaskListRepository)
	service := NewTaskManagerService(mockRepo)

	taskList := task_list.NewTaskListEntity("Test List")
	task := task_list.NewTaskEntity("Test Task", "Description")
	taskList.AddTask(task)

	mockRepo.On("FindByID", taskList.ID.String()).Return(taskList, nil)
	mockRepo.On("Update", mock.AnythingOfType("*task_list.TaskListEntity")).Return(nil)
	mockRepo.On("Flush").Return(nil)

	// Act
	err := service.UpdateTaskStatus(taskList.ID.String(), task.GetID().String(), string(task_list.StatusInProgress))

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, task_list.StatusInProgress, task.GetStatus())
	mockRepo.AssertExpectations(t)
}

func TestUpdateTaskStatus_TaskNotFound(t *testing.T) {
	// Arrange
	mockRepo := new(MockTaskListRepository)
	service := NewTaskManagerService(mockRepo)

	taskList := task_list.NewTaskListEntity("Test List")
	mockRepo.On("FindByID", taskList.ID.String()).Return(taskList, nil)

	// Act
	err := service.UpdateTaskStatus(taskList.ID.String(), "invalid-task-id", string(task_list.StatusInProgress))

	// Assert
	assert.Error(t, err)
	assert.Equal(t, ErrTaskNotFound, err)
	mockRepo.AssertExpectations(t)
}

func TestUpdateTaskStatus_InvalidStatusChange(t *testing.T) {
	// Arrange
	mockRepo := new(MockTaskListRepository)
	service := NewTaskManagerService(mockRepo)

	taskList := task_list.NewTaskListEntity("Test List")
	task := task_list.NewTaskEntity("Test Task", "Description")
	task.ChangeStatus(task_list.StatusCompleted) // Muda para completed
	taskList.AddTask(task)

	mockRepo.On("FindByID", taskList.ID.String()).Return(taskList, nil)

	// Act - Tenta mudar de completed para pending (não permitido)
	err := service.UpdateTaskStatus(taskList.ID.String(), task.GetID().String(), string(task_list.StatusPending))

	// Assert
	assert.Error(t, err)
	mockRepo.AssertExpectations(t)
}

func TestDeleteTaskList_Success(t *testing.T) {
	// Arrange
	mockRepo := new(MockTaskListRepository)
	service := NewTaskManagerService(mockRepo)

	listID := "test-id"
	mockRepo.On("Remove", listID).Return(nil)
	mockRepo.On("Flush").Return(nil)

	// Act
	err := service.DeleteTaskList(listID)

	// Assert
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestDeleteTaskList_RemoveError(t *testing.T) {
	// Arrange
	mockRepo := new(MockTaskListRepository)
	service := NewTaskManagerService(mockRepo)

	listID := "test-id"
	expectedError := errors.New("remove error")
	mockRepo.On("Remove", listID).Return(expectedError)

	// Act
	err := service.DeleteTaskList(listID)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	mockRepo.AssertExpectations(t)
}

func TestGetTaskListForStats_Success(t *testing.T) {
	// Arrange
	mockRepo := new(MockTaskListRepository)
	service := NewTaskManagerService(mockRepo)

	taskList := task_list.NewTaskListEntity("Test List")
	task1 := task_list.NewTaskEntity("Task 1", "Description")
	task2 := task_list.NewTaskEntity("Task 2", "Description")
	task2.ChangeStatus(task_list.StatusCompleted)
	task3 := task_list.NewTaskEntity("Task 3", "Description")
	task3.ChangeStatus(task_list.StatusCancelled)

	taskList.AddTask(task1)
	taskList.AddTask(task2)
	taskList.AddTask(task3)

	mockRepo.On("FindByID", taskList.ID.String()).Return(taskList, nil)

	// Act
	result, err := service.GetTaskListForStats(taskList.ID.String())

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result.Tasks, 3)
	mockRepo.AssertExpectations(t)
}
