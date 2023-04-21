package grpc

import (
	"google.golang.org/protobuf/types/known/timestamppb"
	"homework9/internal/ads"
	"homework9/internal/users"
)

func AdSuccessResponse(ad *ads.Ad) *AdResponse {
	return &AdResponse{
		Id:           ad.ID,
		Title:        ad.Title,
		Text:         ad.Text,
		UserId:       ad.AuthorID,
		Published:    ad.Published,
		DateUpdate:   timestamppb.New(ad.DateUpdate),
		DateCreating: timestamppb.New(ad.DateCreating),
	}
}

func UserSuccessResponse(user *users.User) *UserResponse {
	return &UserResponse{
		Id:       user.ID,
		Nickname: user.Nickname,
		Email:    user.Email,
	}
}

func AdsSuccessResponse(list []ads.Ad) *ListAdResponse {
	var response ListAdResponse

	for _, v := range list {
		response.List = append(response.List, AdSuccessResponse(&v))
	}

	return &response
}
