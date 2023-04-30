package adrepo

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"homework10/internal/ads"
	"homework10/internal/app"
	"homework10/internal/users"
	"testing"
	"time"
)

type RepositoryMapTestSuite struct {
	suite.Suite
	repo app.Repository
}

func (s *RepositoryMapTestSuite) SetupTest() {
	s.repo = New()
}

func (s *RepositoryMapTestSuite) TearDownTest() {
	s.repo = nil
}

func TestRepoRun(t *testing.T) {
	suite.Run(t, new(RepositoryMapTestSuite))
}

func (s *RepositoryMapTestSuite) TestAddAd() {
	// Test case for adding a new ad
	expectedAd := ads.Ad{ID: 1, Title: "Ad 1", Text: "Ad 1 description", AuthorID: 1, Published: true}
	s.repo.AddAd(&expectedAd)
	ad, err := s.repo.GetAdById(1)
	s.NoError(err)
	s.Equal(expectedAd, ad)
}

func (s *RepositoryMapTestSuite) TestGetAdById() {
	// Test case for a valid ad ID
	expectedAd := ads.Ad{ID: 1, Title: "Ad 1", Text: "Ad 1 description", AuthorID: 1, Published: true}
	s.repo.AddAd(&expectedAd)
	ad, err := s.repo.GetAdById(1)
	s.NoError(err)
	s.Equal(expectedAd, ad)

	// Test case for an invalid ad ID
	_, err = s.repo.GetAdById(2)
	s.ErrorIs(err, app.IncorrectAdId)
}

func (s *RepositoryMapTestSuite) TestChangeAd() {
	// Test case for changing an existing ad
	expectedAd := ads.Ad{ID: 1, Title: "Ad 1", Text: "Ad 1 description", AuthorID: 1, Published: true}
	s.repo.AddAd(&expectedAd)
	expectedAd.Text = "Updated description"
	s.True(s.repo.ChangeAd(&expectedAd))
	ad, err := s.repo.GetAdById(1)
	s.NoError(err)
	s.Equal(expectedAd, ad)

	// Test case for changing a non-existing ad
	expectedAd.ID = 2
	s.False(s.repo.ChangeAd(&expectedAd))
}

func (s *RepositoryMapTestSuite) TestRepositoryMap_DeleteAd() {
	// Test case for delete an existing ad
	expectedAd := ads.Ad{ID: 1, Title: "Ad 1", Text: "Ad 1 description", AuthorID: 1, Published: true}
	s.repo.AddAd(&expectedAd)

	s.repo.DeleteAd(expectedAd.ID)
	_, err := s.repo.GetAdById(1)
	s.ErrorIs(err, app.IncorrectAdId)
}

func (s *RepositoryMapTestSuite) TestRepositoryMap_GetAds() {
	ad1 := ads.Ad{ID: 1, Title: "Ad 1", Text: "Ad 1 description", AuthorID: 1, Published: true}
	ad2 := ads.Ad{ID: 2, Title: "Ad 2", Text: "Ad 2 description", AuthorID: 1, Published: true}
	ad3 := ads.Ad{ID: 3, Title: "Ad 3", Text: "Ad 3 description", AuthorID: 1, Published: true}

	s.repo.AddAd(&ad1)
	s.repo.AddAd(&ad2)
	s.repo.AddAd(&ad3)

	adsList := s.repo.GetAds(nil)
	s.Len(adsList, 3)
}

func (s *RepositoryMapTestSuite) TestRepositoryMap_GetAdsWithFiltersPublished() {
	ad1 := ads.Ad{ID: 1, Title: "Ad 1", Text: "Ad 1 description", AuthorID: 1, Published: true}
	ad2 := ads.Ad{ID: 2, Title: "Ad 2", Text: "Ad 2 description", AuthorID: 1, Published: false}
	ad3 := ads.Ad{ID: 3, Title: "Ad 3", Text: "Ad 3 description", AuthorID: 1, Published: true}

	s.repo.AddAd(&ad1)
	s.repo.AddAd(&ad2)
	s.repo.AddAd(&ad3)

	filters := map[string]any{
		"published": true,
	}
	adsList := s.repo.GetAds(filters)
	s.Len(adsList, 2)
}

func (s *RepositoryMapTestSuite) TestRepositoryMap_GetAdsWithFiltersUserId() {
	ad1 := ads.Ad{ID: 1, Title: "Ad 1", Text: "Ad 1 description", AuthorID: 2, Published: true}
	ad2 := ads.Ad{ID: 2, Title: "Ad 2", Text: "Ad 2 description", AuthorID: 1, Published: false}
	ad3 := ads.Ad{ID: 3, Title: "Ad 3", Text: "Ad 3 description", AuthorID: 1, Published: true}

	s.repo.AddAd(&ad1)
	s.repo.AddAd(&ad2)
	s.repo.AddAd(&ad3)

	filters := map[string]any{
		"user_id": int64(1),
	}
	adsList := s.repo.GetAds(filters)
	s.Len(adsList, 2)
}

func (s *RepositoryMapTestSuite) TestRepositoryMap_GetAdsWithAllFilters() {
	timeStr1 := "2023-04-05"
	timeStr2 := "2023-04-10"

	time1, err := time.Parse(dateFormat, timeStr1)
	s.NoError(err)
	time2, err := time.Parse(dateFormat, timeStr2)
	s.NoError(err)

	ad1 := ads.Ad{ID: 1, Title: "Ad 1", Text: "Ad 1 description", AuthorID: 2, Published: true, DateCreating: time1}
	ad2 := ads.Ad{ID: 2, Title: "Ad 2", Text: "Ad 2 description", AuthorID: 1, Published: false, DateCreating: time1}
	ad3 := ads.Ad{ID: 3, Title: "Ad 3", Text: "Ad 3 description", AuthorID: 1, Published: true, DateCreating: time2}
	ad4 := ads.Ad{ID: 4, Title: "Ad 4", Text: "Ad 4 description", AuthorID: 1, Published: true, DateCreating: time1}
	ad5 := ads.Ad{ID: 5, Title: "Ad 5", Text: "Ad 5 description", AuthorID: 1, Published: true, DateCreating: time1}

	s.repo.AddAd(&ad1)
	s.repo.AddAd(&ad2)
	s.repo.AddAd(&ad3)
	s.repo.AddAd(&ad4)
	s.repo.AddAd(&ad5)

	filters := map[string]any{
		"user_id":       int64(1),
		"published":     true,
		"date_creating": timeStr1,
	}
	adsList := s.repo.GetAds(filters)
	s.Len(adsList, 2)
}

func (s *RepositoryMapTestSuite) TestRepositoryMap_GetAdsPrimaryKey() {
	ad1 := ads.Ad{ID: 1, Title: "Ad 1", Text: "Ad 1 description", AuthorID: 2, Published: true}
	ad2 := ads.Ad{ID: 2, Title: "Ad 2", Text: "Ad 2 description", AuthorID: 1, Published: false}
	ad3 := ads.Ad{ID: 3, Title: "Ad 3", Text: "Ad 3 description", AuthorID: 1, Published: true}
	ad4 := ads.Ad{ID: 4, Title: "Ad 4", Text: "Ad 4 description", AuthorID: 1, Published: true}
	ad5 := ads.Ad{ID: 5, Title: "Ad 5", Text: "Ad 5 description", AuthorID: 1, Published: true}

	s.repo.AddAd(&ad1)
	s.repo.AddAd(&ad2)
	s.repo.AddAd(&ad3)
	s.repo.AddAd(&ad4)
	s.repo.AddAd(&ad5)

	primaryKey := s.repo.GetAdsPrimaryKey()
	s.Equal(primaryKey, int64(5))
}

func (s *RepositoryMapTestSuite) TestRepositoryMap_GetUsersPrimaryKey() {
	user1 := users.User{ID: 1, Nickname: "nickname 1", Email: "email 1"}
	user2 := users.User{ID: 2, Nickname: "nickname 2", Email: "email 2"}
	user3 := users.User{ID: 3, Nickname: "nickname 3", Email: "email 3"}

	s.repo.AddUser(&user1)
	s.repo.AddUser(&user2)
	s.repo.AddUser(&user3)

	primaryKey := s.repo.GetUsersPrimaryKey()
	s.Equal(primaryKey, int64(3))
}

func (s *RepositoryMapTestSuite) TestRepositoryMap_AddUser() {
	expect := users.User{ID: 1, Nickname: "nickname 1", Email: "email 1"}

	s.repo.AddUser(&expect)
	got, err := s.repo.GetUserById(expect.ID)
	s.NoError(err)
	s.Equal(expect, got)
}

func (s *RepositoryMapTestSuite) TestRepositoryMap_GetUser() {
	expect := users.User{ID: 1, Nickname: "nickname 1", Email: "email 1"}
	s.repo.AddUser(&expect)

	got, err := s.repo.GetUserById(expect.ID)
	s.NoError(err)
	s.Equal(expect, got)

	_, err = s.repo.GetUserById(2)
	s.ErrorIs(err, app.IncorrectUserId)
}

func (s *RepositoryMapTestSuite) TestRepositoryMap_DeleteUser() {
	expect := users.User{ID: 1, Nickname: "nickname 1", Email: "email 1"}
	s.repo.AddUser(&expect)

	s.repo.DeleteUser(expect.ID)

	_, err := s.repo.GetUserById(expect.ID)
	s.ErrorIs(err, app.IncorrectUserId)
}

func (s *RepositoryMapTestSuite) TestRepositoryMap_ChangeUser() {
	// Test case for changing an existing user
	expected := users.User{ID: 1, Nickname: "nickname 1", Email: "email 1"}
	s.repo.AddUser(&expected)
	expected.Nickname = "Updated nickname"
	s.True(s.repo.ChangeUser(&expected))
	got, err := s.repo.GetUserById(1)
	s.NoError(err)
	s.Equal(expected, got)

	// Test case for changing a non-existing user
	expected.ID = 2
	s.False(s.repo.ChangeUser(&expected))
}

func (s *RepositoryMapTestSuite) TestRepositoryMap_GetAdsByTitle() {
	ad1 := ads.Ad{ID: 1, Title: "Ad", Text: "Ad 1 description", AuthorID: 1, Published: true}
	ad2 := ads.Ad{ID: 2, Title: "Ads", Text: "Ad 2 description", AuthorID: 1, Published: true}
	ad3 := ads.Ad{ID: 3, Title: "Another", Text: "Ad 3 description", AuthorID: 1, Published: true}
	ad4 := ads.Ad{ID: 4, Title: "All", Text: "Ad 4 description", AuthorID: 1, Published: true}

	s.repo.AddAd(&ad1)
	s.repo.AddAd(&ad2)
	s.repo.AddAd(&ad3)
	s.repo.AddAd(&ad4)

	title := "Ad"
	adsList := s.repo.GetAdsByTitle(title)
	s.Len(adsList, 2)
}

// test for checking speed processing
func BenchmarkRepoRun(b *testing.B) {
	b.Run("Get ad by id", BenchmarkGetAdById)
	b.Run("Get ads By title", BenchmarkGetAdsByTitle)
}

func BenchmarkGetAdById(b *testing.B) {
	repo := New()
	ad := &ads.Ad{ID: 1, Title: "Test ad", Text: "Test text", AuthorID: 1}
	repo.AddAd(ad)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _ = repo.GetAdById(1)
	}
}

func BenchmarkGetAdsByTitle(b *testing.B) {
	repo := New()
	ad1 := &ads.Ad{ID: 1, Title: "Test ad 1", Text: "Test text 1", AuthorID: 1}
	ad2 := &ads.Ad{ID: 2, Title: "Test ad 2", Text: "Test text 2", AuthorID: 1}
	repo.AddAd(ad1)
	repo.AddAd(ad2)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = repo.GetAdsByTitle("Test ad")
	}
}

// Fuzz tests
func FuzzGetAdById(f *testing.F) {
	repo := New()

	// Test for correct ID.
	f.Fuzz(func(t *testing.T, id int64) {
		expect := ads.Ad{ID: id, Title: "Ad 1", Text: "Ad 1 description", AuthorID: 1, Published: true}
		repo.AddAd(&expect)
		got, err := repo.GetAdById(id)

		assert.NoError(t, err)
		assert.Equal(t, got, expect)
	})
}
