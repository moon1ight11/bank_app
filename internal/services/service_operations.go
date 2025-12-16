package services

import (
	"bank_app/internal/storage/repos/operations"

	"github.com/google/uuid"
)

type OperationsService struct {
	operationsRepo *operations.Repo
}

func NewOperationsService(operationsRepo *operations.Repo) *OperationsService {
	return &OperationsService{operationsRepo: operationsRepo}
}

// получение всех операций пользователя
func (o *OperationsService) AllOperationsGet(userID uuid.UUID) ([]operations.Operation, error) {
	operations, err := o.operationsRepo.GetAllUsersOperations(userID)
	if err != nil {
		return nil, err
	}

	return operations, nil
}

// получение операций по счету
func (o *OperationsService) AccountOperationsGet(userID uuid.UUID, accountID uuid.UUID) ([]operations.Operation, error) {
	operations, err := o.operationsRepo.GetOperationsByAccount(userID, accountID)
	if err != nil {
		return nil, err
	}

	return operations, nil
}

// получение одной операции
func (o *OperationsService) OperationByIdGet(userID uuid.UUID, operationID uuid.UUID) (operations.Operation, error) {
	operation, err := o.operationsRepo.GetOperationByID(operationID, userID)
	if err != nil {
		return operations.Operation{}, err
	}

	return operation, nil
}