// midlware.go
package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"net/http"
	"os"
	"strings"
)

type JWTClaims struct {
	UserID int `json:"user_id"`
	jwt.StandardClaims
}

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

// AuthMiddleware checks authentication
func (ctrl *Controller) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// get the token from the Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "missing/invalid authorization header"})
			c.Abort()
			return
		}
		// trim the unnecessary Bearer prefix
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		// validtae token
		token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})
		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			c.Abort()
			return
		}

		_, err = ctrl.PostgreService.GetActiveToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token on db"})
			c.Abort()
			return
		}
		// get userid from token claim
		claims, ok := token.Claims.(*JWTClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token claims"})
			c.Abort()
			return
		}

		// set userID in context
		c.Set("userID", claims.UserID)
		c.Next()
	}
}
