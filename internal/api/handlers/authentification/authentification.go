package authentificationhandlers

import (
	"bank_app/internal/api/jwt"
	"bank_app/internal/monitoring"
	"bank_app/internal/services"
	"bank_app/pkg/logger"
)

type AuthHandler struct {
	userService services.UsersServiceInterface
	jwtService  jwt.TokenService
	logger      logger.Logger
	metrics     *monitoring.Metrics
}

func NewAuthHandler(
	userService services.UsersServiceInterface,
	jwtService jwt.TokenService,
	logger logger.Logger,
	metrics *monitoring.Metrics,
) *AuthHandler {
	return &AuthHandler{
		userService: userService,
		jwtService:  jwtService,
		logger:      logger,
		metrics:     metrics,
	}
}
