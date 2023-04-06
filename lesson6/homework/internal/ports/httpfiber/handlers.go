package httpfiber

import (
	"fmt"
	"github.com/dubter/Validator"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"homework6/internal/app"
)

// Метод для создания объявления (ad)
func createAd(a app.App) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var reqBody createAdRequest
		err := c.BodyParser(&reqBody)
		if err != nil {
			c.Status(http.StatusBadRequest)
			return c.JSON(AdErrorResponse(err))
		}

		ad, errValid := a.CreateAd(reqBody.Title, reqBody.Text, reqBody.UserID)
		if errValid == fmt.Errorf("error validation") {
			c.Status(http.StatusBadRequest)
			return c.JSON(AdErrorResponse(errValid))
		}

		if Validator.Validate(*ad) != nil {
			c.Status(http.StatusBadRequest)
			return c.JSON(AdErrorResponse(Validator.Validate(*ad)))
		}

		if err != nil {
			c.Status(http.StatusInternalServerError)
			return c.JSON(AdErrorResponse(err))
		}
		return c.JSON(AdSuccessResponse(ad))
	}
}

// Метод для изменения статуса объявления (опубликовано - Published = true или снято с публикации Published = false)
func changeAdStatus(a app.App) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var reqBody changeAdStatusRequest
		if err := c.BodyParser(&reqBody); err != nil {
			c.Status(http.StatusBadRequest)
			return c.JSON(AdErrorResponse(err))
		}

		adID, err := c.ParamsInt("ad_id")
		if err != nil {
			c.Status(http.StatusBadRequest)
			return c.JSON(AdErrorResponse(err))
		}

		ad, errorId := a.ChangeAdStatus(int64(adID), reqBody.UserID, reqBody.Published)
		if errorId != nil {
			c.Status(http.StatusForbidden)
			return c.JSON(AdErrorResponse(errorId))
		}

		if Validator.Validate(*ad) != nil {
			c.Status(http.StatusBadRequest)
			return c.JSON(AdErrorResponse(Validator.Validate(*ad)))
		}

		if err != nil {
			c.Status(http.StatusInternalServerError)
			return c.JSON(AdErrorResponse(err))
		}

		return c.JSON(AdSuccessResponse(ad))
	}
}

// Метод для обновления текста(Text) или заголовка(Title) объявления
func updateAd(a app.App) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var reqBody updateAdRequest
		if err := c.BodyParser(&reqBody); err != nil {
			c.Status(http.StatusBadRequest)
			return c.JSON(AdErrorResponse(err))
		}

		adID, err := c.ParamsInt("ad_id")
		if err != nil {
			c.Status(http.StatusBadRequest)
			return c.JSON(AdErrorResponse(err))
		}

		ad, errorId := a.UpdateAd(int64(adID), reqBody.UserID, reqBody.Title, reqBody.Text)
		if errorId != nil {
			c.Status(http.StatusForbidden)
			return c.JSON(AdErrorResponse(errorId))
		}

		if Validator.Validate(*ad) != nil {
			c.Status(http.StatusBadRequest)
			return c.JSON(AdErrorResponse(Validator.Validate(*ad)))
		}

		if err != nil {
			c.Status(http.StatusInternalServerError)
			return c.JSON(AdErrorResponse(err))
		}

		return c.JSON(AdSuccessResponse(ad))
	}
}
