package httpgin

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetAdsByTitle(t *testing.T) {
	client := getTestClient()

	_, errUser := client.createUser("og buda", "buda@phystech.edu")
	assert.NoError(t, errUser)

	_, err := client.createAd(0, "hello", "world")
	assert.NoError(t, err)

	_, err = client.createAd(0, "bye", "forever")
	assert.NoError(t, err)

	_, err = client.createAd(0, "hello", "friend")
	assert.NoError(t, err)

	// "hel" is a prefix of "hello"
	response, err := client.listAdsByTitle("hel")
	assert.NoError(t, err)
	assert.Len(t, response.Data, 2)
	assert.Equal(t, response.Data[0].Title, "hello")
	assert.Equal(t, response.Data[0].Text, "world")
	assert.Equal(t, response.Data[1].Title, "hello")
	assert.Equal(t, response.Data[1].Text, "friend")
}

func TestGetFilteredAdsOnlyUnpublished(t *testing.T) {
	client := getTestClient()

	_, errUser := client.createUser("og buda", "buda@phystech.edu")
	assert.NoError(t, errUser)

	_, err := client.createAd(0, "hello", "world")
	assert.NoError(t, err)

	ad2, err := client.createAd(0, "bye", "forever")
	assert.NoError(t, err)

	_, err = client.createAd(0, "hello", "friend")
	assert.NoError(t, err)

	_, err = client.changeAdStatus(0, ad2.Data.ID, true)
	assert.NoError(t, err)

	filters := map[string]any{
		"published": false,
	}
	response, errRes := client.getListAdsWithFilter(filters)
	assert.NoError(t, errRes)
	assert.Len(t, response.Data, 2)
}

func TestGetFilteredAdsByAuthor(t *testing.T) {
	client := getTestClient()

	user1, errUser := client.createUser("og buda", "buda@phystech.edu")
	assert.NoError(t, errUser)

	user2, errUser1 := client.createUser("oxxxymiron", "oxxxymiron@phystech.edu")
	assert.NoError(t, errUser1)

	_, err := client.createAd(user1.Data.ID, "hello", "world")
	assert.NoError(t, err)

	_, err = client.createAd(user1.Data.ID, "hello", "friend")
	assert.NoError(t, err)

	_, err = client.createAd(user2.Data.ID, "bye", "forever")
	assert.NoError(t, err)

	filters := map[string]any{
		"user_id": user1.Data.ID,
	}
	response, errRes := client.getListAdsWithFilter(filters)
	assert.NoError(t, errRes)
	assert.Len(t, response.Data, 2)
}

func TestGetFilteredAdsByDateCreating(t *testing.T) {
	client := getTestClient()

	user1, errUser := client.createUser("og buda", "buda@phystech.edu")
	assert.NoError(t, errUser)

	ad1, err := client.createAd(user1.Data.ID, "hello", "world")
	assert.NoError(t, err)

	_, err = client.createAd(user1.Data.ID, "hello", "friend")
	assert.NoError(t, err)

	_, err = client.createAd(user1.Data.ID, "bye", "forever")
	assert.NoError(t, err)

	filters := map[string]any{
		"date_creating": ad1.Data.DateCreating,
	}
	response, errRes := client.getListAdsWithFilter(filters)
	assert.NoError(t, errRes)
	assert.Len(t, response.Data, 3)
}

// Let's combine some filters
func TestGetAdsByOnlyUnpublishedByAuthor(t *testing.T) {
	client := getTestClient()

	user1, errUser := client.createUser("og buda", "buda@phystech.edu")
	assert.NoError(t, errUser)

	user2, errUser1 := client.createUser("oxxxymiron", "oxxxymiron@phystech.edu")
	assert.NoError(t, errUser1)

	ad1, err := client.createAd(user1.Data.ID, "hello", "world")
	assert.NoError(t, err)

	ad2, err := client.createAd(user1.Data.ID, "hello", "friend")
	assert.NoError(t, err)

	_, err = client.createAd(user2.Data.ID, "bye", "forever")
	assert.NoError(t, err)

	_, err = client.changeAdStatus(user1.Data.ID, ad1.Data.ID, true)
	assert.NoError(t, err)

	filters := map[string]any{
		"published": false,
		"user_id":   user1.Data.ID,
	}
	response, errRes := client.getListAdsWithFilter(filters)
	assert.NoError(t, errRes)
	assert.Len(t, response.Data, 1)
	assert.Equal(t, response.Data[0], ad2.Data)
}

func TestGetAdsByOnlyUnpublishedByDateCreatingByAuthor(t *testing.T) {
	client := getTestClient()

	user1, errUser := client.createUser("og buda", "buda@phystech.edu")
	assert.NoError(t, errUser)

	user2, errUser1 := client.createUser("oxxxymiron", "oxxxymiron@phystech.edu")
	assert.NoError(t, errUser1)

	_, err := client.createAd(user1.Data.ID, "hello", "world")
	assert.NoError(t, err)

	ad2, err := client.createAd(user2.Data.ID, "hello", "friend")
	assert.NoError(t, err)

	_, err = client.createAd(user2.Data.ID, "bye", "forever")
	assert.NoError(t, err)

	ad4, err := client.createAd(user2.Data.ID, "good", "evening")
	assert.NoError(t, err)

	_, err = client.changeAdStatus(user2.Data.ID, ad2.Data.ID, true)
	assert.NoError(t, err)

	filters := map[string]any{
		"published":     false,
		"user_id":       user2.Data.ID,
		"date_creating": ad4.Data.DateCreating,
	}
	response, errRes := client.getListAdsWithFilter(filters)
	assert.NoError(t, errRes)
	assert.Len(t, response.Data, 2)
}
