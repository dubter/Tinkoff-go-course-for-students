package grpc

import (
	"github.com/stretchr/testify/assert"
	grpcPort "homework10/internal/ports/grpc"
	"strings"
	"testing"
)

func TestGRPCCreateUser_EmptyNickname(t *testing.T) {
	client, ctx := getTestClient(t)

	_, err := client.CreateUser(ctx, &grpcPort.CreateUserRequest{Nickname: "", Email: "buda@phystech.edu"})
	assert.ErrorIs(t, err, grpcPort.ErrValidate.Err())
}

func TestGRPCCreateUser_TooLongNickname(t *testing.T) {
	client, ctx := getTestClient(t)

	nickname := strings.Repeat("a", 101)
	_, err := client.CreateUser(ctx, &grpcPort.CreateUserRequest{Nickname: nickname, Email: "buda@phystech.edu"})
	assert.ErrorIs(t, err, grpcPort.ErrValidate.Err())
}

func TestGRPCCreateUser_EmptyEmail(t *testing.T) {
	client, ctx := getTestClient(t)

	_, err := client.CreateUser(ctx, &grpcPort.CreateUserRequest{Nickname: "nickname", Email: ""})
	assert.ErrorIs(t, err, grpcPort.ErrValidate.Err())
}

func TestGRPCCreateUser_TooLongEmail(t *testing.T) {
	client, ctx := getTestClient(t)

	email := strings.Repeat("a", 501)
	_, err := client.CreateUser(ctx, &grpcPort.CreateUserRequest{Nickname: "nickname", Email: email})
	assert.ErrorIs(t, err, grpcPort.ErrValidate.Err())
}

func TestGRPCUpdateUser_EmptyNickname(t *testing.T) {
	client, ctx := getTestClient(t)

	_, err := client.CreateUser(ctx, &grpcPort.CreateUserRequest{Nickname: "nickname", Email: "email"})
	assert.NoError(t, err)

	_, err = client.UpdateUser(ctx, &grpcPort.UpdateUserRequest{UserId: 0, Nickname: "", Email: "new_world"})
	assert.ErrorIs(t, err, grpcPort.ErrValidate.Err())
}

func TestGRPCUpdateUser_TooLongNickname(t *testing.T) {
	client, ctx := getTestClient(t)

	_, err := client.CreateUser(ctx, &grpcPort.CreateUserRequest{Nickname: "nickname", Email: "email"})
	assert.NoError(t, err)

	nickname := strings.Repeat("a", 101)

	_, err = client.UpdateUser(ctx, &grpcPort.UpdateUserRequest{UserId: 0, Nickname: nickname, Email: "world"})
	assert.ErrorIs(t, err, grpcPort.ErrValidate.Err())
}

func TestGRPCUpdateAd_EmptyEmail(t *testing.T) {
	client, ctx := getTestClient(t)

	_, err := client.CreateUser(ctx, &grpcPort.CreateUserRequest{Nickname: "nickname", Email: "email"})
	assert.NoError(t, err)

	_, err = client.UpdateUser(ctx, &grpcPort.UpdateUserRequest{UserId: 0, Nickname: "nickname", Email: ""})
	assert.ErrorIs(t, err, grpcPort.ErrValidate.Err())
}

func TestGRPCUpdateUser_TooLongEmail(t *testing.T) {
	client, ctx := getTestClient(t)

	_, err := client.CreateUser(ctx, &grpcPort.CreateUserRequest{Nickname: "nickname", Email: "email"})
	assert.NoError(t, err)

	email := strings.Repeat("a", 501)

	_, err = client.UpdateUser(ctx, &grpcPort.UpdateUserRequest{UserId: 0, Nickname: "nickname", Email: email})
	assert.ErrorIs(t, err, grpcPort.ErrValidate.Err())
}
