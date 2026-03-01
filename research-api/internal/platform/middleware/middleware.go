package middleware

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
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
		reqLog.Debug("http request started")

		c.Next()

		status := c.Writer.Status()

		reqLog.Info(
			"http request completed",
			"status", status,
		)
	}
}
func CORS(allowedOrigins []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")
		for _, allowed := range allowedOrigins {
			if origin == allowed {
				c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
				break
			}
		}
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET,POST,PUT,PATCH,DELETE,OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Authorization,Content-Type")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
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
	b := make([]byte, 8)
	if _, err := rand.Read(b); err != nil {
		return fmt.Sprintf("req-%d", time.Now().UnixNano())
	}
	return "req-" + hex.EncodeToString(b)
}

type TokenClaims struct {
	UserID string `json:"sub"`
	jwt.RegisteredClaims
}

func validateToken(tokenStr, secret string) (*TokenClaims, error) {
	claims := &TokenClaims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}
	if claims.UserID == "" {
		return nil, fmt.Errorf("missing subject")
	}
	return claims, nil
}
