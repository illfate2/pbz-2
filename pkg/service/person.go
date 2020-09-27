package service

import (
	"pbz2/pkg/entities"
)

func (s *Service) CreatePerson(person entities.Person) (entities.Person, error) {
	person, err := s.repo.InsertPerson(person)
	if err != nil {
		return entities.Person{}, err
	}
	return person, nil
}
