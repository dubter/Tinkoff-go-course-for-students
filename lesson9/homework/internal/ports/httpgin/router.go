package httpgin

import (
	"github.com/gin-gonic/gin"
	"homework9/internal/app"
)

func AppRouter(r *gin.RouterGroup, a app.App) {
	r.Use(Logger())             // Middleware для логгирования всех запросов
	r.Use(RecoveryWithLogger()) // Middleware для обработки panic с логгированием

	adsR := r.Group("/ads")
	adsR.POST("", createAd(a))                    // Метод для создания объявления (ad)
	adsR.PUT("/:ad_id/status", changeAdStatus(a)) // Метод для изменения статуса объявления (опубликовано - Published = true или снято с публикации Published = false)
	adsR.PUT("/:ad_id", updateAd(a))              // Метод для обновления текста(Text) или заголовка(Title) объявления
	adsR.DELETE("/:ad_id", deleteAd(a))           // Метод для удаления объявления по id

	adsR.GET("/:ad_id", getAd(a))                       // Метод для вывода объявления по id
	adsR.GET("", getListAds(a))                         // Метод для вывода списка опубликаванных объявлений
	adsR.GET("/with_filter", getListAds(a))             // Метод для вывода списка опубликаванных объявлений
	adsR.GET("/search/:ad_title", getListAdsByTitle(a)) // Метод для поиска объявлений по названию

	userR := r.Group("/users")
	userR.POST("", createUser(a))            // Метод для создания пользователя (user)
	userR.PUT("/:user_id", updateUser(a))    // Метод для редактирования данных пользователя
	userR.GET("/:user_id", getUser(a))       // Метод для вывода пользователя по id
	userR.DELETE("/:user_id", deleteUser(a)) // Метод для удаления пользователя id
}
