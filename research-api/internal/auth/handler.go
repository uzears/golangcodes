package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/uzears/golangcodes/research-api/internal/platform/logger"
)

type Handler struct {
	service Service
	log     logger.Logger
}

func NewHandler(service Service, log logger.Logger) *Handler {
	return &Handler{
		service: service,
		log:     log,
	}
}

func (h *Handler) Register(c *gin.Context) {
	log := c.MustGet("logger").(logger.Logger)
	log.Debug("register request received")

	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Warn("Invalid registration request", "err", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})
		return
	}

	userID, err := h.service.Register(
		c.Request.Context(),
		req.Email,
		req.Password,
	)
	if err != nil {
		log.Error("Registration failed", "err", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	log.Info("user registered successfully", "user_id", userID)

	c.JSON(http.StatusCreated, gin.H{
		"id": userID,
	})

}

func (h *Handler) Login(c *gin.Context) {
	log := c.MustGet("logger").(logger.Logger)
	log.Debug("login request received")

	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Warn("Invalid login request", "err", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})
		return
	}
	token, err := h.service.Login(
		c.Request.Context(),
		req.Email,
		req.Password,
	)

	if err != nil {
		log.Error("Login failed", "err", err)
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})
		return
	}

	log.Info("user logged in successfully", "email", req.Email)

	c.JSON(http.StatusOK, LoginResponse{
		Token: token,
	})
}

func (h *Handler) Me(c *gin.Context) {
	log := c.MustGet("logger").(logger.Logger)
	log.Debug("me request received")

	userID, exists := c.Get("user_id")
	if !exists {
		log.Debug("me request missing user_id in context")
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "unauthorized",
		})
		return
	}

	user, err := h.service.GetByID(
		c.Request.Context(),
		userID.(string),
	)
	if err != nil {
		log.Error("failed to fetch user", "user_id", userID, "err", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to fetch user",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":    user.ID,
		"email": user.Email,
	})
}
