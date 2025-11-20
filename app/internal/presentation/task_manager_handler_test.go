package presentation

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gsousadev/doolar2/internal/application"
	"github.com/gsousadev/doolar2/internal/domain/entity/task_list"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockTaskManager é um mock da interface TaskManager para testes
type MockTaskManager struct {
	mock.Mock
}

func (m *MockTaskManager) CreateTaskList(dto application.CreateTaskListDTO) (*task_list.TaskListEntity, error) {
	args := m.Called(dto)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*task_list.TaskListEntity), args.Error(1)
}

func (m *MockTaskManager) GetTaskList(id string) (*task_list.TaskListEntity, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*task_list.TaskListEntity), args.Error(1)
}

func (m *MockTaskManager) AddTaskToList(listID string, dto application.CreateTaskDTO) (*task_list.TaskListEntity, error) {
	args := m.Called(listID, dto)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*task_list.TaskListEntity), args.Error(1)
}

func (m *MockTaskManager) GetPendingTasks(listID string) ([]task_list.ITask, error) {
	args := m.Called(listID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]task_list.ITask), args.Error(1)
}

func (m *MockTaskManager) GetTasksByStatus(listID string, status string) ([]task_list.ITask, error) {
	args := m.Called(listID, status)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]task_list.ITask), args.Error(1)
}

func (m *MockTaskManager) UpdateTaskStatus(listID, taskID string, newStatus string) error {
	args := m.Called(listID, taskID, newStatus)
	return args.Error(0)
}

func (m *MockTaskManager) DeleteTaskList(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockTaskManager) GetTaskListForStats(listID string) (*task_list.TaskListEntity, error) {
	args := m.Called(listID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*task_list.TaskListEntity), args.Error(1)
}

func TestCreateTaskList_Success(t *testing.T) {
	// Arrange
	mockService := new(MockTaskManager)
	handler := NewTaskManagerHandler(mockService)

	taskList := task_list.NewTaskListEntity("Test List")
	mockService.On("CreateTaskList", application.CreateTaskListDTO{Title: "Test List"}).Return(taskList, nil)

	reqBody := CreateTaskListRequest{Title: "Test List"}
	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/task-lists", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Act
	handler.CreateTaskList(w, req)

	// Assert
	assert.Equal(t, http.StatusCreated, w.Code)

	var response SuccessResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Task list created successfully", response.Message)
	assert.NotNil(t, response.Data)

	mockService.AssertExpectations(t)
}

func TestCreateTaskList_InvalidJSON(t *testing.T) {
	// Arrange
	mockService := new(MockTaskManager)
	handler := NewTaskManagerHandler(mockService)

	req := httptest.NewRequest(http.MethodPost, "/task-lists", bytes.NewBufferString("invalid json"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Act
	handler.CreateTaskList(w, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Invalid request body", response.Message)
}

func TestCreateTaskList_EmptyTitle(t *testing.T) {
	// Arrange
	mockService := new(MockTaskManager)
	handler := NewTaskManagerHandler(mockService)

	reqBody := CreateTaskListRequest{Title: ""}
	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/task-lists", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Act
	handler.CreateTaskList(w, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Title is required", response.Message)
}

func TestCreateTaskList_MethodNotAllowed(t *testing.T) {
	// Arrange
	mockService := new(MockTaskManager)
	handler := NewTaskManagerHandler(mockService)

	req := httptest.NewRequest(http.MethodGet, "/task-lists", nil)
	w := httptest.NewRecorder()

	// Act
	handler.CreateTaskList(w, req)

	// Assert
	assert.Equal(t, http.StatusMethodNotAllowed, w.Code)
}

func TestGetTaskList_Success(t *testing.T) {
	// Arrange
	mockService := new(MockTaskManager)
	handler := NewTaskManagerHandler(mockService)

	taskList := task_list.NewTaskListEntity("Test List")
	task := task_list.NewTaskEntity("Test Task", "Description")
	taskList.AddTask(task)

	mockService.On("GetTaskList", taskList.ID.String()).Return(taskList, nil)

	req := httptest.NewRequest(http.MethodGet, "/task-lists/"+taskList.ID.String(), nil)
	w := httptest.NewRecorder()

	// Act
	handler.GetTaskList(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	var response SuccessResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Task list retrieved successfully", response.Message)

	mockService.AssertExpectations(t)
}

func TestGetTaskList_NotFound(t *testing.T) {
	// Arrange
	mockService := new(MockTaskManager)
	handler := NewTaskManagerHandler(mockService)

	mockService.On("GetTaskList", "invalid-id").Return(nil, application.ErrTaskListNotFound)

	req := httptest.NewRequest(http.MethodGet, "/task-lists/invalid-id", nil)
	w := httptest.NewRecorder()

	// Act
	handler.GetTaskList(w, req)

	// Assert
	assert.Equal(t, http.StatusNotFound, w.Code)

	var response ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Task list not found", response.Message)

	mockService.AssertExpectations(t)
}

func TestAddTaskToList_Success(t *testing.T) {
	// Arrange
	mockService := new(MockTaskManager)
	handler := NewTaskManagerHandler(mockService)

	taskList := task_list.NewTaskListEntity("Test List")
	task := task_list.NewTaskEntity("New Task", "Description")
	taskList.AddTask(task)

	mockService.On("AddTaskToList", taskList.ID.String(), application.CreateTaskDTO{
		Title:       "New Task",
		Description: "Description",
	}).Return(taskList, nil)

	reqBody := CreateTaskRequest{Title: "New Task", Description: "Description"}
	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/task-lists/"+taskList.ID.String()+"/tasks", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Act
	handler.AddTaskToList(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	var response SuccessResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Task added successfully", response.Message)

	mockService.AssertExpectations(t)
}

func TestAddTaskToList_EmptyTitle(t *testing.T) {
	// Arrange
	mockService := new(MockTaskManager)
	handler := NewTaskManagerHandler(mockService)

	reqBody := CreateTaskRequest{Title: "", Description: "Description"}
	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/task-lists/some-id/tasks", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Act
	handler.AddTaskToList(w, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Title is required", response.Message)
}

func TestGetPendingTasks_Success(t *testing.T) {
	// Arrange
	mockService := new(MockTaskManager)
	handler := NewTaskManagerHandler(mockService)

	task1 := task_list.NewTaskEntity("Pending Task 1", "Description")
	task2 := task_list.NewTaskEntity("Pending Task 2", "Description")
	tasks := []task_list.ITask{task1, task2}

	listID := "test-list-id"
	mockService.On("GetPendingTasks", listID).Return(tasks, nil)

	req := httptest.NewRequest(http.MethodGet, "/task-lists/"+listID+"/tasks/pending", nil)
	w := httptest.NewRecorder()

	// Act
	handler.GetPendingTasks(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	var response SuccessResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Pending tasks retrieved successfully", response.Message)

	mockService.AssertExpectations(t)
}

func TestUpdateTaskStatus_Success(t *testing.T) {
	// Arrange
	mockService := new(MockTaskManager)
	handler := NewTaskManagerHandler(mockService)

	listID := "list-id"
	taskID := "task-id"

	mockService.On("UpdateTaskStatus", listID, taskID, "in_progress").Return(nil)

	reqBody := UpdateTaskStatusRequest{Status: "in_progress"}
	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPatch, "/task-lists/"+listID+"/tasks/"+taskID+"/status", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Act
	handler.UpdateTaskStatus(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	var response SuccessResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Task status updated successfully", response.Message)

	mockService.AssertExpectations(t)
}

func TestUpdateTaskStatus_EmptyStatus(t *testing.T) {
	// Arrange
	mockService := new(MockTaskManager)
	handler := NewTaskManagerHandler(mockService)

	reqBody := UpdateTaskStatusRequest{Status: ""}
	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPatch, "/task-lists/list-id/tasks/task-id/status", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Act
	handler.UpdateTaskStatus(w, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Status is required", response.Message)
}

func TestUpdateTaskStatus_TaskNotFound(t *testing.T) {
	// Arrange
	mockService := new(MockTaskManager)
	handler := NewTaskManagerHandler(mockService)

	listID := "list-id"
	taskID := "invalid-task-id"

	mockService.On("UpdateTaskStatus", listID, taskID, "completed").Return(application.ErrTaskNotFound)

	reqBody := UpdateTaskStatusRequest{Status: "completed"}
	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPatch, "/task-lists/"+listID+"/tasks/"+taskID+"/status", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Act
	handler.UpdateTaskStatus(w, req)

	// Assert
	assert.Equal(t, http.StatusNotFound, w.Code)

	var response ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Task not found", response.Message)

	mockService.AssertExpectations(t)
}

func TestDeleteTaskList_Success(t *testing.T) {
	// Arrange
	mockService := new(MockTaskManager)
	handler := NewTaskManagerHandler(mockService)

	listID := "test-list-id"
	mockService.On("DeleteTaskList", listID).Return(nil)

	req := httptest.NewRequest(http.MethodDelete, "/task-lists/"+listID, nil)
	w := httptest.NewRecorder()

	// Act
	handler.DeleteTaskList(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	var response SuccessResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Task list deleted successfully", response.Message)

	mockService.AssertExpectations(t)
}

func TestDeleteTaskList_Error(t *testing.T) {
	// Arrange
	mockService := new(MockTaskManager)
	handler := NewTaskManagerHandler(mockService)

	listID := "test-list-id"
	expectedError := errors.New("delete error")
	mockService.On("DeleteTaskList", listID).Return(expectedError)

	req := httptest.NewRequest(http.MethodDelete, "/task-lists/"+listID, nil)
	w := httptest.NewRecorder()

	// Act
	handler.DeleteTaskList(w, req)

	// Assert
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var response ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, expectedError.Error(), response.Message)

	mockService.AssertExpectations(t)
}

func TestGetStatistics_Success(t *testing.T) {
	// Arrange
	mockService := new(MockTaskManager)
	handler := NewTaskManagerHandler(mockService)

	taskList := task_list.NewTaskListEntity("Test List")
	task1 := task_list.NewTaskEntity("Task 1", "Description")
	task2 := task_list.NewTaskEntity("Task 2", "Description")
	task2.ChangeStatus(task_list.StatusCompleted)
	task3 := task_list.NewTaskEntity("Task 3", "Description")
	task3.ChangeStatus(task_list.StatusInProgress)

	taskList.AddTask(task1)
	taskList.AddTask(task2)
	taskList.AddTask(task3)

	mockService.On("GetTaskListForStats", taskList.ID.String()).Return(taskList, nil)

	req := httptest.NewRequest(http.MethodGet, "/task-lists/"+taskList.ID.String()+"/statistics", nil)
	w := httptest.NewRecorder()

	// Act
	handler.GetStatistics(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	var response SuccessResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Statistics retrieved successfully", response.Message)

	// Verifica as estatísticas - Data é um map[string]interface{}
	statsMap, ok := response.Data.(map[string]interface{})
	assert.True(t, ok)
	assert.Equal(t, float64(3), statsMap["total"])
	assert.Equal(t, float64(1), statsMap["pending"])
	assert.Equal(t, float64(1), statsMap["in_progress"])
	assert.Equal(t, float64(1), statsMap["completed"])
	assert.Equal(t, float64(0), statsMap["cancelled"])

	mockService.AssertExpectations(t)
}

func TestGetStatistics_NotFound(t *testing.T) {
	// Arrange
	mockService := new(MockTaskManager)
	handler := NewTaskManagerHandler(mockService)

	mockService.On("GetTaskListForStats", "invalid-id").Return(nil, application.ErrTaskListNotFound)

	req := httptest.NewRequest(http.MethodGet, "/task-lists/invalid-id/statistics", nil)
	w := httptest.NewRecorder()

	// Act
	handler.GetStatistics(w, req)

	// Assert
	assert.Equal(t, http.StatusNotFound, w.Code)

	var response ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Task list not found", response.Message)

	mockService.AssertExpectations(t)
}
