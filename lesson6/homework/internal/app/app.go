package app

import (
	"fmt"
	"homework6/internal/ads"
)

type App interface {
	CreateAd(string, string, int64) (*ads.Ad, error)
	ChangeAdStatus(int64, int64, bool) (*ads.Ad, error)
	UpdateAd(int64, int64, string, string) (*ads.Ad, error)
}

type Repository interface {
	GetAdById(int64) (*ads.Ad, error)
	AddAd(*ads.Ad)
}

func NewApp(repo Repository) App {
	return &appRepo{repo, 0}
}

type appRepo struct {
	repository Repository
	counter    int64
}

func (a *appRepo) CreateAd(title string, text string, userId int64) (*ads.Ad, error) {
	ad := ads.Ad{ID: a.counter, Title: title, Text: text, AuthorID: userId}
	a.repository.AddAd(&ad)
	a.counter++
	return &ad, nil
}

func (a *appRepo) ChangeAdStatus(adId int64, userId int64, published bool) (*ads.Ad, error) {
	ad, err := a.repository.GetAdById(adId)
	if err != nil {
		return ad, err
	}
	if userId != ad.AuthorID {
		return ad, fmt.Errorf("incorrect userId")
	}
	ad.Published = published
	return ad, nil
}

func (a *appRepo) UpdateAd(adId int64, userId int64, title string, text string) (*ads.Ad, error) {
	ad, err := a.repository.GetAdById(adId)
	if err != nil {
		return ad, err
	}
	if userId != ad.AuthorID {
		return ad, fmt.Errorf("incorrect userId")
	}
	ad.Text = text
	ad.Title = title
	return ad, nil
}
