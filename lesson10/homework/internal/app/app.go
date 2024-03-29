package app

import (
	"errors"
	"github.com/dubter/Validator"
	"homework10/internal/ads"
	"homework10/internal/users"
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
	DeleteAd(adId int64, userId int64) error

	GetAd(id int64) (*ads.Ad, error)
	GetListAds(filters map[string]any) []ads.Ad
	GetListAdsByTitle(pattern string) []ads.Ad

	CreateUser(nickname string, email string) (*users.User, error)
	UpdateUser(userId int64, nickname string, email string) (*users.User, error)
	DeleteUser(userId int64) error
	GetUser(userId int64) (*users.User, error)
}

type Repository interface {
	GetAdById(id int64) (ads.Ad, error)
	AddAd(ad *ads.Ad)
	GetAdsPrimaryKey() int64
	ChangeAd(ad *ads.Ad) bool

	GetAds(filters map[string]any) []ads.Ad
	GetAdsByTitle(pattern string) []ads.Ad

	GetUserById(id int64) (users.User, error)
	AddUser(user *users.User)
	ChangeUser(user *users.User) bool
	GetUsersPrimaryKey() int64
	DeleteAd(adId int64)
	DeleteUser(uerId int64)
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
	now := time.Now().UTC()
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
	ad.DateUpdate = time.Now().UTC()
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
	ad.DateUpdate = time.Now().UTC()
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

func (a *appRepo) GetListAdsByTitle(pattern string) []ads.Ad {
	return a.repository.GetAdsByTitle(pattern)
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

func (a *appRepo) GetUser(userId int64) (*users.User, error) {
	user, err := a.repository.GetUserById(userId)
	return &user, err
}

func (a *appRepo) DeleteUser(userId int64) error {
	_, err := a.repository.GetUserById(userId)
	if err != nil {
		return err
	}

	a.repository.DeleteUser(userId)
	return nil
}

func (a *appRepo) DeleteAd(adId int64, userId int64) error {
	ad, err := a.repository.GetAdById(adId)
	if err != nil {
		return err
	}

	if ad.AuthorID != userId {
		return IncorrectUserId
	}

	a.repository.DeleteAd(adId)
	return nil
}
