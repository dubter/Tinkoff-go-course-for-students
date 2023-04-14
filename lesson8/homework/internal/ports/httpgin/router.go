package httpgin

import (
	"github.com/gin-gonic/gin"

	"homework8/internal/app"
)

func AppRouter(r *gin.RouterGroup, a app.App) {
	r.POST("/ads", createAd(a))                    // Метод для создания объявления (ad)
	r.PUT("/ads/:ad_id/status", changeAdStatus(a)) // Метод для изменения статуса объявления (опубликовано - Published = true или снято с публикации Published = false)
	r.PUT("/ads/:ad_id", updateAd(a))              // Метод для обновления текста(Text) или заголовка(Title) объявления

	r.GET("/ads/:ad_id", getAd(a))                      // Метод для вывода объявления по id
	r.GET("/ads", getListAds(a))                        // Метод для вывода списка опубликаванных объявлений
	r.POST("/ads/with_filter", getListAds(a))           // Метод для вывода списка опубликаванных объявлений
	r.POST("/ads/search_by_name", getListAdsByTitle(a)) // Метод для поиска объявлений по названию

	r.POST("/users", createUser(a))         // Метод для создания пользователя (user)
	r.PUT("/users/:user_id", updateUser(a)) // Метод для редактирования данных пользователя
}
