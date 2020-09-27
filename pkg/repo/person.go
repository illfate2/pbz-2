package repo

import (
	"context"

	"pbz2/pkg/entities"
)

func (r *Repo) InsertPerson(person entities.Person) (entities.Person, error) {
	err := r.conn.QueryRow(context.Background(),
		`INSERT INTO persons(first_name,second_name,middle_name) VALUES($1,$2,$3) RETURNING id`,
		person.FirstName, person.LastName, person.MiddleName).
		Scan(&person.ID)
	if err != nil {
		return entities.Person{}, err
	}
	return person, nil
}
