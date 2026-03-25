package middleware

import (
	"bank_app/internal/api/jwt"
	"bank_app/internal/api/models"
	"log"
	"net/http"
	"github.com/gin-gonic/gin"
)

// аус-мидлвар для всех
func Auth(jwtService jwt.TokenService) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer c.Next()

		value, err := c.Cookie("cookie")
		if err != nil {
			log.Println("Error in get value from cookie")
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
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

		c.Set("UserId", *claims.UserId)
		c.Set("UserRole", claims.Role)
	}
}

// role-check для базового и выше
func AuthUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer c.Next()

		UserRole, exist := c.Get("UserRole")
		if !exist {
			log.Println("Error in get User role")
			c.JSON(http.StatusForbidden, gin.H{"error": "User role not found"})
			c.Abort()
			return
		}

		if UserRole != models.RoleBasic && UserRole != models.RoleVerificator && UserRole != models.RoleAdmin{
			log.Println("Access violation: you are not 'User'")
			c.JSON(http.StatusForbidden, gin.H{"error": "Access violation"})
			c.Abort()
			return
		}
	}
}

// role-check для верификатора и выше
func AuthVerificator() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer c.Next()

		UserRole, exist := c.Get("UserRole")
		if !exist {
			log.Println("Error in get User role")
			c.JSON(http.StatusForbidden, gin.H{"error": "User role not found"})
			c.Abort()
			return
		}

		if UserRole != models.RoleVerificator && UserRole != models.RoleAdmin {
			log.Println("Access violation: you are not 'Verificator' or 'Admin")
			c.JSON(http.StatusForbidden, gin.H{"error": "Access violation"})
			c.Abort()
			return
		}
	}
}

// role-check для админа
func AuthAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer c.Next()

		UserRole, exist := c.Get("UserRole")
		if !exist {
			log.Println("Error in get User role")
			c.JSON(http.StatusForbidden, gin.H{"error": "User role not found"})
			c.Abort()
			return
		}

		if UserRole != models.RoleAdmin {
			log.Println("Access violation: you are not 'Admin'")
			c.JSON(http.StatusForbidden, gin.H{"error": "Access violation"})
			c.Abort()
			return
		}
	}
}
