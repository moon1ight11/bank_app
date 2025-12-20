package middleware

import (
	"bank_app/internal/jwt"
	"bank_app/internal/storage/repos/users"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// аус-мидлвар для юзера
func AuthUser(jwtService jwt.TokenService) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer c.Next()

		value, err := c.Cookie("cookie")
		if err != nil {
			log.Println("Error in get value from cookie")
			c.JSON((http.StatusForbidden), gin.H{"error": err.Error()})
			c.Abort()
			return
		}

		claims := jwt.Claims{}

		token, err := jwtService.ParseToken(value, &claims)
		if err != nil {
			log.Println("Error in parse token", err)
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			c.Abort()
			return
		}

		if !token.Valid {
			log.Println("Token not valid")
			c.JSON(http.StatusForbidden, gin.H{"error": "Token not valid"})
			c.Abort()
			return
		}

		if claims.Roles != users.RoleUser {
			log.Println("Access violation: you are not 'User'")
			c.JSON(http.StatusForbidden, gin.H{"error": "Access violation"})
			c.Abort()
			return
		}

		c.Set("UserId", claims.UserId)
	}
}

// аус-мидлвар для верификатора
func AuthVerificator(jwtService jwt.TokenService) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer c.Next()

		value, err := c.Cookie("cookie")
		if err != nil {
			log.Println("Error in get value from cookie")
			c.JSON((http.StatusForbidden), gin.H{"error": err.Error()})
			c.Abort()
			return
		}

		claims := jwt.Claims{}

		token, err := jwtService.ParseToken(value, &claims)
		if err != nil {
			log.Println("Error in parse token", err)
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			c.Abort()
			return
		}

		if !token.Valid {
			log.Println("Token not valid")
			c.JSON(http.StatusForbidden, gin.H{"error": "Token not valid"})
			c.Abort()
			return
		}

		if claims.Roles != users.RoleVerificator {
			log.Println("Access violation: you are not 'Verificator'")
			c.JSON(http.StatusForbidden, gin.H{"error": "Access violation"})
			c.Abort()
			return
		}

		c.Set("UserId", claims.UserId)
	}
}

// аус-мидлвар для админа
func AuthAdmin(jwtService jwt.TokenService) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer c.Next()

		value, err := c.Cookie("cookie")
		if err != nil {
			log.Println("Error in get value from cookie")
			c.JSON((http.StatusForbidden), gin.H{"error": err.Error()})
			c.Abort()
			return
		}

		claims := jwt.Claims{}

		token, err := jwtService.ParseToken(value, &claims)
		if err != nil {
			log.Println("Error in parse token", err)
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			c.Abort()
			return
		}

		if !token.Valid {
			log.Println("Token not valid")
			c.JSON(http.StatusForbidden, gin.H{"error": "Token not valid"})
			c.Abort()
			return
		}

		if claims.Roles != users.RoleAdmin {
			log.Println("Access violation: you are not 'Admin'")
			c.JSON(http.StatusForbidden, gin.H{"error": "Access violation"})
			c.Abort()
			return
		}

		c.Set("UserId", claims.UserId)
	}
}
