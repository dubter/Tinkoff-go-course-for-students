package adrepo

import (
	"fmt"
	"homework6/internal/ads"
	"homework6/internal/app"
)

type repositoryMap struct {
	dict map[int64]ads.Ad
}

func (repo *repositoryMap) GetAdById(id int64) (*ads.Ad, error) {
	value, ok := repo.dict[id]
	if !ok {
		return &value, fmt.Errorf("there no such id in the table")
	}
	return &value, nil
}

func (repo *repositoryMap) AddAd(ad *ads.Ad) {
	repo.dict[ad.ID] = *ad
}

func New() app.Repository {
	return &repositoryMap{dict: make(map[int64]ads.Ad)}
}
