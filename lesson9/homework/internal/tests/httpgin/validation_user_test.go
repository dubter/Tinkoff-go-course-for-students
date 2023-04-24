package httpgin

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateUser_EmptyNickname(t *testing.T) {
	client := getTestClient()

	_, err := client.createUser("", "buda@phystech.edu")
	assert.ErrorIs(t, err, ErrBadRequest)
}

func TestCreateUser_TooLongNickname(t *testing.T) {
	client := getTestClient()

	nickname := strings.Repeat("a", 101)
	_, err := client.createUser(nickname, "buda@phystech.edu")
	assert.ErrorIs(t, err, ErrBadRequest)
}

func TestCreateUser_EmptyEmail(t *testing.T) {
	client := getTestClient()

	_, err := client.createUser("og buda", "")
	assert.ErrorIs(t, err, ErrBadRequest)
}

func TestCreateUser_TooLongEmail(t *testing.T) {
	client := getTestClient()

	email := strings.Repeat("a", 501)
	_, err := client.createUser("og buda", email)
	assert.ErrorIs(t, err, ErrBadRequest)
}

func TestUpdateUser_EmptyNickname(t *testing.T) {
	client := getTestClient()

	_, errUser0 := client.createUser("og buda", "buda@phystech.edu")
	assert.NoError(t, errUser0)

	_, err := client.updateUser(0, "", "new_world")
	assert.ErrorIs(t, err, ErrBadRequest)
}

func TestUpdateUser_TooLongNickname(t *testing.T) {
	client := getTestClient()

	_, errUser0 := client.createUser("og buda", "buda@phystech.edu")
	assert.NoError(t, errUser0)

	nickname := strings.Repeat("a", 101)

	_, err := client.updateUser(0, nickname, "world")
	assert.ErrorIs(t, err, ErrBadRequest)
}

func TestUpdateAd_EmptyEmail(t *testing.T) {
	client := getTestClient()

	_, errUser0 := client.createUser("og buda", "buda@phystech.edu")
	assert.NoError(t, errUser0)

	_, err := client.updateUser(0, "og buda", "")
	assert.ErrorIs(t, err, ErrBadRequest)
}

func TestUpdateUser_TooLongEmail(t *testing.T) {
	client := getTestClient()

	_, errUser0 := client.createUser("og buda", "buda@phystech.edu")
	assert.NoError(t, errUser0)

	email := strings.Repeat("a", 501)

	_, err := client.updateUser(0, "og buda", email)
	assert.ErrorIs(t, err, ErrBadRequest)
}
