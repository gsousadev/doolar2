package task_list

import (
	"testing"
	"time"

	house_entity "github.com/gsousadev/doolar2/internal/house/domain/entity"
	"github.com/stretchr/testify/assert"
)

func Test_NewHomeTaskWithTimeLimitEntity_generateSuccess(t *testing.T) {
	room := house_entity.NewRoom("Living Room")
	task := NewTimedHomeTask(room, "Test Home Task", "This is a test home task", time.Now(), time.Now().Add(2*time.Hour))
	assert.IsType(t, task, &TimedHomeTask{})
	assert.IsType(t, task.TimedTaskEntity, &TimedTaskEntity{})
	assert.IsType(t, task.Room, &house_entity.Room{})
	assert.Equal(t, "Living Room", task.Room.Name, "Expected room name to be 'Living Room'")
	assert.Equal(t, StatusPending, task.GetStatus(), "Expected new home task to have status 'pending'")
}
