package grpc

import (
	"github.com/stretchr/testify/assert"
	grpcPort "homework10/internal/ports/grpc"
	"strings"
	"testing"
)

func TestGRPCCreateAd_EmptyTitle(t *testing.T) {
	client, ctx := getTestClient(t)

	_, errUser0 := client.CreateUser(ctx, &grpcPort.CreateUserRequest{Nickname: "og buda", Email: "buda@phystech.edu"})
	assert.NoError(t, errUser0)
	_, errUser1 := client.CreateUser(ctx, &grpcPort.CreateUserRequest{Nickname: "oxxxymiron", Email: "oxxxymiron@phystech.edu"})
	assert.NoError(t, errUser1)

	_, err := client.CreateAd(ctx, &grpcPort.CreateAdRequest{UserId: 1, Title: "", Text: "world"})
	assert.ErrorIs(t, err, grpcPort.ErrValidate.Err())
}

func TestGRPCCreateAd_TooLongTitle(t *testing.T) {
	client, ctx := getTestClient(t)

	_, errUser0 := client.CreateUser(ctx, &grpcPort.CreateUserRequest{Nickname: "og buda", Email: "buda@phystech.edu"})
	assert.NoError(t, errUser0)
	_, errUser1 := client.CreateUser(ctx, &grpcPort.CreateUserRequest{Nickname: "oxxxymiron", Email: "oxxxymiron@phystech.edu"})
	assert.NoError(t, errUser1)

	title := strings.Repeat("a", 101)

	_, err := client.CreateAd(ctx, &grpcPort.CreateAdRequest{UserId: 1, Title: title, Text: "world"})
	assert.ErrorIs(t, err, grpcPort.ErrValidate.Err())
}

func TestGRPCCreateAd_EmptyText(t *testing.T) {
	client, ctx := getTestClient(t)

	_, errUser0 := client.CreateUser(ctx, &grpcPort.CreateUserRequest{Nickname: "og buda", Email: "buda@phystech.edu"})
	assert.NoError(t, errUser0)
	_, errUser1 := client.CreateUser(ctx, &grpcPort.CreateUserRequest{Nickname: "oxxxymiron", Email: "oxxxymiron@phystech.edu"})
	assert.NoError(t, errUser1)

	_, err := client.CreateAd(ctx, &grpcPort.CreateAdRequest{UserId: 1, Title: "title", Text: ""})
	assert.ErrorIs(t, err, grpcPort.ErrValidate.Err())
}

func TestGRPCCreateAd_TooLongText(t *testing.T) {
	client, ctx := getTestClient(t)

	_, errUser0 := client.CreateUser(ctx, &grpcPort.CreateUserRequest{Nickname: "og buda", Email: "buda@phystech.edu"})
	assert.NoError(t, errUser0)
	_, errUser1 := client.CreateUser(ctx, &grpcPort.CreateUserRequest{Nickname: "oxxxymiron", Email: "oxxxymiron@phystech.edu"})
	assert.NoError(t, errUser1)

	text := strings.Repeat("a", 501)

	_, err := client.CreateAd(ctx, &grpcPort.CreateAdRequest{UserId: 1, Title: "title", Text: text})
	assert.ErrorIs(t, err, grpcPort.ErrValidate.Err())
}

func TestGRPCUpdateAd_EmptyTitle(t *testing.T) {
	client, ctx := getTestClient(t)

	_, errUser0 := client.CreateUser(ctx, &grpcPort.CreateUserRequest{Nickname: "og buda", Email: "buda@phystech.edu"})
	assert.NoError(t, errUser0)
	_, errUser1 := client.CreateUser(ctx, &grpcPort.CreateUserRequest{Nickname: "oxxxymiron", Email: "oxxxymiron@phystech.edu"})
	assert.NoError(t, errUser1)

	resp, err := client.CreateAd(ctx, &grpcPort.CreateAdRequest{UserId: 1, Title: "title", Text: "text"})
	assert.NoError(t, err)

	_, err = client.UpdateAd(ctx, &grpcPort.UpdateAdRequest{UserId: 1, AdId: resp.Id, Title: "", Text: "new_world"})
	assert.ErrorIs(t, err, grpcPort.ErrValidate.Err())
}

func TestGRPCUpdateAd_TooLongTitle(t *testing.T) {
	client, ctx := getTestClient(t)

	_, errUser0 := client.CreateUser(ctx, &grpcPort.CreateUserRequest{Nickname: "og buda", Email: "buda@phystech.edu"})
	assert.NoError(t, errUser0)
	_, errUser1 := client.CreateUser(ctx, &grpcPort.CreateUserRequest{Nickname: "oxxxymiron", Email: "oxxxymiron@phystech.edu"})
	assert.NoError(t, errUser1)

	resp, err := client.CreateAd(ctx, &grpcPort.CreateAdRequest{UserId: 1, Title: "title", Text: "text"})
	assert.NoError(t, err)

	title := strings.Repeat("a", 101)

	_, err = client.UpdateAd(ctx, &grpcPort.UpdateAdRequest{UserId: 1, AdId: resp.Id, Title: title, Text: "new_world"})
	assert.ErrorIs(t, err, grpcPort.ErrValidate.Err())
}

func TestGRPCUpdateAd_EmptyText(t *testing.T) {
	client, ctx := getTestClient(t)

	_, errUser0 := client.CreateUser(ctx, &grpcPort.CreateUserRequest{Nickname: "og buda", Email: "buda@phystech.edu"})
	assert.NoError(t, errUser0)
	_, errUser1 := client.CreateUser(ctx, &grpcPort.CreateUserRequest{Nickname: "oxxxymiron", Email: "oxxxymiron@phystech.edu"})
	assert.NoError(t, errUser1)

	resp, err := client.CreateAd(ctx, &grpcPort.CreateAdRequest{UserId: 1, Title: "title", Text: "text"})
	assert.NoError(t, err)

	_, err = client.UpdateAd(ctx, &grpcPort.UpdateAdRequest{UserId: 1, AdId: resp.Id, Title: "title", Text: ""})
	assert.ErrorIs(t, err, grpcPort.ErrValidate.Err())
}

func TestGRPCUpdateAd_TooLongText(t *testing.T) {
	client, ctx := getTestClient(t)

	_, errUser0 := client.CreateUser(ctx, &grpcPort.CreateUserRequest{Nickname: "og buda", Email: "buda@phystech.edu"})
	assert.NoError(t, errUser0)
	_, errUser1 := client.CreateUser(ctx, &grpcPort.CreateUserRequest{Nickname: "oxxxymiron", Email: "oxxxymiron@phystech.edu"})
	assert.NoError(t, errUser1)

	text := strings.Repeat("a", 501)

	resp, err := client.CreateAd(ctx, &grpcPort.CreateAdRequest{UserId: 1, Title: "title", Text: "text"})
	assert.NoError(t, err)

	_, err = client.UpdateAd(ctx, &grpcPort.UpdateAdRequest{UserId: 1, AdId: resp.Id, Title: "title", Text: text})
	assert.ErrorIs(t, err, grpcPort.ErrValidate.Err())
}
