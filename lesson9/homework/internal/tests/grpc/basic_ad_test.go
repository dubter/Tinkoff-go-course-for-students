package grpc

import (
	"github.com/stretchr/testify/assert"
	grpcPort "homework9/internal/ports/grpc"
	"testing"
)

func TestGRPCCreateAd(t *testing.T) {
	client, ctx := getTestClient(t)

	_, err := client.CreateUser(ctx, &grpcPort.CreateUserRequest{Nickname: "Oleg", Email: "oleg@phystech.edu"})
	assert.NoError(t, err, "client.GetUser")

	resp, err := client.CreateAd(ctx, &grpcPort.CreateAdRequest{UserId: 0, Title: "hello", Text: "world"})
	assert.NoError(t, err)
	assert.Zero(t, resp.Id)
	assert.Equal(t, resp.Title, "hello")
	assert.Equal(t, resp.Text, "world")
	assert.Equal(t, resp.UserId, int64(0))
	assert.False(t, resp.Published)
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
	assert.True(t, response.Published)

	response, err = client.ChangeAdStatus(ctx, &grpcPort.ChangeAdStatusRequest{UserId: 1, AdId: response.Id, Published: false})
	assert.NoError(t, err)
	assert.False(t, response.Published)

	response, err = client.ChangeAdStatus(ctx, &grpcPort.ChangeAdStatusRequest{UserId: 1, AdId: response.Id, Published: false})
	assert.NoError(t, err)
	assert.False(t, response.Published)
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
	assert.Equal(t, response.Title, "привет")
	assert.Equal(t, response.Text, "мир")
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
	assert.Equal(t, ads.List[0].Id, publishedAd.Id)
	assert.Equal(t, ads.List[0].Title, publishedAd.Title)
	assert.Equal(t, ads.List[0].Text, publishedAd.Text)
	assert.Equal(t, ads.List[0].UserId, publishedAd.UserId)
	assert.True(t, ads.List[0].Published)
	assert.Equal(t, ads.List[0].DateCreating, publishedAd.DateCreating)
	assert.Equal(t, ads.List[0].DateUpdate, publishedAd.DateUpdate)
}

func TestGRPCAdById(t *testing.T) {
	client, ctx := getTestClient(t)

	_, errUser1 := client.CreateUser(ctx, &grpcPort.CreateUserRequest{Nickname: "vasa", Email: "oxxx@phystech.edu"})
	assert.NoError(t, errUser1)

	response, err := client.CreateAd(ctx, &grpcPort.CreateAdRequest{UserId: 0, Title: "hello", Text: "world"})
	assert.NoError(t, err)

	ad, err := client.GetAd(ctx, &grpcPort.GetAdRequest{AdId: 0})
	assert.NoError(t, err)
	assert.Equal(t, ad.Id, response.Id)
	assert.Equal(t, ad.Title, response.Title)
	assert.Equal(t, ad.Text, response.Text)
	assert.Equal(t, ad.UserId, response.UserId)
	assert.False(t, ad.Published)
	assert.Equal(t, ad.DateCreating, response.DateCreating)
	assert.Equal(t, ad.DateUpdate, response.DateUpdate)
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
	assert.Len(t, res.List, 1)

	_, err = client.DeleteAd(ctx, &grpcPort.DeleteAdRequest{AdId: response.Id, UserId: 0})
	assert.NoError(t, err)

	res, err = client.ListAdsWithFilter(ctx, &grpcPort.GetListAdsWithFilterRequest{})
	assert.NoError(t, err)
	assert.Len(t, res.List, 0)
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

	userId := &user2.Id
	published := &ad4.Published
	dateCreating := ad4.DateCreating

	response, errRes := client.ListAdsWithFilter(ctx, &grpcPort.GetListAdsWithFilterRequest{UserId: userId, Published: published, DateCreating: dateCreating})
	assert.NoError(t, errRes)
	assert.Len(t, response.List, 2)
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
	assert.Len(t, response.List, 2)
	assert.Equal(t, response.List[0].Title, "hello")
	assert.Equal(t, response.List[0].Text, "world")
	assert.Equal(t, response.List[1].Title, "hello")
	assert.Equal(t, response.List[1].Text, "friend")
}
