package tests

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateUser(t *testing.T) {
	client := getTestClient()

	response, err := client.createUser("og buda", "buda@phystech.edu")
	assert.NoError(t, err)
	assert.Zero(t, response.Data.ID)
	assert.Equal(t, response.Data.Nickname, "og buda")
	assert.Equal(t, response.Data.Email, "buda@phystech.edu")
}

func TestUpdateUser(t *testing.T) {
	client := getTestClient()

	_, errUser0 := client.createUser("og buda", "buda@phystech.edu")
	assert.NoError(t, errUser0)
	res, errUser1 := client.createUser("mayot", "mayot@phystech.edu")
	assert.NoError(t, errUser1)

	response, err := client.updateUser(res.Data.ID, "oxxximiron", "oxxximiron@phystech.edu")
	assert.NoError(t, err)
	assert.Equal(t, response.Data.Nickname, "oxxximiron")
	assert.Equal(t, response.Data.Email, "oxxximiron@phystech.edu")
}
