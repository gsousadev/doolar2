package entity

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestPerson_GetAge_WhenValidBirthDate_ShouldReturnCorrectAge(t *testing.T) {
	birthDate := time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)

	person, err := NewPerson("John Doe", "johndoe", birthDate)
	age := person.GetAge()

	assert.NoError(t, err, "Expected no error when creating a valid person")
	assert.Equal(t, uint8(time.Now().Year()-2000), age, "Expected age to be the difference between current year and birth year")
}

func TestPerson_NewPerson_WhenBirthDateInFuture_ShouldReturnError(t *testing.T) {
	birthDate := time.Date(3000, 1, 1, 0, 0, 0, 0, time.UTC)
	person, err := NewPerson("Jane Doe", "janedoe", birthDate)
	assert.Nil(t, person, "Expected person to be nil when birth date is in the future")
	assert.Error(t, err, "Expected error when birth date is in the future")
	assert.EqualError(t, err, "birth date should not be in the future", "Expected specific error message for future birth date")
}

func TestPerson_NewPerson_WhenBirthDateIsZero_ShouldReturnError(t *testing.T) {
	birthDate := time.Time{} // Zero value for time.Time
	person, err := NewPerson("Alice", "alice", birthDate)
	assert.Nil(t, person, "Expected person to be nil when birth date is zero")
	assert.Error(t, err, "Expected error when birth date is zero")
	assert.EqualError(t, err, "birth date cannot be zero", "Expected specific error message for zero birth date")
}

func TestPerson_NewPerson_WhenNameIsEmpty_ShouldReturnError(t *testing.T) {
	birthDate := time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)
	person, err := NewPerson("", "nickname", birthDate)
	assert.Nil(t, person, "Expected person to be nil when name is empty")
	assert.Error(t, err, "Expected error when name is empty")
	assert.EqualError(t, err, "name cannot be empty", "Expected specific error message for empty name")
}

func TestPerson_NewPerson_WhenNicknameIsEmpty_ShouldReturnError(t *testing.T) {
	birthDate := time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)
	person, err := NewPerson("Name", "", birthDate)
	assert.Nil(t, person, "Expected person to be nil when nickname is empty")
	assert.Error(t, err, "Expected error when nickname is empty")
	assert.EqualError(t, err, "nickname cannot be empty", "Expected specific error message for empty nickname")
}
