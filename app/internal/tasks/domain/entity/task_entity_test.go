package entity

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestNewTaskEntity(t *testing.T) {
	task := NewTaskEntity("Test Task", "This is a test task")
	assert.IsType(t, task, &TaskEntity{})
	assert.IsType(t, task.ID, uuid.UUID{})
	assert.Equal(t, task.Status, StatusPending, "Expected new task to have status 'pending'")
}

func Test_whenChangeStatusFromACompletedTask_generateError(t *testing.T) {
	task := NewTaskEntity("Test Task", "This is a test task")
	err := task.ChangeStatus(StatusCompleted)
	assert.Equal(t, nil, err, "Expected no error when changing status to completed")
	err = task.ChangeStatus(StatusCompleted)
	assert.Equal(t, ErrorChangingFinalStatus, err, "Expected error when changing status of a completed task")
}

func Test_whenChangeStatusFromAPendingTaskToInProgress_shouldChangeSuccessfully(t *testing.T) {
	task := NewTaskEntity("Test Task", "This is a test task")
	err := task.ChangeStatus(StatusInProgress)
	assert.Equal(t, nil, err, "Expected no error when changing status from pending to in_progress")
	assert.Equal(t, StatusInProgress, task.GetStatus(), "Expected task status to be in_progress")
}
