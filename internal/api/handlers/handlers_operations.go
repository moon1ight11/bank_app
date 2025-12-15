package handlers

import (
	"bank_app/internal/jwt"
	"bank_app/internal/services"
)

type OperationsHandler struct {
	operationsService *services.OperationsService
	jwtService jwt.TokenService
}

func NewOperationsHandler(operationsService *services.OperationsService, jwtService jwt.TokenService) *OperationsHandler {
	return &OperationsHandler{
		operationsService: operationsService,
		jwtService:  jwtService,
	}
}
