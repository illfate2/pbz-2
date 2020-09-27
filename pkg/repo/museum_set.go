package repo

import (
	"context"

	"pbz2/pkg/entities"
)

func (r *Repo) InsertMuseumSet(set entities.MuseumSet) (entities.MuseumSet, error) {
	err := r.conn.QueryRow(context.Background(),
		`INSERT INTO museum_item_sets(name)
		VALUES($1)
		ON CONFLICT (name)
		DO UPDATE SET name=EXCLUDED.name
		RETURNING id`,
		set.Name).
		Scan(&set.ID)
	if err != nil {
		return entities.MuseumSet{}, err
	}
	return set, nil
}

func (r *Repo) FindMuseumSets() ([]entities.MuseumSet, error) {
	rows, err := r.conn.Query(context.Background(),
		`SELECT 
			id, name
			FROM museum_item_sets
`)
	if err != nil {
		return nil, err
	}
	var sets []entities.MuseumSet
	for rows.Next() {
		var set entities.MuseumSet
		err := rows.Scan(
			&set.ID, &set.Name,
		)
		if err != nil {
			return nil, err
		}
		sets = append(sets, set)
	}
	return sets, nil
}

func (r *Repo) FindMuseumSet(id int) (entities.MuseumSetWithDetails, error) {
	rows, err := r.conn.Query(context.Background(),
		`SELECT 
      mis.id, mis.name,
      mi.id ,mi.name, mi.creation_date, mi.annotation,
      p.id, p.first_name, p.second_name, p.middle_name
      FROM museum_item_sets mis
      LEFT JOIN museum_items mi ON mis.id=mi.set_id
      LEFT JOIN persons p ON mi.keeper_id=p.id
	WHERE mis.id = $1`, id)
	if err != nil {
		return entities.MuseumSetWithDetails{}, err
	}
	var curSet entities.MuseumSetWithDetails
	for rows.Next() {
		var item entities.MuseumItemWithKeeper
		err := rows.Scan(
			&curSet.ID, &curSet.Name,
			&item.ID, &item.Name, &item.CreationDate.Time, &item.Annotation,
			&item.Keeper.ID, &item.Keeper.FirstName, &item.Keeper.LastName, &item.Keeper.MiddleName,
		)
		if err != nil {
			return entities.MuseumSetWithDetails{}, err
		}
		curSet.Items = append(curSet.Items, item)
	}
	return curSet, nil
}
