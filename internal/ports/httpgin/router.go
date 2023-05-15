package httpgin

import (
	"github.com/gin-gonic/gin"
	"homework9/internal/app"
)

func AppRouter(r *gin.RouterGroup, a app.App) {
	r.POST("/ads", createAd(a))                    // Метод для создания объявления (ad)
	r.PUT("/ads/:ad_id/status", changeAdStatus(a)) // Метод для изменения статуса объявления (опубликовано - Published = true или снято с публикации Published = false)
	r.PUT("/ads/:ad_id", updateAd(a))              // Метод для обновления текста(Text) или заголовка(Title) объявления
	r.GET("/ads/:ad_id", getAd(a))                 // Метод для доступа к объявления по ID
	r.GET("/ads/title/:title", getAdByTitle(a))    // Метод для доступа к объявлению по Title
	r.GET("/ads", getAds(a))
	r.GET("/ads/params", getAdsFilter(a))
	r.DELETE("/ads/:ad_id", deleteAd(a))

	r.POST("/users", createUser(a))            // Метод для создания пользователя (user)
	r.DELETE("/users/:user_id", deleteUser(a)) // Метод для удаления пользователя (user)
	r.GET("/users/:user_id", getUser(a))       // Метод для доступа к пользователю по ID
}
