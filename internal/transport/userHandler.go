package transport

import (
	"fmt"
	"log/slog"
	"net/http"
	"net/url"

	"github.com/Vin-Xi/auth/internal/service"
	"github.com/Vin-Xi/auth/internal/util"
	"github.com/Vin-Xi/auth/pkg/logger"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	UserService service.Service
	JwtEngine   *util.JWTEngine
}

func (h *UserHandler) RegisterRoutes(router *gin.Engine) {
	router.POST("/register", h.register)
	router.POST("/login", h.login)
}

func (h *UserHandler) register(c *gin.Context) {
	email := c.PostForm("email")
	password := c.PostForm("password")
	firstName := c.PostForm("fName")
	lastName := c.PostForm("lName")

	_, err := h.UserService.Register(c.Request.Context(), email, password, firstName, lastName)
	fmt.Println(email, password, firstName, lastName)

	if err != nil {
		fmt.Println(email, password, err)
		if err == service.ErrUserAlreadyExists {
			slog.Info("user already exists")
			c.Redirect(http.StatusFound, "/register?error=Already+Exists")
			return
		}
		logger.Log.ErrorWithStack("failed to create user", err)

		c.Redirect(http.StatusFound, "/register?error=Registration+Failed")
		return
	}
	fmt.Println(email, password)

	c.Redirect(http.StatusFound, "/log")
}

func (h *UserHandler) login(c *gin.Context) {
	redirectURI := c.Query("redirect_uri")

	email := c.PostForm("email")
	password := c.PostForm("password")

	u, err := h.UserService.Login(c.Request.Context(), email, password)

	if err != nil {
		logger.Log.Info("invalid credentials", "email", email)
		c.Redirect(http.StatusFound, "/login?error=Invalid+credentials&redirect_uri="+url.QueryEscape(redirectURI))
		return
	}

	token, err := h.JwtEngine.Generate(u.ID)
	if err != nil {
		logger.Log.ErrorWithStack("login failed", err)
		c.Redirect(http.StatusFound, "/login?error=Login+Failed&redirect_uri="+url.QueryEscape(redirectURI))
		return
	}

	c.SetCookie("auth_token", token, 3600, "/", "localhost", true, true)
	c.Redirect(http.StatusFound, redirectURI)
}
