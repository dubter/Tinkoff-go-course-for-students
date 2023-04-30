package grpc

import (
	"google.golang.org/protobuf/types/known/timestamppb"
	"homework10/internal/ads"
	"homework10/internal/ports/grpc/proto"
	"homework10/internal/users"
)

func AdSuccessResponse(ad *ads.Ad) *proto.AdResponse {
	return &proto.AdResponse{
		Id:           ad.ID,
		Title:        ad.Title,
		Text:         ad.Text,
		UserId:       ad.AuthorID,
		Published:    ad.Published,
		DateUpdate:   timestamppb.New(ad.DateUpdate),
		DateCreating: timestamppb.New(ad.DateCreating),
	}
}

func UserSuccessResponse(user *users.User) *proto.UserResponse {
	return &proto.UserResponse{
		Id:       user.ID,
		Nickname: user.Nickname,
		Email:    user.Email,
	}
}

func AdsSuccessResponse(list []ads.Ad) *proto.ListAdResponse {
	var response proto.ListAdResponse

	for _, v := range list {
		response.List = append(response.List, AdSuccessResponse(&v))
	}

	return &response
}
