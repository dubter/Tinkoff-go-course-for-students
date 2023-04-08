package adrepo

import (
	"fmt"
	"homework6/internal/ads"
	"homework6/internal/app"
)

type repositoryMap struct {
	dict    map[int64]ads.Ad
	counter int64
}

func (repo *repositoryMap) GetAdById(id int64) (ads.Ad, error) {
	ad, ok := repo.dict[id]
	if !ok {
		return ad, fmt.Errorf("there no such id in the table")
	}
	return ad, nil
}

func (repo *repositoryMap) AddAd(ad *ads.Ad) {
	repo.dict[ad.ID] = *ad
	repo.counter++
}

func New() app.Repository {
	return &repositoryMap{dict: make(map[int64]ads.Ad), counter: 0}
}

func (repo *repositoryMap) GetPrimaryKey() int64 {
	return repo.counter
}

func (repo *repositoryMap) Update(ad *ads.Ad) bool {
	_, ok := repo.dict[ad.ID]
	if ok {
		repo.dict[ad.ID] = *ad
	}
	return ok
}
