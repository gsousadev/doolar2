package task_list

import (
	"encoding/json"
	"errors"
	"time"
)

var ErrEndDateBeforeStartDate = errors.New("end date cannot be before start date")
var ErrEndDateBeforeNow = errors.New("end date cannot be before current date")
var ErrStartDateAfterEndDate = errors.New("start date cannot be after end date")
var ErrStartDateBeforeNow = errors.New("start date cannot be before current date")

type TimedTaskEntity struct {
	*TaskEntity
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
}

func NewTimedTaskEntity(title, description string, startDate, endDate time.Time) *TimedTaskEntity {
	return &TimedTaskEntity{
		TaskEntity: NewTaskEntity(title, description),
		StartDate:  startDate,
		EndDate:    endDate,
	}
}

func (t *TimedTaskEntity) changeStartDate(newStartDate time.Time) error {

	if newStartDate.After(t.EndDate) {
		return ErrStartDateAfterEndDate
	}

	if newStartDate.Before(time.Now()) {
		return ErrStartDateBeforeNow
	}

	t.StartDate = newStartDate
	return nil
}

func (t *TimedTaskEntity) changeEndDate(newEndDate time.Time) error {

	if newEndDate.Before(t.StartDate) {
		return ErrEndDateBeforeStartDate
	}

	if newEndDate.Before(time.Now()) {
		return ErrEndDateBeforeNow
	}

	t.EndDate = newEndDate
	return nil
}

func (t *TimedTaskEntity) ToJSONString() (string, error) {
	jsonBytes, err := t.ToJSON()
	if err != nil {
		return "", err
	}
	return string(jsonBytes), nil
}

func (t *TimedTaskEntity) ToJSON() ([]byte, error) {
	return json.Marshal(t)
}
