package service

import (
	"context"
	"log"
)

func (s *Service) FindMuseumItemMovement(id int) (MuseumItemMovement, error) {
	var m MuseumItemMovement
	err := s.conn.QueryRow(context.Background(),
		`SELECT id,item_id,responsible_person_id,accept_date,exhibit_transfer_date,exhibit_return_date FROM museum_item_movements WHERE id = $1`, id).
		Scan(&m.ID, &m.MuseumItemID, &m.ResponsiblePersonID, &m.AcceptDate, &m.ExhibitTransferDate, &m.ExhibitReturnDate)
	if err != nil {
		return MuseumItemMovement{}, err
	}
	return m, nil
}

func (s *Service) FindMuseumItemMovements() ([]MuseumItemMovement, error) {
	var movements []MuseumItemMovement
	rows, err := s.conn.Query(context.Background(),
		`SELECT id,item_id,responsible_person_id,accept_date,exhibit_transfer_date,exhibit_return_date FROM museum_item_movements`)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var m MuseumItemMovement
		err := rows.Scan(&m.ID, &m.MuseumItemID, &m.ResponsiblePersonID, &m.AcceptDate, &m.ExhibitTransferDate, &m.ExhibitReturnDate)
		if err != nil {
			return nil, err
		}
		movements = append(movements, m)
	}
	return movements, nil
}

func (s *Service) CreateMuseumItemMovement(movement MuseumItemMovement) (MuseumItemMovement, error) {
	var err error
	movement.ResponsiblePerson, err = s.InsertPerson(movement.ResponsiblePerson)
	if err != nil {
		log.Printf("failed to insert person: %s", err)
		return MuseumItemMovement{}, err
	}
	movement.ResponsiblePersonID = movement.ResponsiblePerson.ID
	item, err := s.FindMuseumItemByName(movement.Item.Name)
	if err != nil {
		log.Printf("failed to find movement: %s", err)
		return MuseumItemMovement{}, err
	}
	movement.MuseumItemID = item.ID
	movement, err = s.InsertMuseumItemMovement(movement)
	if err != nil {
		log.Print(err)
		return MuseumItemMovement{}, err
	}
	return movement, nil
}

func (s *Service) InsertMuseumItemMovement(movement MuseumItemMovement) (MuseumItemMovement, error) {
	err := s.conn.QueryRow(context.Background(),
		`INSERT INTO museum_item_movements(item_id,responsible_person_id,accept_date,exhibit_transfer_date,exhibit_return_date)
		VALUES($1,$2,$3,$4,$5) RETURNING id`,
		movement.MuseumItemID, movement.ResponsiblePersonID, movement.AcceptDate, movement.ExhibitTransferDate,
		movement.ExhibitReturnDate).
		Scan(&movement.ID)
	if err != nil {
		return MuseumItemMovement{}, err
	}
	return movement, nil
}

func (s *Service) UpdateMuseumItemMovement(item MuseumItem) error {
	_, err := s.conn.Exec(context.Background(),
		`UPDATE museum_items 
			SET name = $1, creation_date = $2, annotation = $3
			WHERE id = $4`,
		item.Name, item.CreationDate.time, item.Annotation, item.ID)
	return err
}

func (s *Service) DeleteMuseumItemMovement(id int) error {
	_, err := s.conn.Exec(context.Background(),
		`DELETE FROM museum_item_movements WHERE id = $1`, id)
	return err
}
