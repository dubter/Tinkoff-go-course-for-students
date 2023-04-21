package httpgin

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"homework9/internal/app"
)

func NewHTTPServer(port string, a app.App) *http.Server {
	gin.SetMode(gin.ReleaseMode)
	handler := gin.New()
	s := &http.Server{Addr: port, Handler: handler}
	api := handler.Group("/api/v1")
	AppRouter(api, a)
	return s
}
