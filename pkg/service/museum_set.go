package service

import "context"

func (s *Service) InsertMuseumSet(set MuseumSet) (MuseumSet, error) {
	err := s.conn.QueryRow(context.Background(),
		`INSERT INTO museum_item_sets(name)
		VALUES($1)
		ON CONFLICT (name)
		DO UPDATE SET name=EXCLUDED.name
		RETURNING id`,
		set.Name).
		Scan(&set.ID)
	if err != nil {
		return MuseumSet{}, err
	}
	return set, nil
}

func (s *Service) FindMuseumSets() ([]MuseumSet, error) {
	rows, err := s.conn.Query(context.Background(),
		`SELECT 
			id, name
			FROM museum_item_sets
`)
	if err != nil {
		return nil, err
	}
	var sets []MuseumSet
	for rows.Next() {
		var set MuseumSet
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

func (s *Service) FindMuseumSet(id int) (MuseumSetWithDetails, error) {
	rows, err := s.conn.Query(context.Background(),
		`SELECT 
      mis.id, mis.name,
      mi.id ,mi.name, mi.creation_date, mi.annotation,
      p.id, p.first_name, p.second_name, p.middle_name
      FROM museum_item_sets mis
      LEFT JOIN museum_items mi ON mis.id=mi.set_id
      LEFT JOIN persons p ON mi.keeper_id=p.id
	WHERE mis.id = $1`, id)
	if err != nil {
		return MuseumSetWithDetails{}, err
	}
	var curSet MuseumSetWithDetails
	for rows.Next() {
		var item MuseumItemWithKeeper
		err := rows.Scan(
			&curSet.ID, &curSet.Name,
			&item.ID, &item.Name, &item.CreationDate.time, &item.Annotation,
			&item.Keeper.ID, &item.Keeper.FirstName, &item.Keeper.LastName, &item.Keeper.MiddleName,
		)
		if err != nil {
			return MuseumSetWithDetails{}, err
		}
		curSet.Items = append(curSet.Items, item)
	}
	return curSet, nil
}
