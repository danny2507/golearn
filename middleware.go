// midlware.go
package main

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"net/http"
	"strings"
)

// AuthMiddleware checks authentication
func AuthMiddleware() gin.HandlerFunc {
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
