package repo

import (
	"context"

	"pbz2/pkg/entities"
)

func (r *Repo) FindMuseumItemMovement(id int) (entities.MuseumItemMovement, error) {
	var m entities.MuseumItemMovement
	err := r.conn.QueryRow(context.Background(),
		`SELECT id,item_id,responsible_person_id,accept_date,exhibit_transfer_date,exhibit_return_date FROM museum_item_movements WHERE id = $1`, id).
		Scan(&m.ID, &m.MuseumItemID, &m.ResponsiblePersonID, &m.AcceptDate, &m.ExhibitTransferDate, &m.ExhibitReturnDate)
	if err != nil {
		return entities.MuseumItemMovement{}, err
	}
	return m, nil
}

func (r *Repo) FindMuseumItemMovements() ([]entities.MuseumItemMovement, error) {
	var movements []entities.MuseumItemMovement
	rows, err := r.conn.Query(context.Background(),
		`SELECT id,item_id,responsible_person_id,accept_date,exhibit_transfer_date,exhibit_return_date FROM museum_item_movements`)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var m entities.MuseumItemMovement
		err := rows.Scan(&m.ID, &m.MuseumItemID, &m.ResponsiblePersonID, &m.AcceptDate, &m.ExhibitTransferDate, &m.ExhibitReturnDate)
		if err != nil {
			return nil, err
		}
		movements = append(movements, m)
	}
	return movements, nil
}

func (r *Repo) InsertMuseumItemMovement(movement entities.MuseumItemMovement) (entities.MuseumItemMovement, error) {
	err := r.conn.QueryRow(context.Background(),
		`INSERT INTO museum_item_movements(item_id,responsible_person_id,accept_date,exhibit_transfer_date,exhibit_return_date)
		VALUES($1,$2,$3,$4,$5) RETURNING id`,
		movement.MuseumItemID, movement.ResponsiblePersonID, movement.AcceptDate, movement.ExhibitTransferDate,
		movement.ExhibitReturnDate).
		Scan(&movement.ID)
	if err != nil {
		return entities.MuseumItemMovement{}, err
	}
	return movement, nil
}

func (r *Repo) UpdateMuseumItemMovement(item entities.MuseumItem) error {
	_, err := r.conn.Exec(context.Background(),
		`UPDATE museum_items 
			SET name = $1, creation_date = $2, annotation = $3
			WHERE id = $4`,
		item.Name, item.CreationDate.Time, item.Annotation, item.ID)
	return err
}

func (r *Repo) DeleteMuseumItemMovement(id int) error {
	_, err := r.conn.Exec(context.Background(),
		`DELETE FROM museum_item_movements WHERE id = $1`, id)
	return err
}
