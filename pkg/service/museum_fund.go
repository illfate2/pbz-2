package service

import "context"

func (s *Service) InsertMuseumFund(fund MuseumFund) (MuseumFund, error) {
	err := s.conn.QueryRow(context.Background(),
		`INSERT INTO museum_funds(name)
			VALUES($1)
			ON CONFLICT (name)
			DO UPDATE SET name=EXCLUDED.name
			RETURNING id`,
		fund.Name).
		Scan(&fund.ID)
	if err != nil {
		return MuseumFund{}, err
	}
	return fund, nil
}
