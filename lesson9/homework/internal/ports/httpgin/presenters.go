package httpgin

import (
	"github.com/gin-gonic/gin"
	"homework9/internal/ads"
	"homework9/internal/users"
	"time"
)

type createAdRequest struct {
	Title  string `json:"title"`
	Text   string `json:"text"`
	UserID int64  `json:"user_id"`
}

type adResponse struct {
	ID           int64     `json:"id"`
	Title        string    `json:"title"`
	Text         string    `json:"text"`
	AuthorID     int64     `json:"author_id"`
	Published    bool      `json:"published"`
	DateUpdate   time.Time `json:"date_update"`
	DateCreating time.Time `json:"date_creating"`
}

type userResponse struct {
	ID       int64  `json:"id"`
	Nickname string `json:"nickname"`
	Email    string `json:"email"`
}

type changeAdStatusRequest struct {
	Published bool  `json:"published"`
	UserID    int64 `json:"user_id"`
}

type updateAdRequest struct {
	Title  string `json:"title"`
	Text   string `json:"text"`
	UserID int64  `json:"user_id"`
}

type createUpdateUserRequest struct {
	NickName string `json:"nickname"`
	Email    string `json:"email"`
}

func AdSuccessResponse(ad *ads.Ad) *gin.H {
	return &gin.H{
		"data": adResponse{
			ID:           ad.ID,
			Title:        ad.Title,
			Text:         ad.Text,
			AuthorID:     ad.AuthorID,
			Published:    ad.Published,
			DateUpdate:   ad.DateUpdate,
			DateCreating: ad.DateCreating,
		},
		"error": nil,
	}
}

func UserSuccessResponse(ad *users.User) *gin.H {
	return &gin.H{
		"data": userResponse{
			ID:       ad.ID,
			Nickname: ad.Nickname,
			Email:    ad.Email,
		},
		"error": nil,
	}
}

func ErrorResponse(err error) *gin.H {
	return &gin.H{
		"data":  nil,
		"error": err.Error(),
	}
}

func AdsSuccessResponse(a []ads.Ad) *gin.H {
	var response []adResponse
	for i := range a {
		response = append(response, adResponse{
			ID:           a[i].ID,
			Title:        a[i].Title,
			Text:         a[i].Text,
			AuthorID:     a[i].AuthorID,
			Published:    a[i].Published,
			DateUpdate:   a[i].DateUpdate,
			DateCreating: a[i].DateCreating,
		})
	}

	return &gin.H{
		"data":  response,
		"error": nil,
	}
}
