package repo

import (
	"context"

	"pbz2/pkg/entities"
)

func (r *Repo) InsertMuseumFund(fund entities.MuseumFund) (entities.MuseumFund, error) {
	err := r.conn.QueryRow(context.Background(),
		`INSERT INTO museum_funds(name)
			VALUES($1)
			ON CONFLICT (name)
			DO UPDATE SET name=EXCLUDED.name
			RETURNING id`,
		fund.Name).
		Scan(&fund.ID)
	if err != nil {
		return entities.MuseumFund{}, err
	}
	return fund, nil
}
