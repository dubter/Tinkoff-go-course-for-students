package grpc

import (
	"context"
	"errors"
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"homework9/internal/app"
)

var ErrValidate = status.New(codes.InvalidArgument, "validation error")
var ErrIncorrectUserId = status.New(codes.PermissionDenied, "incorrect user id")
var ErrIncorrectAdId = status.New(codes.NotFound, "id is not found")
var OkStatus = status.New(codes.OK, "success")

type AdService struct {
	a app.App
}

func NewService(a app.App) *AdService {
	return &AdService{a}
}

func (service *AdService) CreateAd(_ context.Context, req *CreateAdRequest) (*AdResponse, error) {
	ad, ok := service.a.CreateAd(req.Title, req.Text, req.UserId)

	if errors.Is(ok, app.ValidateError) {
		return nil, ErrValidate.Err()
	}

	if errors.Is(ok, app.IncorrectUserId) {
		return nil, ErrIncorrectUserId.Err()
	}

	return AdSuccessResponse(ad), OkStatus.Err()
}

func (service *AdService) ChangeAdStatus(_ context.Context, req *ChangeAdStatusRequest) (*AdResponse, error) {
	ad, ok := service.a.ChangeAdStatus(req.AdId, req.UserId, req.Published)

	if errors.Is(ok, app.IncorrectUserId) {
		return nil, ErrIncorrectUserId.Err()
	}

	return AdSuccessResponse(ad), OkStatus.Err()
}

func (service *AdService) UpdateAd(_ context.Context, req *UpdateAdRequest) (*AdResponse, error) {
	ad, ok := service.a.UpdateAd(req.AdId, req.UserId, req.Title, req.Text)

	if errors.Is(ok, app.ValidateError) {
		return nil, ErrValidate.Err()
	}

	if errors.Is(ok, app.IncorrectUserId) {
		return nil, ErrIncorrectUserId.Err()
	}

	return AdSuccessResponse(ad), OkStatus.Err()
}

func (service *AdService) ListAdsWithFilter(_ context.Context, req *GetListAdsWithFilterRequest) (*ListAdResponse, error) {
	filters := make(map[string]any)

	if req.UserId != nil {
		filters["user_id"] = *req.UserId
	}

	if req.DateCreating != nil {
		filters["date_creating"] = fmt.Sprint((*req.DateCreating).AsTime().UTC())
	}

	if req.Published != nil {
		filters["published"] = *req.Published
	}

	list := service.a.GetListAds(filters)
	return AdsSuccessResponse(list), OkStatus.Err()
}

func (service *AdService) ListAdsByTitle(_ context.Context, req *GetListAdsByTitleRequest) (*ListAdResponse, error) {
	list := service.a.GetListAdsByTitle(req.Title)
	return AdsSuccessResponse(list), OkStatus.Err()
}

func (service *AdService) CreateUser(_ context.Context, req *CreateUserRequest) (*UserResponse, error) {
	user, ok := service.a.CreateUser(req.Nickname, req.Email)

	if errors.Is(ok, app.ValidateError) {
		return nil, ErrValidate.Err()
	}

	return UserSuccessResponse(user), OkStatus.Err()
}

func (service *AdService) UpdateUser(_ context.Context, req *UpdateUserRequest) (*UserResponse, error) {
	user, ok := service.a.UpdateUser(req.UserId, req.Nickname, req.Email)

	if errors.Is(ok, app.IncorrectUserId) {
		return nil, ErrValidate.Err()
	}

	if errors.Is(ok, app.ValidateError) {
		return nil, ErrValidate.Err()
	}

	return UserSuccessResponse(user), OkStatus.Err()
}

func (service *AdService) GetUser(_ context.Context, req *GetUserRequest) (*UserResponse, error) {
	user, ok := service.a.GetUser(req.Id)

	if errors.Is(ok, app.IncorrectUserId) {
		return nil, ErrIncorrectUserId.Err()
	}

	return UserSuccessResponse(user), OkStatus.Err()
}

func (service *AdService) DeleteUser(_ context.Context, req *DeleteUserRequest) (*emptypb.Empty, error) {
	ok := service.a.DeleteUser(req.Id)

	if errors.Is(ok, app.IncorrectUserId) {
		return nil, ErrIncorrectUserId.Err()
	}

	return new(emptypb.Empty), OkStatus.Err()
}

func (service *AdService) DeleteAd(_ context.Context, req *DeleteAdRequest) (*emptypb.Empty, error) {
	ok := service.a.DeleteAd(req.AdId, req.UserId)

	if errors.Is(ok, app.IncorrectUserId) {
		return nil, ErrIncorrectUserId.Err()
	}

	if errors.Is(ok, app.IncorrectAdId) {
		return nil, ErrIncorrectAdId.Err()
	}

	return new(emptypb.Empty), OkStatus.Err()
}

func (service *AdService) GetAd(_ context.Context, req *GetAdRequest) (*AdResponse, error) {
	ad, ok := service.a.GetAd(req.AdId)

	if errors.Is(ok, app.IncorrectAdId) {
		return nil, ErrIncorrectAdId.Err()
	}

	return AdSuccessResponse(ad), OkStatus.Err()
}
