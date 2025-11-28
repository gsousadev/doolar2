package entity

import (
	"errors"
	"time"
)

type PersonRepositoryInterface interface {
	Save(person *Person) error
}

type Person struct {
	Name      string    `json:"name"`
	Nickname  string    `json:"nickname"`
	BirthDate time.Time `json:"birth_date"`
}

func (p *Person) validate() error {
	if p.Name == "" {
		return errors.New("name cannot be empty")
	}
	if p.Nickname == "" {
		return errors.New("nickname cannot be empty")
	}
	if p.BirthDate.IsZero() {
		return errors.New("birth date cannot be zero")
	}
	if p.BirthDate.After(time.Now()) {
		return errors.New("birth date should not be in the future")
	}
	return nil
}

func NewPerson(name, nickname string, birthDate time.Time) (*Person, error) {

	person := &Person{
		Name:      name,
		Nickname:  nickname,
		BirthDate: birthDate,
	}

	if err := person.validate(); err != nil {
		return nil, err
	}
	return person, nil
}

func (p *Person) GetAge() uint8 {
	age := time.Now().Year() - p.BirthDate.Year()
	if time.Now().YearDay() < p.BirthDate.YearDay() {
		age--
	}
	return uint8(age)
}
