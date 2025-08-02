package transport

import (
	"fmt"
	"net/http"

	"github.com/Vin-Xi/auth/internal/service"
	"github.com/Vin-Xi/auth/internal/user"
	"github.com/Vin-Xi/auth/internal/util"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	UserService service.Service
	JwtEngine   *util.JWTEngine
}

func (h *UserHandler) RegisterRoutes(router *gin.Engine) {
	api := router.Group("/api/v1")
	{
		api.POST("/register", h.register)
		api.POST("/login", h.login)

		protected := api.Group("/").Use(h.authMiddleware())
		{
			protected.GET("/profile", h.getProfile)
		}
	}
}

type registerReq struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

func (h *UserHandler) register(c *gin.Context) {
	var req registerReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err := h.UserService.Register(c.Request.Context(), req.Email, req.Password)

	if err != nil {
		fmt.Println(req.Email, req.Password, err)
		if err == service.ErrUserAlreadyExists {
			c.JSON(http.StatusConflict, gin.H{"error": "user already exists"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to register user"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "user successfully registered"})
}

type loginReq struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func (h *UserHandler) login(c *gin.Context) {
	var req loginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	u, err := h.UserService.Login(c.Request.Context(), req.Email, req.Password)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	token, err := h.JwtEngine.Generate(u.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})

}

func (h *UserHandler) getProfile(c *gin.Context) {
	u, exists := c.Get("user")

	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	currentUser := u.(*user.User)

	c.JSON(http.StatusOK, gin.H{
		"id":    currentUser.ID,
		"email": currentUser.Email,
	})
}
