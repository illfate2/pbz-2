package service

import "context"

func (s *Service) InsertPerson(person Person) (Person, error) {
	err := s.conn.QueryRow(context.Background(),
		`INSERT INTO persons(first_name,second_name,middle_name) VALUES($1,$2,$3) RETURNING id`,
		person.FirstName, person.LastName, person.MiddleName).
		Scan(&person.ID)
	if err != nil {
		return Person{}, err
	}
	return person, nil
}
