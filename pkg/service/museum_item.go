package service

import (
	"github.com/pkg/errors"

	"pbz2/pkg/entities"
)

// todo: add tx
func (s *Service) CreateMuseumItem(item entities.MuseumItemWithDetails) (entities.MuseumItemWithDetails, error) {
	var err error
	item.Keeper, err = s.CreatePerson(item.Keeper)
	if err != nil {
		return entities.MuseumItemWithDetails{}, err
	}
	item.KeeperID = item.Keeper.ID

	item.Fund, err = s.CreateMuseumFund(item.Fund)
	if err != nil {
		return entities.MuseumItemWithDetails{}, errors.Wrap(err, "failed to insert musuem fund")
	}
	item.MuseumFundID = item.Fund.ID

	item.Set, err = s.CreateMuseumSet(item.Set)
	if err != nil {
		return entities.MuseumItemWithDetails{}, errors.Wrap(err, "failed to insert set")
	}
	item.MuseumSetID = item.Set.ID

	item.MuseumItem, err = s.insertMuseumItem(item.MuseumItem)
	if err != nil {
		return entities.MuseumItemWithDetails{}, errors.Wrap(err, "failed to insert item")
	}

	return item, nil
}

func (s *Service) GetMuseumItem(id int) (entities.MuseumItem, error) {
	return s.repo.FindMuseumItem(id)
}

func (s *Service) GetMuseumItemByName(name string) (entities.MuseumItem, error) {
	return s.repo.FindMuseumItemByName(name)
}

func (s *Service) GetMuseumItemWithDetails(id int) (entities.MuseumItemWithDetails, error) {
	return s.repo.FindMuseumItemWithDetails(id)
}

func (s *Service) FindMuseumItems(args entities.SearchMuseumItemsArgs) ([]entities.MuseumItem, error) {
	return s.repo.FindMuseumItems(args)
}

func (s *Service) UpdateMuseumItem(item entities.MuseumItem) error {
	return s.repo.UpdateMuseumItem(item)
}

func (s *Service) DeleteMuseumItem(id int) error {
	return s.repo.DeleteMuseumItem(id)
}

func (s *Service) insertMuseumItem(item entities.MuseumItem) (entities.MuseumItem, error) {
	return s.repo.InsertMuseumItem(item)
}
