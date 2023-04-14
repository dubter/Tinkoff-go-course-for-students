package app

import (
	"errors"
	"github.com/dubter/Validator"
	"homework8/internal/ads"
	"homework8/internal/users"
	"time"
)

// some errors

var IncorrectUserId = errors.New("incorrect user id")
var ValidateError = errors.New("validation error")
var IncorrectAdId = errors.New("id is not found")

type App interface {
	CreateAd(title string, text string, userId int64) (*ads.Ad, error)
	ChangeAdStatus(adId int64, userId int64, published bool) (*ads.Ad, error)
	UpdateAd(adId int64, userId int64, title string, text string) (*ads.Ad, error)

	GetAd(id int64) (*ads.Ad, error)
	GetListAds(filters map[string]any) []ads.Ad
	GetListAdsByTitle(title string) []ads.Ad

	CreateUser(nickname string, email string) (*users.User, error)
	UpdateUser(userId int64, nickname string, email string) (*users.User, error)
}

type Repository interface {
	GetAdById(id int64) (ads.Ad, error)
	AddAd(ad *ads.Ad)
	GetAdsPrimaryKey() int64
	ChangeAd(ad *ads.Ad) bool

	GetAds(filters map[string]any) []ads.Ad
	GetAdsByTitle(title string) []ads.Ad

	GetUserById(id int64) (users.User, error)
	AddUser(user *users.User)
	ChangeUser(user *users.User) bool
	GetUsersPrimaryKey() int64
}

func NewApp(repo Repository) App {
	return &appRepo{repo}
}

type appRepo struct {
	repository Repository
}

func (a *appRepo) CreateAd(title string, text string, userId int64) (*ads.Ad, error) {
	if _, err := a.repository.GetUserById(userId); err != nil {
		return nil, IncorrectUserId
	}
	now := time.Now()
	ad := ads.Ad{ID: a.repository.GetAdsPrimaryKey(), Title: title, Text: text, AuthorID: userId, DateCreating: now, DateUpdate: now, Published: false}
	if Validator.Validate(ad) != nil {
		return nil, ValidateError
	}
	a.repository.AddAd(&ad)
	return &ad, nil
}

func (a *appRepo) CreateUser(nickname string, email string) (*users.User, error) {
	user := users.User{ID: a.repository.GetUsersPrimaryKey(), Nickname: nickname, Email: email}
	if Validator.Validate(user) != nil {
		return nil, ValidateError
	}

	a.repository.AddUser(&user)
	return &user, nil
}

func (a *appRepo) ChangeAdStatus(adId int64, userId int64, published bool) (*ads.Ad, error) {
	ad, err := a.repository.GetAdById(adId)
	if err != nil {
		return nil, err
	}

	ad.Published = published
	ad.DateUpdate = time.Now()
	if userId != ad.AuthorID {
		return nil, IncorrectUserId
	}
	if Validator.Validate(ad) != nil {
		return nil, ValidateError
	}

	a.repository.ChangeAd(&ad)
	return &ad, nil
}

func (a *appRepo) UpdateAd(adId int64, userId int64, title string, text string) (*ads.Ad, error) {
	ad, err := a.repository.GetAdById(adId)
	if err != nil {
		return nil, err
	}

	ad.Text = text
	ad.Title = title
	ad.DateUpdate = time.Now()
	if userId != ad.AuthorID {
		return nil, IncorrectUserId
	}
	if Validator.Validate(ad) != nil {
		return nil, ValidateError
	}

	a.repository.ChangeAd(&ad)
	return &ad, nil
}

func (a *appRepo) GetAd(id int64) (*ads.Ad, error) {
	ad, err := a.repository.GetAdById(id)
	return &ad, err
}

func (a *appRepo) GetListAds(filters map[string]any) []ads.Ad {
	return a.repository.GetAds(filters)
}

func (a *appRepo) GetListAdsByTitle(title string) []ads.Ad {
	return a.repository.GetAdsByTitle(title)
}

func (a *appRepo) UpdateUser(userId int64, nickname string, email string) (*users.User, error) {
	user, err := a.repository.GetUserById(userId)
	if err != nil {
		return nil, err
	}

	user.Nickname = nickname
	user.Email = email
	if Validator.Validate(user) != nil {
		return nil, ValidateError
	}

	a.repository.ChangeUser(&user)
	return &user, nil
}
