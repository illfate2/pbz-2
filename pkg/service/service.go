package service

import (
	"pbz2/pkg/repo"
)

func NewService(repo *repo.Repo) *Service {
	return &Service{
		repo: repo,
	}
}

type Service struct {
	repo *repo.Repo
}
