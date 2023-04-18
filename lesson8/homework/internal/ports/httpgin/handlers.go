package httpgin

import (
	"errors"
	"github.com/gin-gonic/gin"
	"homework8/internal/app"
	"net/http"
	"strconv"
)

// Метод для создания объявления (ad)
func createAd(a app.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		var reqBody createAdRequest
		err := c.Bind(&reqBody)
		if err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse(err))
			return
		}

		ad, ok := a.CreateAd(reqBody.Title, reqBody.Text, reqBody.UserID)

		if errors.Is(ok, app.ValidateError) {
			c.JSON(http.StatusBadRequest, ErrorResponse(ok))
			return
		}

		if errors.Is(ok, app.IncorrectUserId) {
			c.JSON(http.StatusForbidden, ErrorResponse(ok))
			return
		}

		if err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResponse(err))
			return
		}
		c.JSON(http.StatusOK, AdSuccessResponse(ad))
	}
}

// Метод для создания пользователя (user)
func createUser(a app.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		var reqBody createUpdateUserRequest
		err := c.Bind(&reqBody)
		if err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse(err))
			return
		}

		user, ok := a.CreateUser(reqBody.NickName, reqBody.Email)

		if errors.Is(ok, app.ValidateError) {
			c.JSON(http.StatusBadRequest, ErrorResponse(ok))
			return
		}

		if err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResponse(err))
			return
		}
		c.JSON(http.StatusOK, UserSuccessResponse(user))
	}
}

// Метод для изменения статуса объявления (опубликовано - Published = true или снято с публикации Published = false)
func changeAdStatus(a app.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		var reqBody changeAdStatusRequest
		err := c.Bind(&reqBody)
		if err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse(err))
			return
		}

		adID := c.Param("ad_id")
		num, errToInt := strconv.Atoi(adID)
		if errToInt != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse(errToInt))
			return
		}

		ad, ok := a.ChangeAdStatus(int64(num), reqBody.UserID, reqBody.Published)
		if errors.Is(ok, app.IncorrectUserId) {
			c.JSON(http.StatusForbidden, ErrorResponse(ok))
			return
		}

		if err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResponse(err))
			return
		}

		c.JSON(http.StatusOK, AdSuccessResponse(ad))
	}
}

// Метод для обновления текста(Text) или заголовка(Title) объявления
func updateAd(a app.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		var reqBody updateAdRequest
		err := c.Bind(&reqBody)
		if err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse(err))
			return
		}

		adID := c.Param("ad_id")
		num, errToInt := strconv.Atoi(adID)
		if errToInt != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse(errToInt))
			return
		}

		ad, ok := a.UpdateAd(int64(num), reqBody.UserID, reqBody.Title, reqBody.Text)
		if errors.Is(ok, app.IncorrectUserId) {
			c.JSON(http.StatusForbidden, ErrorResponse(ok))
			return
		}

		if errors.Is(ok, app.ValidateError) {
			c.JSON(http.StatusBadRequest, ErrorResponse(ok))
			return
		}

		if err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResponse(err))
			return
		}

		c.JSON(http.StatusOK, AdSuccessResponse(ad))
	}
}

// Метод для редактирования данных пользователя
func updateUser(a app.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		var reqBody createUpdateUserRequest
		err := c.Bind(&reqBody)
		if err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse(err))
			return
		}

		adID := c.Param("user_id")
		num, errToInt := strconv.Atoi(adID)
		if errToInt != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse(errToInt))
			return
		}

		user, ok := a.UpdateUser(int64(num), reqBody.NickName, reqBody.Email)
		if errors.Is(ok, app.IncorrectUserId) {
			c.JSON(http.StatusForbidden, ErrorResponse(ok))
			return
		}

		if errors.Is(ok, app.ValidateError) {
			c.JSON(http.StatusBadRequest, ErrorResponse(ok))
			return
		}

		if err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResponse(err))
			return
		}

		c.JSON(http.StatusOK, UserSuccessResponse(user))
	}
}

// Метод для получения списка выложенных объявлений
func getListAds(a app.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		filters := make(map[string]any)
		published := c.Query("published")
		if published != "" {
			if published == "true" {
				filters["published"] = true
			} else {
				filters["published"] = false
			}
		}

		userId := c.Query("user_id")
		if userId != "" {
			id, errToInt := strconv.Atoi(userId)
			if errToInt != nil {
				c.JSON(http.StatusBadRequest, ErrorResponse(errToInt))
				return
			}
			filters["user_id"] = int64(id)
		}

		dateCreating := c.Query("date_creating")
		if dateCreating != "" {
			filters["date_creating"] = dateCreating
		}

		ads := a.GetListAds(filters)

		c.JSON(http.StatusOK, AdsSuccessResponse(ads))
	}
}

// Метод для вывода объявления по id
func getAd(a app.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		adID := c.Param("ad_id")
		num, errToInt := strconv.Atoi(adID)
		if errToInt != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse(errToInt))
			return
		}

		ad, ok := a.GetAd(int64(num))
		if errors.Is(ok, app.IncorrectAdId) {
			c.JSON(http.StatusForbidden, ErrorResponse(ok))
			return
		}

		c.JSON(http.StatusOK, AdSuccessResponse(ad))
	}
}

// Метод для поиска объявлений по названию
func getListAdsByTitle(a app.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		title := c.Param("ad_title")
		ads := a.GetListAdsByTitle(title)

		c.JSON(http.StatusOK, AdsSuccessResponse(ads))
	}
}
