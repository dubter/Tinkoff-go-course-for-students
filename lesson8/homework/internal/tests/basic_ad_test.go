package tests

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const dateFormat string = "2006-01-02"

func TestCreateAd(t *testing.T) {
	client := getTestClient()

	_, errUser := client.createUser("og buda", "buda@phystech.edu")
	assert.NoError(t, errUser)

	now := time.Now()
	response, err := client.createAd(0, "hello", "world")
	assert.NoError(t, err)
	assert.Zero(t, response.Data.ID)
	assert.Equal(t, response.Data.Title, "hello")
	assert.Equal(t, response.Data.Text, "world")
	assert.Equal(t, response.Data.AuthorID, int64(0))
	assert.False(t, response.Data.Published)
	assert.Equal(t, response.Data.DateCreating.Format(dateFormat), now.Format(dateFormat))
	assert.Equal(t, response.Data.DateUpdate.Format(dateFormat), now.Format(dateFormat))
}

func TestChangeAdStatus(t *testing.T) {
	client := getTestClient()

	_, errUser0 := client.createUser("og buda", "buda@phystech.edu")
	assert.NoError(t, errUser0)
	_, errUser1 := client.createUser("oxxxymiron", "oxxxymiron@phystech.edu")
	assert.NoError(t, errUser1)

	response, err := client.createAd(1, "hello", "world")
	assert.NoError(t, err)

	response, err = client.changeAdStatus(1, response.Data.ID, true)
	assert.NoError(t, err)
	assert.True(t, response.Data.Published)

	response, err = client.changeAdStatus(1, response.Data.ID, false)
	assert.NoError(t, err)
	assert.False(t, response.Data.Published)

	response, err = client.changeAdStatus(1, response.Data.ID, false)
	assert.NoError(t, err)
	assert.False(t, response.Data.Published)
}

func TestUpdateAd(t *testing.T) {
	client := getTestClient()

	_, errUser0 := client.createUser("og buda", "buda@phystech.edu")
	assert.NoError(t, errUser0)
	_, errUser1 := client.createUser("mayot", "mayot@phystech.edu")
	assert.NoError(t, errUser1)

	response, err := client.createAd(1, "hello", "world")
	assert.NoError(t, err)

	response, err = client.updateAd(1, response.Data.ID, "привет", "мир")
	assert.NoError(t, err)
	assert.Equal(t, response.Data.Title, "привет")
	assert.Equal(t, response.Data.Text, "мир")
}

func TestListAds(t *testing.T) {
	client := getTestClient()

	_, errUser := client.createUser("og buda", "buda@phystech.edu")
	assert.NoError(t, errUser)

	now := time.Now()
	response, err := client.createAd(0, "hello", "world")
	assert.NoError(t, err)

	publishedAd, err := client.changeAdStatus(0, response.Data.ID, true)
	assert.NoError(t, err)

	_, err = client.createAd(0, "best cat", "not for sale")
	assert.NoError(t, err)

	ads, err := client.listAds()
	assert.NoError(t, err)
	assert.Len(t, ads.Data, 1)
	assert.Equal(t, ads.Data[0].ID, publishedAd.Data.ID)
	assert.Equal(t, ads.Data[0].Title, publishedAd.Data.Title)
	assert.Equal(t, ads.Data[0].Text, publishedAd.Data.Text)
	assert.Equal(t, ads.Data[0].AuthorID, publishedAd.Data.AuthorID)
	assert.True(t, ads.Data[0].Published)
	assert.Equal(t, ads.Data[0].DateCreating.Format(dateFormat), now.Format(dateFormat))
	assert.Equal(t, ads.Data[0].DateUpdate.Format(dateFormat), now.Format(dateFormat))
}

func TestAdById(t *testing.T) {
	client := getTestClient()

	_, errUser := client.createUser("og buda", "buda@phystech.edu")
	assert.NoError(t, errUser)

	now := time.Now()
	response, err := client.createAd(0, "hello", "world")
	assert.NoError(t, err)

	ad, err := client.getAdById(0)
	assert.NoError(t, err)
	assert.Equal(t, ad.Data.ID, response.Data.ID)
	assert.Equal(t, ad.Data.Title, response.Data.Title)
	assert.Equal(t, ad.Data.Text, response.Data.Text)
	assert.Equal(t, ad.Data.AuthorID, response.Data.AuthorID)
	assert.False(t, ad.Data.Published)
	assert.Equal(t, ad.Data.DateCreating.Format(dateFormat), now.Format(dateFormat))
	assert.Equal(t, ad.Data.DateUpdate.Format(dateFormat), now.Format(dateFormat))
}

func TestUpdatedAdById(t *testing.T) {
	client := getTestClient()

	_, errUser := client.createUser("og buda", "buda@phystech.edu")
	assert.NoError(t, errUser)

	now := time.Now()
	response1, err := client.createAd(0, "hello", "world")
	assert.NoError(t, err)

	response2, err := client.changeAdStatus(0, response1.Data.ID, true)
	assert.NoError(t, err)

	ad, err := client.getAdById(0)
	assert.NoError(t, err)
	assert.Equal(t, ad.Data.ID, response2.Data.ID)
	assert.Equal(t, ad.Data.Title, response2.Data.Title)
	assert.Equal(t, ad.Data.Text, response2.Data.Text)
	assert.Equal(t, ad.Data.AuthorID, response2.Data.AuthorID)
	assert.True(t, ad.Data.Published)
	assert.Equal(t, ad.Data.DateCreating.Format(dateFormat), now.Format(dateFormat))
	assert.Equal(t, ad.Data.DateUpdate.Format(dateFormat), now.Format(dateFormat))
}
