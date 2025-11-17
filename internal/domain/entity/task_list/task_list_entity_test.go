package tasklist

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewTaskList(t *testing.T) {
	// Arrange
	title := "My Task List"

	// Act
	taskList := NewTaskList(title)

	// Assert
	assert.Equal(t, title, taskList.Title, "Expected Title to match")
	assert.NotNil(t, taskList.Tasks, "Expected Tasks to be initialized")
	assert.Empty(t, taskList.Tasks, "Expected Tasks to be empty")
}

func TestAddTask(t *testing.T) {
	// Arrange
	taskList := NewTaskList("Test List")
	task := NewTaskEntity("Task 1", "Description 1")

	// Act
	taskList.AddTask(task)

	// Assert
	assert.Len(t, taskList.Tasks, 1, "Expected 1 task in the list")
	assert.Equal(t, task, taskList.Tasks[0], "Expected added task to match")
}

func TestAddMultipleTasks(t *testing.T) {
	// Arrange
	taskList := NewTaskList("Test List")
	task1 := NewTaskEntity("Task 1", "Description 1")
	task2 := NewTaskEntity("Task 2", "Description 2")
	task3 := NewTaskEntity("Task 3", "Description 3")

	// Act
	taskList.AddTask(task1)
	taskList.AddTask(task2)
	taskList.AddTask(task3)

	// Assert
	assert.Len(t, taskList.Tasks, 3, "Expected 3 tasks in the list")
	assert.Equal(t, task1, taskList.Tasks[0], "Expected first task to match task1")
	assert.Equal(t, task2, taskList.Tasks[1], "Expected second task to match task2")
	assert.Equal(t, task3, taskList.Tasks[2], "Expected third task to match task3")
}

func TestMultipleTypesOfTasks(t *testing.T) {
	// Arrange
	taskList := NewTaskList("Mixed Task List")
	task1 := NewTaskEntity("Simple Task", "A simple task")
	task2 := NewTimedTaskEntity("Timed Task", "A task with time limit", time.Now(), time.Now().Add(2*time.Hour))

	// Act
	taskList.AddTask(task1)
	taskList.AddTask(task2)

	// Assert
	assert.Len(t, taskList.Tasks, 2, "Expected 2 tasks in the list")
	assert.Equal(t, task1, taskList.Tasks[0], "Expected first task to match simple task")
	assert.IsType(t, &TaskEntity{}, taskList.Tasks[0], "Expected first task to be of type TaskEntity")
	assert.Equal(t, task2, taskList.Tasks[1], "Expected second task to match timed task")
	assert.IsType(t, &TimedTaskEntity{}, taskList.Tasks[1], "Expected second task to be of type TimedTaskEntity")
}
