package httpgin

import (
	"github.com/gin-gonic/gin"
	"log"
	"os"

	"homework8/internal/app"
)

func AppRouter(r *gin.RouterGroup, a app.App) {
	r.Use(Logger())                                        // Middleware для логгирования всех запросов
	r.Use(RecoveryWithLogger(log.New(os.Stdout, "\n", 0))) // Middleware для обработки panic с логгированием

	adsR := r.Group("/ads")
	adsR.POST("", createAd(a))                    // Метод для создания объявления (ad)
	adsR.PUT("/:ad_id/status", changeAdStatus(a)) // Метод для изменения статуса объявления (опубликовано - Published = true или снято с публикации Published = false)
	adsR.PUT("/:ad_id", updateAd(a))              // Метод для обновления текста(Text) или заголовка(Title) объявления

	adsR.GET("/:ad_id", getAd(a))                       // Метод для вывода объявления по id
	adsR.GET("", getListAds(a))                         // Метод для вывода списка опубликаванных объявлений
	adsR.GET("/with_filter", getListAds(a))             // Метод для вывода списка опубликаванных объявлений
	adsR.GET("/search/:ad_title", getListAdsByTitle(a)) // Метод для поиска объявлений по названию

	userR := r.Group("/users")
	userR.POST("", createUser(a))         // Метод для создания пользователя (user)
	userR.PUT("/:user_id", updateUser(a)) // Метод для редактирования данных пользователя
}
