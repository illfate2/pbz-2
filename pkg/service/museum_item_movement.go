package service

import (
	"log"

	"pbz2/pkg/entities"
)

func (s *Service) CreateMuseumItemMovement(movement entities.MuseumItemMovement) (entities.MuseumItemMovement, error) {
	var err error
	movement.ResponsiblePerson, err = s.CreatePerson(movement.ResponsiblePerson)
	if err != nil {
		log.Printf("failed to insert person: %s", err)
		return entities.MuseumItemMovement{}, err
	}
	movement.ResponsiblePersonID = movement.ResponsiblePerson.ID
	item, err := s.GetMuseumItemByName(movement.Item.Name)
	if err != nil {
		log.Printf("failed to find movement: %s", err)
		return entities.MuseumItemMovement{}, err
	}
	movement.MuseumItemID = item.ID
	movement, err = s.insertMuseumItemMovement(movement)
	if err != nil {
		log.Print(err)
		return entities.MuseumItemMovement{}, err
	}
	return movement, nil
}

func (s *Service) GetMuseumItemMovement(id int) (entities.MuseumItemMovement, error) {
	return s.repo.FindMuseumItemMovement(id)
}

func (s *Service) GetMuseumItemMovements() ([]entities.MuseumItemMovement, error) {
	return s.repo.FindMuseumItemMovements()
}

func (s *Service) insertMuseumItemMovement(movement entities.MuseumItemMovement) (entities.MuseumItemMovement, error) {
	return s.repo.InsertMuseumItemMovement(movement)
}

func (s *Service) UpdateMuseumItemMovement(item entities.MuseumItem) error {
	return s.UpdateMuseumItemMovement(item)
}

func (s *Service) DeleteMuseumItemMovement(id int) error {
	return s.repo.DeleteMuseumItemMovement(id)
}
