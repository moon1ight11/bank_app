package middleware

import (
	"bank_app/internal/api/jwt"
	"bank_app/internal/api/models"
	"bank_app/pkg/logger"
	"github.com/gin-gonic/gin"
	"net/http"
)

// аус-мидлвар для всех
func Auth(jwtService jwt.TokenService, logger logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer c.Next()

		value, err := c.Cookie("cookie")
		if err != nil {
			logger.Error("Error in Auth-middleware", "error:", err)
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			c.Abort()
			return
		}

		claims := jwt.Claims{}
		token, err := jwtService.ParseToken(value, &claims)
		if err != nil {
			logger.Error("Error in Auth-middleware", "error:", err)
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			c.Abort()
			return
		}

		if !token.Valid {
			logger.Error("Error in Auth-middleware", "error:", "token not valid")
			c.JSON(http.StatusForbidden, gin.H{"error": "Token not valid"})
			c.Abort()
			return
		}

		c.Set("UserId", *claims.UserId)
		c.Set("UserRole", claims.Role)
	}
}

// role-check для базового и выше
func AuthUser(logger logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer c.Next()

		UserRole, exist := c.Get("UserRole")
		if !exist {
			logger.Error("Error in AuthUser-middleware", "error:", "user role not found")
			c.JSON(http.StatusForbidden, gin.H{"error": "User role not found"})
			c.Abort()
			return
		}

		if UserRole != models.RoleBasic && UserRole != models.RoleVerificator && UserRole != models.RoleAdmin {
			logger.Error("Error in Auth-middleware", "error:", "access violation")
			c.JSON(http.StatusForbidden, gin.H{"error": "Access violation"})
			c.Abort()
			return
		}
	}
}

// role-check для верификатора и выше
func AuthVerificator(logger logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer c.Next()

		UserRole, exist := c.Get("UserRole")
		if !exist {
			logger.Error("Error in AuthVerificator-middleware", "error:", "user role not found")
			c.JSON(http.StatusForbidden, gin.H{"error": "User role not found"})
			c.Abort()
			return
		}

		if UserRole != models.RoleVerificator && UserRole != models.RoleAdmin {
			logger.Error("Error in AuthVerificator-middleware", "error:", "access violation")
			c.JSON(http.StatusForbidden, gin.H{"error": "Access violation"})
			c.Abort()
			return
		}
	}
}

// role-check для админа
func AuthAdmin(logger logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer c.Next()

		UserRole, exist := c.Get("UserRole")
		if !exist {
			logger.Error("Error in AuthAdmin-middleware", "error:", "user role not found")
			c.JSON(http.StatusForbidden, gin.H{"error": "User role not found"})
			c.Abort()
			return
		}

		if UserRole != models.RoleAdmin {
			logger.Error("Error in AuthVerificator-middleware", "error:", "access violation")
			c.JSON(http.StatusForbidden, gin.H{"error": "Access violation"})
			c.Abort()
			return
		}
	}
}
