package entity

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_newTimedTaskEntity_generateSuccess(t *testing.T) {
	task := NewTimedTaskEntity("Test Task", "This is a test task", time.Now(), time.Now().Add(2*time.Hour))
	assert.IsType(t, task, &TimedTaskEntity{})
	assert.IsType(t, task.TaskEntity, &TaskEntity{})
	assert.Equal(t, StatusPending, task.GetStatus(), "Expected new task to have status 'pending'")
}

func Test_whenChangeStartDateAfterEndDate_generateError(t *testing.T) {
	task := NewTimedTaskEntity("Test Task", "This is a test task", time.Now(), time.Now().Add(2*time.Hour))
	err := task.changeStartDate(time.Now().Add(3 * time.Hour))
	assert.Equal(t, ErrStartDateAfterEndDate, err, "Expected error when changing start date after end date")
}

func Test_whenChangeStartDateBeforeCurrentDate_generateError(t *testing.T) {
	task := NewTimedTaskEntity("Test Task", "This is a test task", time.Now(), time.Now().Add(2*time.Hour))
	err := task.changeStartDate(time.Now().Add(-3 * time.Hour))
	assert.Equal(t, ErrStartDateBeforeNow, err, "Expected error when changing start date before current date")
}

func Test_whenChangeEndDateBeforeStartDate_generateError(t *testing.T) {
	task := NewTimedTaskEntity("Test Task", "This is a test task", time.Now(), time.Now().Add(2*time.Hour))
	err := task.changeEndDate(time.Now().Add(-1 * time.Hour))
	assert.Equal(t, ErrEndDateBeforeStartDate, err, "Expected error when changing end date before start date")
}

func Test_whenChangeEndDateBeforeCurrentDate_generateError(t *testing.T) {

	task := NewTimedTaskEntity("Test Task", "This is a test task", time.Now().Add(1*time.Second), time.Now().Add(5*time.Second))

	time.Sleep(2 * time.Second)

	err := task.changeEndDate(time.Now().Add(-1 * time.Second))
	assert.Equal(t, ErrEndDateBeforeNow, err, "Expected error when changing end date before current date")
}

func Test_toJSONString_generateSuccess(t *testing.T) {
	task := NewTimedTaskEntity("Test Task", "This is a test task", time.Now(), time.Now().Add(2*time.Hour))
	jsonString, err := task.ToJSONString()
	assert.Nil(t, err, "Expected no error when converting to JSON string")
	assert.Contains(t, jsonString, `"title":"Test Task"`, "Expected JSON string to contain task title")
	assert.Contains(t, jsonString, `"description":"This is a test task"`, "Expected JSON string to contain task description")
}

func Test_whenChangeEndDateWithValidDate_generateSuccess(t *testing.T) {
	task := NewTimedTaskEntity("Test Task", "This is a test task", time.Now(), time.Now().Add(5*time.Hour))
	newEndDate := time.Now().Add(10 * time.Hour)
	err := task.changeEndDate(newEndDate)
	assert.Nil(t, err, "Expected no error when changing end date to a valid date")
	assert.Equal(t, newEndDate, task.EndDate, "Expected end date to be updated to the new value")
}

func Test_whenChangeStartDateWithValidDate_generateSuccess(t *testing.T) {
	task := NewTimedTaskEntity("Test Task", "This is a test task", time.Now().Add(1*time.Hour), time.Now().Add(5*time.Hour))
	newStartDate := time.Now().Add(2 * time.Hour)
	err := task.changeStartDate(newStartDate)
	assert.Nil(t, err, "Expected no error when changing start date to a valid date")
	assert.Equal(t, newStartDate, task.StartDate, "Expected start date to be updated to the new value")
}
