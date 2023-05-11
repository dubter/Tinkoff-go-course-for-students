package grpc

import (
	"context"
	"errors"
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"homework10/internal/app"
	"homework10/internal/ports/grpc/proto"
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

func (service *AdService) CreateAd(_ context.Context, req *proto.CreateAdRequest) (*proto.AdResponse, error) {
	ad, ok := service.a.CreateAd(req.GetTitle(), req.GetText(), req.GetUserId())

	if errors.Is(ok, app.ValidateError) {
		return nil, ErrValidate.Err()
	}

	if errors.Is(ok, app.IncorrectUserId) {
		return nil, ErrIncorrectUserId.Err()
	}

	return AdSuccessResponse(ad), OkStatus.Err()
}

func (service *AdService) ChangeAdStatus(_ context.Context, req *proto.ChangeAdStatusRequest) (*proto.AdResponse, error) {
	ad, ok := service.a.ChangeAdStatus(req.GetAdId(), req.GetUserId(), req.GetPublished())

	if errors.Is(ok, app.IncorrectUserId) {
		return nil, ErrIncorrectUserId.Err()
	}

	return AdSuccessResponse(ad), OkStatus.Err()
}

func (service *AdService) UpdateAd(_ context.Context, req *proto.UpdateAdRequest) (*proto.AdResponse, error) {
	ad, ok := service.a.UpdateAd(req.GetAdId(), req.GetUserId(), req.GetTitle(), req.GetText())

	if errors.Is(ok, app.ValidateError) {
		return nil, ErrValidate.Err()
	}

	if errors.Is(ok, app.IncorrectUserId) {
		return nil, ErrIncorrectUserId.Err()
	}

	return AdSuccessResponse(ad), OkStatus.Err()
}

func (service *AdService) ListAdsWithFilter(_ context.Context, req *proto.GetListAdsWithFilterRequest) (*proto.ListAdResponse, error) {
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

func (service *AdService) ListAdsByTitle(_ context.Context, req *proto.GetListAdsByTitleRequest) (*proto.ListAdResponse, error) {
	list := service.a.GetListAdsByTitle(req.GetTitle())
	return AdsSuccessResponse(list), OkStatus.Err()
}

func (service *AdService) CreateUser(_ context.Context, req *proto.CreateUserRequest) (*proto.UserResponse, error) {
	user, ok := service.a.CreateUser(req.GetNickname(), req.GetEmail())

	if errors.Is(ok, app.ValidateError) {
		return nil, ErrValidate.Err()
	}

	return UserSuccessResponse(user), OkStatus.Err()
}

func (service *AdService) UpdateUser(_ context.Context, req *proto.UpdateUserRequest) (*proto.UserResponse, error) {
	user, ok := service.a.UpdateUser(req.GetUserId(), req.GetNickname(), req.GetEmail())

	if errors.Is(ok, app.IncorrectUserId) {
		return nil, ErrIncorrectUserId.Err()
	}

	if errors.Is(ok, app.ValidateError) {
		return nil, ErrValidate.Err()
	}

	return UserSuccessResponse(user), OkStatus.Err()
}

func (service *AdService) GetUser(_ context.Context, req *proto.GetUserRequest) (*proto.UserResponse, error) {
	user, ok := service.a.GetUser(req.GetId())

	if errors.Is(ok, app.IncorrectUserId) {
		return nil, ErrIncorrectUserId.Err()
	}

	return UserSuccessResponse(user), OkStatus.Err()
}

func (service *AdService) DeleteUser(_ context.Context, req *proto.DeleteUserRequest) (*emptypb.Empty, error) {
	ok := service.a.DeleteUser(req.GetId())

	if errors.Is(ok, app.IncorrectUserId) {
		return nil, ErrIncorrectUserId.Err()
	}

	return new(emptypb.Empty), OkStatus.Err()
}

func (service *AdService) DeleteAd(_ context.Context, req *proto.DeleteAdRequest) (*emptypb.Empty, error) {
	ok := service.a.DeleteAd(req.GetAdId(), req.GetUserId())

	if errors.Is(ok, app.IncorrectUserId) {
		return nil, ErrIncorrectUserId.Err()
	}

	if errors.Is(ok, app.IncorrectAdId) {
		return nil, ErrIncorrectAdId.Err()
	}

	return new(emptypb.Empty), OkStatus.Err()
}

func (service *AdService) GetAd(_ context.Context, req *proto.GetAdRequest) (*proto.AdResponse, error) {
	ad, ok := service.a.GetAd(req.GetAdId())

	if errors.Is(ok, app.IncorrectAdId) {
		return nil, ErrIncorrectAdId.Err()
	}

	return AdSuccessResponse(ad), OkStatus.Err()
}
