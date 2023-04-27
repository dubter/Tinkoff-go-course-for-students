package grpc

import (
	"github.com/stretchr/testify/assert"
	grpcPort "homework10/internal/ports/grpc"
	"testing"
)

func TestGRPCCreateAd(t *testing.T) {
	client, ctx := getTestClient(t)

	_, err := client.CreateUser(ctx, &grpcPort.CreateUserRequest{Nickname: "Oleg", Email: "oleg@phystech.edu"})
	assert.NoError(t, err, "client.GetUser")

	resp, err := client.CreateAd(ctx, &grpcPort.CreateAdRequest{UserId: 0, Title: "hello", Text: "world"})
	assert.NoError(t, err)
	assert.Zero(t, resp.Id)
	assert.Equal(t, resp.GetTitle(), "hello")
	assert.Equal(t, resp.GetText(), "world")
	assert.Equal(t, resp.GetUserId(), int64(0))
	assert.False(t, resp.GetPublished())
}

func TestGRPCChangeAdStatus(t *testing.T) {
	client, ctx := getTestClient(t)

	_, errUser0 := client.CreateUser(ctx, &grpcPort.CreateUserRequest{Nickname: "peter", Email: "buda@phystech.edu"})
	assert.NoError(t, errUser0)
	_, errUser1 := client.CreateUser(ctx, &grpcPort.CreateUserRequest{Nickname: "vasa", Email: "oxxx@phystech.edu"})
	assert.NoError(t, errUser1)

	response, err := client.CreateAd(ctx, &grpcPort.CreateAdRequest{UserId: 1, Title: "hello", Text: "world"})
	assert.NoError(t, err)

	response, err = client.ChangeAdStatus(ctx, &grpcPort.ChangeAdStatusRequest{UserId: 1, AdId: response.Id, Published: true})
	assert.NoError(t, err)
	assert.True(t, response.GetPublished())

	response, err = client.ChangeAdStatus(ctx, &grpcPort.ChangeAdStatusRequest{UserId: 1, AdId: response.Id, Published: false})
	assert.NoError(t, err)
	assert.False(t, response.GetPublished())

	response, err = client.ChangeAdStatus(ctx, &grpcPort.ChangeAdStatusRequest{UserId: 1, AdId: response.Id, Published: false})
	assert.NoError(t, err)
	assert.False(t, response.GetPublished())
}

func TestGRPCUpdateAd(t *testing.T) {
	client, ctx := getTestClient(t)

	_, errUser0 := client.CreateUser(ctx, &grpcPort.CreateUserRequest{Nickname: "peter", Email: "buda@phystech.edu"})
	assert.NoError(t, errUser0)
	_, errUser1 := client.CreateUser(ctx, &grpcPort.CreateUserRequest{Nickname: "vasa", Email: "oxxx@phystech.edu"})
	assert.NoError(t, errUser1)

	response, err := client.CreateAd(ctx, &grpcPort.CreateAdRequest{UserId: 1, Title: "hello", Text: "world"})
	assert.NoError(t, err)

	response, err = client.UpdateAd(ctx, &grpcPort.UpdateAdRequest{UserId: 1, AdId: response.Id, Title: "привет", Text: "мир"})
	assert.NoError(t, err)
	assert.Equal(t, response.GetTitle(), "привет")
	assert.Equal(t, response.GetText(), "мир")
}

func TestGRPCUpdateAdNotFound(t *testing.T) {
	client, ctx := getTestClient(t)

	_, errUser0 := client.CreateUser(ctx, &grpcPort.CreateUserRequest{Nickname: "peter", Email: "buda@phystech.edu"})
	assert.NoError(t, errUser0)

	response, err := client.CreateAd(ctx, &grpcPort.CreateAdRequest{UserId: 0, Title: "hello", Text: "world"})
	assert.NoError(t, err)

	response, err = client.UpdateAd(ctx, &grpcPort.UpdateAdRequest{UserId: 100, AdId: response.Id, Title: "привет", Text: "мир"})
	assert.ErrorIs(t, err, grpcPort.ErrIncorrectUserId.Err())
}

func TestGRPCListAds(t *testing.T) {
	client, ctx := getTestClient(t)

	_, errUser1 := client.CreateUser(ctx, &grpcPort.CreateUserRequest{Nickname: "vasa", Email: "oxxx@phystech.edu"})
	assert.NoError(t, errUser1)

	response, err := client.CreateAd(ctx, &grpcPort.CreateAdRequest{UserId: 0, Title: "hello", Text: "world"})
	assert.NoError(t, err)

	publishedAd, err := client.ChangeAdStatus(ctx, &grpcPort.ChangeAdStatusRequest{UserId: 0, AdId: response.Id, Published: true})
	assert.NoError(t, err)

	_, err = client.CreateAd(ctx, &grpcPort.CreateAdRequest{UserId: 0, Title: "bye", Text: "bro"})
	assert.NoError(t, err)

	ads, err := client.ListAdsWithFilter(ctx, &grpcPort.GetListAdsWithFilterRequest{})
	assert.NoError(t, err)
	assert.Len(t, ads.List, 1)
	assert.Equal(t, ads.List[0].GetId(), publishedAd.GetId())
	assert.Equal(t, ads.List[0].GetTitle(), publishedAd.GetTitle())
	assert.Equal(t, ads.List[0].GetText(), publishedAd.GetText())
	assert.Equal(t, ads.List[0].GetUserId(), publishedAd.GetUserId())
	assert.True(t, ads.List[0].GetPublished())
	assert.Equal(t, ads.List[0].GetDateCreating(), publishedAd.GetDateCreating())
	assert.Equal(t, ads.List[0].GetDateUpdate(), publishedAd.GetDateUpdate())
}

func TestGRPCAdById(t *testing.T) {
	client, ctx := getTestClient(t)

	_, errUser1 := client.CreateUser(ctx, &grpcPort.CreateUserRequest{Nickname: "vasa", Email: "oxxx@phystech.edu"})
	assert.NoError(t, errUser1)

	response, err := client.CreateAd(ctx, &grpcPort.CreateAdRequest{UserId: 0, Title: "hello", Text: "world"})
	assert.NoError(t, err)

	ad, err := client.GetAd(ctx, &grpcPort.GetAdRequest{AdId: 0})
	assert.NoError(t, err)
	assert.Equal(t, ad.GetId(), response.GetId())
	assert.Equal(t, ad.GetTitle(), response.GetTitle())
	assert.Equal(t, ad.GetText(), response.GetText())
	assert.Equal(t, ad.GetUserId(), response.GetUserId())
	assert.False(t, ad.GetPublished())
	assert.Equal(t, ad.GetDateCreating(), response.GetDateCreating())
	assert.Equal(t, ad.GetDateUpdate(), response.GetDateUpdate())
}

func TestGRPCDeleteAdById(t *testing.T) {
	client, ctx := getTestClient(t)

	_, errUser1 := client.CreateUser(ctx, &grpcPort.CreateUserRequest{Nickname: "vasa", Email: "oxxx@phystech.edu"})
	assert.NoError(t, errUser1)

	response, err := client.CreateAd(ctx, &grpcPort.CreateAdRequest{UserId: 0, Title: "hello", Text: "world"})
	assert.NoError(t, err)

	_, err = client.ChangeAdStatus(ctx, &grpcPort.ChangeAdStatusRequest{UserId: 0, AdId: response.Id, Published: true})
	assert.NoError(t, err)

	res, err := client.ListAdsWithFilter(ctx, &grpcPort.GetListAdsWithFilterRequest{})
	assert.NoError(t, err)
	assert.Len(t, res.GetList(), 1)

	_, err = client.DeleteAd(ctx, &grpcPort.DeleteAdRequest{AdId: response.Id, UserId: 0})
	assert.NoError(t, err)

	res, err = client.ListAdsWithFilter(ctx, &grpcPort.GetListAdsWithFilterRequest{})
	assert.NoError(t, err)
	assert.Len(t, res.GetList(), 0)
}

func TestGRPCGetAdsByOnlyUnpublishedByDateCreatingByAuthor(t *testing.T) {
	client, ctx := getTestClient(t)

	user1, errUser0 := client.CreateUser(ctx, &grpcPort.CreateUserRequest{Nickname: "og buda", Email: "buda@phystech.edu"})
	assert.NoError(t, errUser0)
	user2, errUser1 := client.CreateUser(ctx, &grpcPort.CreateUserRequest{Nickname: "oxxxymiron", Email: "oxxxymiron@phystech.edu"})
	assert.NoError(t, errUser1)

	_, err := client.CreateAd(ctx, &grpcPort.CreateAdRequest{UserId: user1.Id, Title: "hello", Text: "world"})
	assert.NoError(t, err)

	ad2, err := client.CreateAd(ctx, &grpcPort.CreateAdRequest{UserId: user2.Id, Title: "hello", Text: "friend"})
	assert.NoError(t, err)

	_, err = client.CreateAd(ctx, &grpcPort.CreateAdRequest{UserId: user2.Id, Title: "bye", Text: "forever"})
	assert.NoError(t, err)

	ad4, err := client.CreateAd(ctx, &grpcPort.CreateAdRequest{UserId: user2.Id, Title: "good", Text: "evening"})
	assert.NoError(t, err)

	_, err = client.ChangeAdStatus(ctx, &grpcPort.ChangeAdStatusRequest{UserId: user2.Id, AdId: ad2.Id, Published: true})
	assert.NoError(t, err)

	userId := user2.GetId()
	published := ad4.GetPublished()
	dateCreating := ad4.GetDateCreating()

	response, errRes := client.ListAdsWithFilter(ctx, &grpcPort.GetListAdsWithFilterRequest{UserId: &userId, Published: &published, DateCreating: dateCreating})
	assert.NoError(t, errRes)
	assert.Len(t, response.GetList(), 2)
}

func TestGRPCGetAdsByTitle(t *testing.T) {
	client, ctx := getTestClient(t)

	user1, errUser0 := client.CreateUser(ctx, &grpcPort.CreateUserRequest{Nickname: "og buda", Email: "buda@phystech.edu"})
	assert.NoError(t, errUser0)

	_, err := client.CreateAd(ctx, &grpcPort.CreateAdRequest{UserId: user1.Id, Title: "hello", Text: "world"})
	assert.NoError(t, err)

	_, err = client.CreateAd(ctx, &grpcPort.CreateAdRequest{UserId: user1.Id, Title: "hello", Text: "friend"})
	assert.NoError(t, err)

	_, err = client.CreateAd(ctx, &grpcPort.CreateAdRequest{UserId: user1.Id, Title: "bye", Text: "forever"})
	assert.NoError(t, err)

	// "hel" is a prefix of "hello"
	response, err := client.ListAdsByTitle(ctx, &grpcPort.GetListAdsByTitleRequest{Title: "hel"})
	assert.NoError(t, err)
	assert.Len(t, response.GetList(), 2)
	assert.Equal(t, response.GetList()[0].GetTitle(), "hello")
	assert.Equal(t, response.GetList()[0].GetText(), "world")
	assert.Equal(t, response.GetList()[1].GetTitle(), "hello")
	assert.Equal(t, response.GetList()[1].GetText(), "friend")
}

func TestGRPCGetAdByIncorrectId(t *testing.T) {
	client, ctx := getTestClient(t)

	_, errUser0 := client.CreateUser(ctx, &grpcPort.CreateUserRequest{Nickname: "peter", Email: "buda@phystech.edu"})
	assert.NoError(t, errUser0)

	_, err := client.CreateAd(ctx, &grpcPort.CreateAdRequest{UserId: 0, Title: "hello", Text: "world"})
	assert.NoError(t, err)

	_, err = client.GetAd(ctx, &grpcPort.GetAdRequest{AdId: 100})
	assert.ErrorIs(t, err, grpcPort.ErrIncorrectAdId.Err())
}
