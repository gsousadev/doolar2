package database

import (
	"github.com/gsousadev/doolar-golang/internal/domain/entity"
	"database/sql"
)

type PersonRepository struct {
	Db *sql.DB
}

func (r *PersonRepository) Save(person *entity.Person) error {
	if person == nil {
		return nil // or return an error if you prefer
	}

	prepare, err := r.Db.Prepare("INSERT INTO people (name, nickname, birth_date) VALUES (?, ?, ?)")
	if err != nil {
		return err
	}
	
	defer prepare.Close()

	_, err = prepare.Exec(person.Name, person.Nickname, person.BirthDate)

	if err != nil {
		return err
	}
	
	return nil
}

func NewPersonRepository(db *sql.DB) *PersonRepository {
	return &PersonRepository{Db: db}
}