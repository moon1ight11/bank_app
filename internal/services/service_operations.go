package services

import "bank_app/internal/storage/repos/operations"

type OperationsService struct {
	operationsRepo *operations.Repo
}

func NewOperationsService(operationsRepo *operations.Repo) *OperationsService {
	return &OperationsService{operationsRepo: operationsRepo}
}