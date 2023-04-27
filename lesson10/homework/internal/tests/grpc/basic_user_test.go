package grpc

import (
	"github.com/stretchr/testify/assert"
	grpcPort "homework10/internal/ports/grpc"
	"testing"
)

func TestGRPCCreateUser(t *testing.T) {
	client, ctx := getTestClient(t)

	res, err := client.CreateUser(ctx, &grpcPort.CreateUserRequest{Nickname: "Oleg", Email: "oleg@phystech.edu"})
	assert.NoError(t, err, "client.GetUser")

	assert.Equal(t, "Oleg", res.Nickname)
	assert.Equal(t, "oleg@phystech.edu", res.Email)
}

func TestGRPCUpdateUser(t *testing.T) {
	client, ctx := getTestClient(t)

	_, errUser0 := client.CreateUser(ctx, &grpcPort.CreateUserRequest{Nickname: "og buda", Email: "buda@phystech.edu"})
	assert.NoError(t, errUser0)
	_, errUser1 := client.CreateUser(ctx, &grpcPort.CreateUserRequest{Nickname: "oxxxymiron", Email: "oxxxymiron@phystech.edu"})
	assert.NoError(t, errUser1)

	response, err := client.UpdateUser(ctx, &grpcPort.UpdateUserRequest{UserId: 1, Nickname: "hello", Email: "hello@yandex.ru"})
	assert.NoError(t, err)
	assert.Equal(t, response.Nickname, "hello")
	assert.Equal(t, response.Email, "hello@yandex.ru")
}

func TestGRPCGetUser(t *testing.T) {
	client, ctx := getTestClient(t)

	_, errUser0 := client.CreateUser(ctx, &grpcPort.CreateUserRequest{Nickname: "og buda", Email: "buda@phystech.edu"})
	assert.NoError(t, errUser0)

	resp, err := client.GetUser(ctx, &grpcPort.GetUserRequest{Id: 0})
	assert.NoError(t, err)
	assert.Zero(t, resp.Id)
	assert.Equal(t, resp.Nickname, "og buda")
	assert.Equal(t, resp.Email, "buda@phystech.edu")
}

func TestGRPCDeleteUser(t *testing.T) {
	client, ctx := getTestClient(t)

	_, errUser0 := client.CreateUser(ctx, &grpcPort.CreateUserRequest{Nickname: "og buda", Email: "buda@phystech.edu"})
	assert.NoError(t, errUser0)
	_, errUser1 := client.CreateUser(ctx, &grpcPort.CreateUserRequest{Nickname: "oxxxymiron", Email: "oxxxymiron@phystech.edu"})
	assert.NoError(t, errUser1)

	_, err := client.DeleteUser(ctx, &grpcPort.DeleteUserRequest{Id: 0})
	assert.NoError(t, err)
}

func TestGRPCCreateAdByDeletedUser(t *testing.T) {
	client, ctx := getTestClient(t)

	_, errUser0 := client.CreateUser(ctx, &grpcPort.CreateUserRequest{Nickname: "og buda", Email: "buda@phystech.edu"})
	assert.NoError(t, errUser0)
	_, errUser1 := client.CreateUser(ctx, &grpcPort.CreateUserRequest{Nickname: "oxxxymiron", Email: "oxxxymiron@phystech.edu"})
	assert.NoError(t, errUser1)

	_, err := client.DeleteUser(ctx, &grpcPort.DeleteUserRequest{Id: 0})
	assert.NoError(t, err)

	_, err = client.CreateAd(ctx, &grpcPort.CreateAdRequest{UserId: 0, Title: "ok", Text: "ok"})
	assert.ErrorIs(t, err, grpcPort.ErrIncorrectUserId.Err())
}
