package tests

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUpdateNonexistentUser(t *testing.T) {
	client := getTestClient()

	_, errUser0 := client.createUser("og buda", "buda@phystech.edu")
	assert.NoError(t, errUser0)
	_, errUser1 := client.createUser("oxxxymiron", "oxxxymiron@phystech.edu")
	assert.NoError(t, errUser1)

	_, err := client.updateUser(123, "mayot", "mayot@phystech.edu")
	assert.ErrorIs(t, err, ErrForbidden)
}

func TestCreateAdOfNonexistentUser(t *testing.T) {
	client := getTestClient()

	_, err := client.createAd(123, "hello", "world")
	assert.ErrorIs(t, err, ErrForbidden)
}
