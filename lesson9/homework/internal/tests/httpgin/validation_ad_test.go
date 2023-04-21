package httpgin

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateAd_EmptyTitle(t *testing.T) {
	client := getTestClient()

	_, errUser0 := client.createUser("og buda", "buda@phystech.edu")
	assert.NoError(t, errUser0)
	_, errUser1 := client.createUser("oxxxymiron", "oxxxymiron@phystech.edu")
	assert.NoError(t, errUser1)

	_, err := client.createAd(1, "", "world")
	assert.ErrorIs(t, err, ErrBadRequest)
}

func TestCreateAd_TooLongTitle(t *testing.T) {
	client := getTestClient()

	_, errUser0 := client.createUser("og buda", "buda@phystech.edu")
	assert.NoError(t, errUser0)
	_, errUser1 := client.createUser("oxxxymiron", "oxxxymiron@phystech.edu")
	assert.NoError(t, errUser1)

	title := strings.Repeat("a", 101)

	_, err := client.createAd(1, title, "world")
	assert.ErrorIs(t, err, ErrBadRequest)
}

func TestCreateAd_EmptyText(t *testing.T) {
	client := getTestClient()

	_, errUser0 := client.createUser("og buda", "buda@phystech.edu")
	assert.NoError(t, errUser0)
	_, errUser1 := client.createUser("oxxxymiron", "oxxxymiron@phystech.edu")
	assert.NoError(t, errUser1)

	_, err := client.createAd(1, "title", "")
	assert.ErrorIs(t, err, ErrBadRequest)
}

func TestCreateAd_TooLongText(t *testing.T) {
	client := getTestClient()

	_, errUser0 := client.createUser("og buda", "buda@phystech.edu")
	assert.NoError(t, errUser0)
	_, errUser1 := client.createUser("oxxxymiron", "oxxxymiron@phystech.edu")
	assert.NoError(t, errUser1)

	text := strings.Repeat("a", 501)

	_, err := client.createAd(1, "title", text)
	assert.ErrorIs(t, err, ErrBadRequest)
}

func TestUpdateAd_EmptyTitle(t *testing.T) {
	client := getTestClient()

	_, errUser0 := client.createUser("og buda", "buda@phystech.edu")
	assert.NoError(t, errUser0)
	_, errUser1 := client.createUser("oxxxymiron", "oxxxymiron@phystech.edu")
	assert.NoError(t, errUser1)

	resp, err := client.createAd(1, "hello", "world")
	assert.NoError(t, err)

	_, err = client.updateAd(1, resp.Data.ID, "", "new_world")
	assert.ErrorIs(t, err, ErrBadRequest)
}

func TestUpdateAd_TooLongTitle(t *testing.T) {
	client := getTestClient()

	_, errUser0 := client.createUser("og buda", "buda@phystech.edu")
	assert.NoError(t, errUser0)
	_, errUser1 := client.createUser("oxxxymiron", "oxxxymiron@phystech.edu")
	assert.NoError(t, errUser1)

	resp, err := client.createAd(1, "hello", "world")
	assert.NoError(t, err)

	title := strings.Repeat("a", 101)

	_, err = client.updateAd(1, resp.Data.ID, title, "world")
	assert.ErrorIs(t, err, ErrBadRequest)
}

func TestUpdateAd_EmptyText(t *testing.T) {
	client := getTestClient()

	_, errUser0 := client.createUser("og buda", "buda@phystech.edu")
	assert.NoError(t, errUser0)
	_, errUser1 := client.createUser("oxxxymiron", "oxxxymiron@phystech.edu")
	assert.NoError(t, errUser1)

	resp, err := client.createAd(1, "hello", "world")
	assert.NoError(t, err)

	_, err = client.updateAd(1, resp.Data.ID, "title", "")
	assert.ErrorIs(t, err, ErrBadRequest)
}

func TestUpdateAd_TooLongText(t *testing.T) {
	client := getTestClient()

	_, errUser0 := client.createUser("og buda", "buda@phystech.edu")
	assert.NoError(t, errUser0)
	_, errUser1 := client.createUser("oxxxymiron", "oxxxymiron@phystech.edu")
	assert.NoError(t, errUser1)

	text := strings.Repeat("a", 501)

	resp, err := client.createAd(1, "hello", "world")
	assert.NoError(t, err)

	_, err = client.updateAd(1, resp.Data.ID, "title", text)
	assert.ErrorIs(t, err, ErrBadRequest)
}
