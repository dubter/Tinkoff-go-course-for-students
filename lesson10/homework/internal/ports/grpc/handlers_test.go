package grpc

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/suite"
	"google.golang.org/protobuf/types/known/timestamppb"
	"homework10/internal/ads"
	"homework10/internal/app"
	"homework10/internal/ports/grpc/proto"
	"homework10/internal/ports/httpgin/mocks"
	"homework10/internal/users"
	"testing"
)

type AdServiceTestSuite struct {
	suite.Suite
	app mocks.App
}

func (s *AdServiceTestSuite) SetupTest() {
	s.app = mocks.App{}
}

func TestRepoRun(t *testing.T) {
	suite.Run(t, new(AdServiceTestSuite))
}

func (s *AdServiceTestSuite) TestAdService_CreateAd() {
	request := &proto.CreateAdRequest{Title: "title 1", Text: "text 1", UserId: 1}
	expect := &ads.Ad{ID: 1, Title: "title 1", Text: "text 1", AuthorID: 1}
	s.app.On("CreateAd", request.Title, request.Text, request.UserId).Return(expect, nil)

	service := NewService(&s.app)
	response, err := service.CreateAd(context.TODO(), request)
	s.NoError(err)
	s.Equal(response, AdSuccessResponse(expect))
}

func (s *AdServiceTestSuite) TestAdService_CreateAdValidationErr() {
	request := &proto.CreateAdRequest{Title: "", Text: "text 1", UserId: 1}
	expect := &ads.Ad{ID: 1, Title: "", Text: "text 1", AuthorID: 1}
	s.app.On("CreateAd", request.Title, request.Text, request.UserId).Return(expect, app.ValidateError)

	service := NewService(&s.app)
	_, err := service.CreateAd(context.TODO(), request)
	s.ErrorIs(err, ErrValidate.Err())
}

func (s *AdServiceTestSuite) TestAdService_CreateAdIncorrectUserId() {
	request := &proto.CreateAdRequest{Title: "title", Text: "text 1", UserId: 3}
	expect := &ads.Ad{ID: 1, Title: "title", Text: "text 1", AuthorID: 3}
	s.app.On("CreateAd", request.Title, request.Text, request.UserId).Return(expect, app.IncorrectUserId)

	service := NewService(&s.app)
	_, err := service.CreateAd(context.TODO(), request)
	s.ErrorIs(err, ErrIncorrectUserId.Err())
}

func (s *AdServiceTestSuite) TestAdService_ChangeAdStatus() {
	request := &proto.ChangeAdStatusRequest{AdId: 1, Published: true, UserId: 1}
	expect := &ads.Ad{ID: 1, Title: "ad 1", Text: "text 1", AuthorID: 1, Published: true}
	s.app.On("ChangeAdStatus", request.AdId, request.UserId, request.Published).Return(expect, nil)

	service := NewService(&s.app)
	response, err := service.ChangeAdStatus(context.TODO(), request)
	s.NoError(err)
	s.Equal(response, AdSuccessResponse(expect))
}

func (s *AdServiceTestSuite) TestAdService_ChangeAdStatusIncorrectUserId() {
	request := &proto.ChangeAdStatusRequest{AdId: 1, Published: true, UserId: 10}
	expect := &ads.Ad{ID: 1, Title: "ad 1", Text: "text 1", AuthorID: 10, Published: true}
	s.app.On("ChangeAdStatus", request.AdId, request.UserId, request.Published).Return(expect, app.IncorrectUserId)

	service := NewService(&s.app)
	_, err := service.ChangeAdStatus(context.TODO(), request)
	s.ErrorIs(err, ErrIncorrectUserId.Err())
}

func (s *AdServiceTestSuite) TestAdService_UpdateAd() {
	request := &proto.UpdateAdRequest{AdId: 1, Title: "updated ad", UserId: 1, Text: "updated text"}
	expect := &ads.Ad{ID: 1, Title: "updated ad", Text: "updated text", AuthorID: 1, Published: false}
	s.app.On("UpdateAd", request.AdId, request.UserId, request.Title, request.Text).Return(expect, nil)

	service := NewService(&s.app)
	response, err := service.UpdateAd(context.TODO(), request)
	s.NoError(err)
	s.Equal(response, AdSuccessResponse(expect))
}

func (s *AdServiceTestSuite) TestAdService_UpdateAdValidationErr() {
	request := &proto.UpdateAdRequest{AdId: 1, Title: "", UserId: 1, Text: "updated text"}
	expect := &ads.Ad{ID: 1, Title: "updated ad", Text: "", AuthorID: 1, Published: false}
	s.app.On("UpdateAd", request.AdId, request.UserId, request.Title, request.Text).Return(expect, app.ValidateError)

	service := NewService(&s.app)
	_, err := service.UpdateAd(context.TODO(), request)
	s.ErrorIs(err, ErrValidate.Err())
}

func (s *AdServiceTestSuite) TestAdService_UpdateAdIncorrectUserId() {
	request := &proto.UpdateAdRequest{AdId: 1, Title: "", UserId: 1, Text: "updated text"}
	expect := &ads.Ad{ID: 1, Title: "updated ad", Text: "", AuthorID: 1, Published: false}
	s.app.On("UpdateAd", request.AdId, request.UserId, request.Title, request.Text).Return(expect, app.IncorrectUserId)

	service := NewService(&s.app)
	_, err := service.UpdateAd(context.TODO(), request)
	s.ErrorIs(err, ErrIncorrectUserId.Err())
}

func (s *AdServiceTestSuite) TestAdService_ListAdsWithFilter() {
	userId := int64(1)
	published := true
	var time timestamppb.Timestamp
	request := &proto.GetListAdsWithFilterRequest{UserId: &userId, Published: &published, DateCreating: &time}

	expect1 := ads.Ad{ID: 1, Title: "title 1", Text: "text 1", AuthorID: 1, Published: true}
	expect2 := ads.Ad{ID: 2, Title: "title 2", Text: "text 2", AuthorID: 1, Published: true}
	adsList := []ads.Ad{expect1, expect2}

	filters := map[string]any{"user_id": *request.UserId, "published": *request.Published, "date_creating": fmt.Sprint((*request.DateCreating).AsTime().UTC())}
	s.app.On("GetListAds", filters).Return(adsList, nil)

	service := NewService(&s.app)
	response, err := service.ListAdsWithFilter(context.TODO(), request)
	s.NoError(err)
	s.Equal(response, AdsSuccessResponse(adsList))
}

func (s *AdServiceTestSuite) TestAdService_ListAdsByTitle() {
	expect1 := ads.Ad{ID: 1, Title: "title", Text: "text 1", AuthorID: 1, Published: true}
	expect2 := ads.Ad{ID: 2, Title: "title", Text: "text 2", AuthorID: 1, Published: true}
	adsList := []ads.Ad{expect1, expect2}

	title := expect1.Title
	request := &proto.GetListAdsByTitleRequest{Title: title}

	s.app.On("GetListAdsByTitle", title).Return(adsList, nil)

	service := NewService(&s.app)
	response, err := service.ListAdsByTitle(context.TODO(), request)
	s.NoError(err)
	s.Equal(response, AdsSuccessResponse(adsList))
}

func (s *AdServiceTestSuite) TestAdService_CreateUser() {
	expect := &users.User{ID: 1, Nickname: "nickname", Email: "email"}
	request := &proto.CreateUserRequest{Nickname: expect.Nickname, Email: expect.Email}

	s.app.On("CreateUser", request.Nickname, request.Email).Return(expect, nil)

	service := NewService(&s.app)
	response, err := service.CreateUser(context.TODO(), request)
	s.NoError(err)
	s.Equal(response, UserSuccessResponse(expect))
}

func (s *AdServiceTestSuite) TestAdService_CreateUserValidationErr() {
	expect := &users.User{ID: 1, Nickname: "", Email: "email"}
	request := &proto.CreateUserRequest{Nickname: expect.Nickname, Email: expect.Email}

	s.app.On("CreateUser", request.Nickname, request.Email).Return(expect, app.ValidateError)

	service := NewService(&s.app)
	_, err := service.CreateUser(context.TODO(), request)
	s.ErrorIs(err, ErrValidate.Err())
}

func (s *AdServiceTestSuite) TestAdService_UpdateUser() {
	expect := &users.User{ID: 1, Nickname: "nickname", Email: "email"}
	request := &proto.UpdateUserRequest{UserId: expect.ID, Nickname: expect.Nickname, Email: expect.Email}

	s.app.On("UpdateUser", request.UserId, request.Nickname, request.Email).Return(expect, nil)

	service := NewService(&s.app)
	response, err := service.UpdateUser(context.TODO(), request)
	s.NoError(err)
	s.Equal(response, UserSuccessResponse(expect))
}

func (s *AdServiceTestSuite) TestAdService_UpdateUserIncorrectUserId() {
	expect := &users.User{ID: 10, Nickname: "nickname", Email: "email"}
	request := &proto.UpdateUserRequest{UserId: expect.ID, Nickname: expect.Nickname, Email: expect.Email}

	s.app.On("UpdateUser", request.UserId, request.Nickname, request.Email).Return(expect, app.IncorrectUserId)

	service := NewService(&s.app)
	_, err := service.UpdateUser(context.TODO(), request)
	s.ErrorIs(err, ErrIncorrectUserId.Err())
}

func (s *AdServiceTestSuite) TestAdService_UpdateUserValidationErr() {
	expect := &users.User{ID: 10, Nickname: "", Email: "email"}
	request := &proto.UpdateUserRequest{UserId: expect.ID, Nickname: expect.Nickname, Email: expect.Email}

	s.app.On("UpdateUser", request.UserId, request.Nickname, request.Email).Return(expect, app.ValidateError)

	service := NewService(&s.app)
	_, err := service.UpdateUser(context.TODO(), request)
	s.ErrorIs(err, ErrValidate.Err())
}

func (s *AdServiceTestSuite) TestAdService_GetUser() {
	expect := &users.User{ID: 1, Nickname: "nickname", Email: "email"}
	request := &proto.GetUserRequest{Id: expect.ID}

	s.app.On("GetUser", request.Id).Return(expect, nil)

	service := NewService(&s.app)
	response, err := service.GetUser(context.TODO(), request)
	s.NoError(err)
	s.Equal(response, UserSuccessResponse(expect))
}

func (s *AdServiceTestSuite) TestAdService_GetUserIncorrectUserId() {
	expect := &users.User{ID: 1, Nickname: "nickname", Email: "email"}
	request := &proto.GetUserRequest{Id: expect.ID}

	s.app.On("GetUser", request.Id).Return(expect, app.IncorrectUserId)

	service := NewService(&s.app)
	_, err := service.GetUser(context.TODO(), request)
	s.ErrorIs(err, ErrIncorrectUserId.Err())
}

func (s *AdServiceTestSuite) TestAdService_DeleteUser() {
	request := &proto.DeleteUserRequest{Id: 1}

	s.app.On("DeleteUser", request.Id).Return(nil)

	service := NewService(&s.app)
	_, err := service.DeleteUser(context.TODO(), request)
	s.NoError(err)
}

func (s *AdServiceTestSuite) TestAdService_DeleteUserIncorrectUserId() {
	request := &proto.DeleteUserRequest{Id: 10}

	s.app.On("DeleteUser", request.Id).Return(app.IncorrectUserId)

	service := NewService(&s.app)
	_, err := service.DeleteUser(context.TODO(), request)
	s.ErrorIs(err, ErrIncorrectUserId.Err())
}

func (s *AdServiceTestSuite) TestAdService_DeleteAd() {
	request := &proto.DeleteAdRequest{AdId: 1, UserId: 1}
	s.app.On("DeleteAd", request.AdId, request.UserId).Return(nil)

	service := NewService(&s.app)
	_, err := service.DeleteAd(context.TODO(), request)
	s.NoError(err)
}

func (s *AdServiceTestSuite) TestAdService_DeleteAdIncorrectUserId() {
	request := &proto.DeleteAdRequest{AdId: 1, UserId: 10}
	s.app.On("DeleteAd", request.AdId, request.UserId).Return(app.IncorrectUserId)

	service := NewService(&s.app)
	_, err := service.DeleteAd(context.TODO(), request)
	s.ErrorIs(err, ErrIncorrectUserId.Err())
}

func (s *AdServiceTestSuite) TestAdService_DeleteAdIncorrectAdId() {
	request := &proto.DeleteAdRequest{AdId: 10, UserId: 1}
	s.app.On("DeleteAd", request.AdId, request.UserId).Return(app.IncorrectAdId)

	service := NewService(&s.app)
	_, err := service.DeleteAd(context.TODO(), request)
	s.ErrorIs(err, ErrIncorrectAdId.Err())
}

func (s *AdServiceTestSuite) TestAdService_GetAd() {
	request := &proto.GetAdRequest{AdId: 10}
	expect := &ads.Ad{ID: 1, Title: "title", Text: "text 1", AuthorID: 3}

	s.app.On("GetAd", request.AdId).Return(expect, nil)

	service := NewService(&s.app)
	response, err := service.GetAd(context.TODO(), request)
	s.NoError(err)
	s.Equal(response, AdSuccessResponse(expect))
}

func (s *AdServiceTestSuite) TestAdService_GetAdIncorrectAdId() {
	request := &proto.GetAdRequest{AdId: 10}
	expect := &ads.Ad{ID: 1, Title: "title", Text: "text 1", AuthorID: 3}

	s.app.On("GetAd", request.AdId).Return(expect, app.IncorrectAdId)

	service := NewService(&s.app)
	_, err := service.GetAd(context.TODO(), request)
	s.ErrorIs(err, ErrIncorrectAdId.Err())
}
