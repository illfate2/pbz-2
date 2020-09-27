package service

import (
	"pbz2/pkg/entities"
)

func (s *Service) CreateMuseumFund(fund entities.MuseumFund) (entities.MuseumFund, error) {
	return s.repo.InsertMuseumFund(fund)
}
