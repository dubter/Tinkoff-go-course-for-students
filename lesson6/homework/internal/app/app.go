package app

import (
	"fmt"
	"github.com/dubter/Validator"
	"homework6/internal/ads"
)

type App interface {
	CreateAd(title string, text string, userId int64) (*ads.Ad, error)
	ChangeAdStatus(adId int64, userId int64, published bool) (*ads.Ad, error)
	UpdateAd(adId int64, userId int64, title string, text string) (*ads.Ad, error)
}

type Repository interface {
	GetAdById(id int64) (ads.Ad, error)
	AddAd(ad *ads.Ad)
	GetPrimaryKey() int64
	Update(ad *ads.Ad) bool
}

func NewApp(repo Repository) App {
	return &appRepo{repo}
}

type appRepo struct {
	repository Repository
}

func (a *appRepo) CreateAd(title string, text string, userId int64) (*ads.Ad, error) {
	ad := ads.Ad{ID: a.repository.GetPrimaryKey() + 1, Title: title, Text: text, AuthorID: userId}
	if Validator.Validate(ad) != nil {
		return &ad, fmt.Errorf("validation error")
	}
	a.repository.AddAd(&ad)
	return &ad, nil
}

func (a *appRepo) ChangeAdStatus(adId int64, userId int64, published bool) (*ads.Ad, error) {
	ad, err := a.repository.GetAdById(adId)
	if err != nil {
		return &ad, err
	}

	ad.Published = published
	if userId != ad.AuthorID {
		return &ad, fmt.Errorf("incorrect userId")
	}
	if Validator.Validate(ad) != nil {
		return nil, fmt.Errorf("validation error")
	}

	a.repository.Update(&ad)
	return &ad, nil
}

func (a *appRepo) UpdateAd(adId int64, userId int64, title string, text string) (*ads.Ad, error) {
	ad, err := a.repository.GetAdById(adId)
	if err != nil {
		return &ad, err
	}

	ad.Text = text
	ad.Title = title
	if userId != ad.AuthorID {
		return &ad, fmt.Errorf("incorrect userId")
	}
	if Validator.Validate(ad) != nil {
		return &ad, fmt.Errorf("validation error")
	}

	a.repository.Update(&ad)
	return &ad, nil
}
