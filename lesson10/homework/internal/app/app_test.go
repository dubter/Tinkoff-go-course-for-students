package app

import (
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"homework10/internal/ads"
	"homework10/internal/app/mocks"
	"homework10/internal/users"
	"testing"
	"time"
)

const dateFormat string = "2006-01-02"
const one int64 = 1

type AppRepoTestSuite struct {
	suite.Suite
	repo mocks.Repository
}

func (s *AppRepoTestSuite) SetupTest() {
	s.repo = mocks.Repository{}
}

func TestRepoRun(t *testing.T) {
	suite.Run(t, new(AppRepoTestSuite))
}

func (s *AppRepoTestSuite) TestAppRepo_CreateAd() {
	now := time.Now().UTC()
	expect := ads.Ad{ID: one, Title: "ad 1", Text: "text 1", AuthorID: one, Published: false, DateCreating: now, DateUpdate: now}
	s.repo.On("GetUserById", one).Return(users.User{}, nil)
	s.repo.On("GetAdsPrimaryKey").Return(one)
	s.repo.On("AddAd", mock.AnythingOfType("*ads.Ad"))

	service := NewApp(&s.repo)
	got, err := service.CreateAd(expect.Title, expect.Text, expect.AuthorID)
	s.NoError(err)

	expect.DateCreating = CutTime(expect.DateCreating)
	got.DateCreating = CutTime(got.DateCreating)
	expect.DateUpdate = CutTime(expect.DateUpdate)
	got.DateUpdate = CutTime(got.DateUpdate)

	s.Equal(*got, expect)
}

func (s *AppRepoTestSuite) TestAppRepo_ChangeAdStatus() {
	now := time.Now().UTC()
	expect := ads.Ad{ID: one, Title: "ad 1", Text: "text 1", AuthorID: one, Published: false, DateCreating: now, DateUpdate: now}
	s.repo.On("GetAdById", one).Return(expect, nil)
	s.repo.On("ChangeAd", mock.AnythingOfType("*ads.Ad")).Return(true)

	service := NewApp(&s.repo)
	expect.Published = true
	got, err := service.ChangeAdStatus(expect.ID, expect.AuthorID, expect.Published)
	s.NoError(err)

	got.DateUpdate = CutTime(got.DateUpdate)
	expect.DateUpdate = CutTime(expect.DateUpdate)

	s.Equal(*got, expect)
}

func (s *AppRepoTestSuite) TestAppRepo_UpdateAd() {
	now := time.Now().UTC()
	expect := ads.Ad{ID: one, Title: "ad 1", Text: "text 1", AuthorID: one, Published: false, DateCreating: now, DateUpdate: now}
	s.repo.On("GetAdById", one).Return(expect, nil)
	s.repo.On("ChangeAd", mock.AnythingOfType("*ads.Ad")).Return(true)

	service := NewApp(&s.repo)
	expect.Text = "text 2"
	expect.Title = "ad 2"
	got, err := service.UpdateAd(expect.ID, expect.AuthorID, expect.Title, expect.Text)
	s.NoError(err)

	got.DateUpdate = CutTime(got.DateUpdate)
	expect.DateUpdate = CutTime(expect.DateUpdate)

	s.Equal(*got, expect)
}

func (s *AppRepoTestSuite) TestAppRepo_DeleteAd() {
	now := time.Now().UTC()
	expect := ads.Ad{ID: one, Title: "ad 1", Text: "text 1", AuthorID: one, Published: false, DateCreating: now, DateUpdate: now}
	s.repo.On("GetAdById", one).Return(expect, nil)
	s.repo.On("DeleteAd", one)

	service := NewApp(&s.repo)
	err := service.DeleteAd(expect.ID, expect.AuthorID)
	s.NoError(err)
}

func (s *AppRepoTestSuite) TestAppRepo_DeleteAdIncorrectAdId() {
	now := time.Now().UTC()
	expect := ads.Ad{ID: one, Title: "ad 1", Text: "text 1", AuthorID: one, Published: false, DateCreating: now, DateUpdate: now}
	s.repo.On("GetAdById", one).Return(expect, IncorrectAdId)
	s.repo.On("DeleteAd", one)

	service := NewApp(&s.repo)
	err := service.DeleteAd(expect.ID, expect.AuthorID)
	s.ErrorIs(err, IncorrectAdId)
}

func (s *AppRepoTestSuite) TestAppRepo_DeleteAdIncorrectUserId() {
	now := time.Now().UTC()
	expect := ads.Ad{ID: one, Title: "ad 1", Text: "text 1", AuthorID: one, Published: false, DateCreating: now, DateUpdate: now}
	s.repo.On("GetAdById", one).Return(expect, nil)
	s.repo.On("DeleteAd", one)

	service := NewApp(&s.repo)
	err := service.DeleteAd(expect.ID, int64(2))
	s.ErrorIs(err, IncorrectUserId)
}

func (s *AppRepoTestSuite) TestAppRepo_GetAd() {
	now := time.Now().UTC()
	expect := ads.Ad{ID: one, Title: "ad 1", Text: "text 1", AuthorID: one, Published: false, DateCreating: now, DateUpdate: now}
	s.repo.On("GetAdById", one).Return(expect, nil)

	service := NewApp(&s.repo)
	got, err := service.GetAd(expect.ID)
	s.NoError(err)
	s.Equal(expect, *got)
}

func (s *AppRepoTestSuite) TestAppRepo_GetListAds() {
	now := time.Now().UTC()
	expect1 := ads.Ad{ID: one, Title: "ad 1", Text: "text 1", AuthorID: one, Published: true, DateCreating: now, DateUpdate: now}
	expect2 := ads.Ad{ID: one, Title: "ad 1", Text: "text 1", AuthorID: one, Published: true, DateCreating: now, DateUpdate: now}
	filters := map[string]any{
		"published":     true,
		"user_id":       one,
		"date_creating": now,
	}

	expectedList := []ads.Ad{expect1, expect2}
	s.repo.On("GetAds", filters).Return(expectedList, nil)

	service := NewApp(&s.repo)
	gotList := service.GetListAds(filters)
	s.Equal(gotList, expectedList)
}

func (s *AppRepoTestSuite) TestAppRepo_GetListAdsByTitle() {
	now := time.Now().UTC()
	expect1 := ads.Ad{ID: one, Title: "ad 1", Text: "text 1", AuthorID: one, Published: true, DateCreating: now, DateUpdate: now}
	expect2 := ads.Ad{ID: one, Title: "ad 1", Text: "text 1", AuthorID: one, Published: true, DateCreating: now, DateUpdate: now}
	pattern := "ad"

	expectedList := []ads.Ad{expect1, expect2}
	s.repo.On("GetAdsByTitle", pattern).Return(expectedList, nil)

	service := NewApp(&s.repo)
	gotList := service.GetListAdsByTitle(pattern)
	s.Equal(gotList, expectedList)
}

func CutTime(t time.Time) time.Time {
	tmp := t.Format(dateFormat)
	cutTime, _ := time.Parse(dateFormat, tmp)
	return cutTime
}

func (s *AppRepoTestSuite) TestAppRepo_CreateUser() {
	expect := users.User{ID: one, Nickname: "nickname 1", Email: "email"}

	s.repo.On("GetUsersPrimaryKey").Return(one)
	s.repo.On("AddUser", mock.AnythingOfType("*users.User"))

	service := NewApp(&s.repo)
	got, err := service.CreateUser(expect.Nickname, expect.Email)
	s.NoError(err)
	s.Equal(*got, expect)
}

func (s *AppRepoTestSuite) TestAppRepo_UpdateUser() {
	expect := users.User{ID: one, Nickname: "nickname 1", Email: "email 1"}

	s.repo.On("GetUserById", one).Return(expect, nil)
	s.repo.On("ChangeUser", mock.AnythingOfType("*users.User")).Return(true)

	service := NewApp(&s.repo)
	expect.Nickname = "nickname 2"
	expect.Email = "email 2"
	got, err := service.UpdateUser(expect.ID, expect.Nickname, expect.Email)
	s.NoError(err)
	s.Equal(*got, expect)
}

func (s *AppRepoTestSuite) TestAppRepo_UpdateUserIncorrectAdId() {
	expect := users.User{ID: one, Nickname: "nickname 1", Email: "email 1"}

	s.repo.On("GetUserById", one).Return(expect, IncorrectAdId)
	s.repo.On("ChangeUser", mock.AnythingOfType("*users.User")).Return(true)

	service := NewApp(&s.repo)
	expect.Nickname = "nickname 2"
	expect.Email = "email 2"
	_, err := service.UpdateUser(expect.ID, expect.Nickname, expect.Email)
	s.ErrorIs(err, IncorrectAdId)
}

func (s *AppRepoTestSuite) TestAppRepo_UpdateUserValidationErr() {
	expect := users.User{ID: one, Nickname: "nickname 1", Email: "email 1"}

	s.repo.On("GetUserById", one).Return(expect, nil)
	s.repo.On("ChangeUser", mock.AnythingOfType("*users.User")).Return(true)

	service := NewApp(&s.repo)
	expect.Nickname = ""
	expect.Email = "email 2"
	_, err := service.UpdateUser(expect.ID, expect.Nickname, expect.Email)
	s.ErrorIs(err, ValidateError)
}

func (s *AppRepoTestSuite) TestAppRepo_GetUser() {
	expect := users.User{ID: one, Nickname: "nickname 1", Email: "email 1"}

	s.repo.On("GetUserById", one).Return(expect, nil)

	service := NewApp(&s.repo)
	got, err := service.GetUser(expect.ID)
	s.NoError(err)
	s.Equal(*got, expect)
}

func (s *AppRepoTestSuite) TestAppRepo_DeleteUser() {
	expect := users.User{ID: one, Nickname: "nickname 1", Email: "email 1"}

	s.repo.On("GetUserById", expect.ID).Return(expect, nil)
	s.repo.On("DeleteUser", expect.ID)

	service := NewApp(&s.repo)
	err := service.DeleteUser(expect.ID)
	s.NoError(err)
}

func (s *AppRepoTestSuite) TestAppRepo_DeleteUserIncorrectUserId() {
	expect := users.User{ID: one, Nickname: "nickname 1", Email: "email 1"}

	s.repo.On("GetUserById", expect.ID).Return(expect, IncorrectUserId)
	s.repo.On("DeleteUser", expect.ID)

	service := NewApp(&s.repo)
	err := service.DeleteUser(expect.ID)
	s.ErrorIs(err, IncorrectUserId)
}

// tests which return errors
func (s *AppRepoTestSuite) TestAppRepo_CreateAdIncorrectUserId() {
	now := time.Now().UTC()
	expect := ads.Ad{ID: one, Title: "ad 1", Text: "text 1", AuthorID: one, Published: false, DateCreating: now, DateUpdate: now}
	s.repo.On("GetUserById", one).Return(users.User{}, IncorrectUserId)
	s.repo.On("GetAdsPrimaryKey").Return(one)
	s.repo.On("AddAd", mock.AnythingOfType("*ads.Ad"))

	service := NewApp(&s.repo)
	_, err := service.CreateAd(expect.Title, expect.Text, expect.AuthorID)
	s.ErrorIs(err, IncorrectUserId)
}

func (s *AppRepoTestSuite) TestAppRepo_CreateAdValidationErr() {
	now := time.Now().UTC()
	expect := ads.Ad{ID: one, Title: "", Text: "text 1", AuthorID: one, Published: false, DateCreating: now, DateUpdate: now}
	s.repo.On("GetUserById", one).Return(users.User{}, nil)
	s.repo.On("GetAdsPrimaryKey").Return(one)
	s.repo.On("AddAd", mock.AnythingOfType("*ads.Ad"))

	service := NewApp(&s.repo)
	_, err := service.CreateAd(expect.Title, expect.Text, expect.AuthorID)
	s.ErrorIs(err, ValidateError)
}

func (s *AppRepoTestSuite) TestAppRepo_CreateUserValidationErr() {
	expect := users.User{ID: one, Nickname: "", Email: "email"}

	s.repo.On("GetUsersPrimaryKey").Return(one)
	s.repo.On("AddUser", mock.AnythingOfType("*users.User"))

	service := NewApp(&s.repo)
	_, err := service.CreateUser(expect.Nickname, expect.Email)
	s.ErrorIs(err, ValidateError)
}

func (s *AppRepoTestSuite) TestAppRepo_ChangeAdStatusValidationErr() {
	now := time.Now().UTC()
	expect := ads.Ad{ID: one, Title: "", Text: "text 1", AuthorID: one, Published: false, DateCreating: now, DateUpdate: now}
	s.repo.On("GetAdById", one).Return(expect, nil)
	s.repo.On("ChangeAd", mock.AnythingOfType("*ads.Ad")).Return(true)

	service := NewApp(&s.repo)
	expect.Published = true
	_, err := service.ChangeAdStatus(expect.ID, expect.AuthorID, expect.Published)
	s.ErrorIs(err, ValidateError)
}

func (s *AppRepoTestSuite) TestAppRepo_ChangeAdStatusIncorrectUserId() {
	now := time.Now().UTC()
	expect := ads.Ad{ID: one, Title: "title", Text: "text 1", AuthorID: one, Published: false, DateCreating: now, DateUpdate: now}
	s.repo.On("GetAdById", one).Return(expect, nil)
	s.repo.On("ChangeAd", mock.AnythingOfType("*ads.Ad")).Return(true)

	service := NewApp(&s.repo)
	expect.Published = true
	_, err := service.ChangeAdStatus(expect.ID, int64(2), expect.Published)
	s.ErrorIs(err, IncorrectUserId)
}

func (s *AppRepoTestSuite) TestAppRepo_ChangeAdStatusIncorrectAdId() {
	now := time.Now().UTC()
	expect := ads.Ad{ID: one, Title: "title", Text: "text 1", AuthorID: one, Published: false, DateCreating: now, DateUpdate: now}
	s.repo.On("GetAdById", one).Return(expect, IncorrectAdId)
	s.repo.On("ChangeAd", mock.AnythingOfType("*ads.Ad")).Return(true)

	service := NewApp(&s.repo)
	expect.Published = true
	_, err := service.ChangeAdStatus(expect.ID, expect.AuthorID, expect.Published)
	s.ErrorIs(err, IncorrectAdId)
}

func (s *AppRepoTestSuite) TestAppRepo_UpdateAdIncorrectAdId() {
	now := time.Now().UTC()
	expect := ads.Ad{ID: one, Title: "title", Text: "text 1", AuthorID: one, Published: false, DateCreating: now, DateUpdate: now}
	s.repo.On("GetAdById", one).Return(expect, IncorrectAdId)

	service := NewApp(&s.repo)
	expect.Title = "new title"
	expect.Text = "new text"
	_, err := service.UpdateAd(expect.ID, expect.AuthorID, expect.Title, expect.Text)
	s.ErrorIs(err, IncorrectAdId)
}

func (s *AppRepoTestSuite) TestAppRepo_UpdateAdValidationErr() {
	now := time.Now().UTC()
	expect := ads.Ad{ID: one, Title: "title", Text: "text 1", AuthorID: one, Published: false, DateCreating: now, DateUpdate: now}
	s.repo.On("GetAdById", one).Return(expect, nil)

	service := NewApp(&s.repo)
	expect.Title = ""
	_, err := service.UpdateAd(expect.ID, expect.AuthorID, expect.Title, expect.Text)
	s.ErrorIs(err, ValidateError)
}

func (s *AppRepoTestSuite) TestAppRepo_UpdateAdIncorrectUserId() {
	now := time.Now().UTC()
	expect := ads.Ad{ID: one, Title: "title", Text: "text 1", AuthorID: one, Published: false, DateCreating: now, DateUpdate: now}
	s.repo.On("GetAdById", one).Return(expect, nil)
	s.repo.On("ChangeAd", mock.AnythingOfType("*ads.Ad")).Return(true)

	service := NewApp(&s.repo)
	expect.Title = "new title"
	expect.Title = "new text"
	_, err := service.UpdateAd(expect.ID, int64(2), expect.Title, expect.Text)
	s.ErrorIs(err, IncorrectUserId)
}
