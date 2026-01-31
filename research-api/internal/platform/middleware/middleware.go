package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/uzears/golangcodes/research-api/internal/platform/logger"
)

func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		reqID := c.GetHeader("X-Request-Id")
		if reqID == "" {
			reqID = generateRequestID()
		}

		c.Set("request_id", reqID)
		c.Writer.Header().Set("X-Request-Id", reqID)

		c.Next()
	}
}
func Logger(log logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		reqID, _ := c.Get("request_id")

		reqLog := log.With(
			"request_id", reqID,
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
		)

		c.Set("logger", reqLog)

		c.Next()

		status := c.Writer.Status()

		reqLog.Info(
			"http request completed",
			"status", status,
		)
	}
}
func JWT(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		if auth == "" || !strings.HasPrefix(auth, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "missing or invalid token",
			})
			return
		}

		tokenStr := strings.TrimPrefix(auth, "Bearer ")

		claims, err := validateToken(tokenStr, secret)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "invalid token",
			})
			return
		}

		// Attach user info to context
		c.Set("user_id", claims.UserID)

		c.Next()
	}
}
func generateRequestID() string {
	// UUID or nanoid (implementation choice)
	return "req-123456"
}

type TokenClaims struct {
	UserID string
}

func validateToken(tokenStr, secret string) (*TokenClaims, error) {
	return &TokenClaims{
		UserID: "user-123",
	}, nil
}
