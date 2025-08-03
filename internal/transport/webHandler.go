package transport

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type WebHandler struct{}

func (w *WebHandler) ServePage(file string) func(c *gin.Context) {
	return func(c *gin.Context) {
		c.HTML(http.StatusOK, file, nil)
	}
}

func (w *WebHandler) RegisterRoutes(router *gin.Engine) {
	router.GET("/login", w.ServePage("login.html"))
	router.GET("/register", w.ServePage("register.html"))
}
