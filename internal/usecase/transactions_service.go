package usecase

import (
	"GoBank/internal/domain"
	"GoBank/internal/repository"
	"context"
	"time"
)

type TransactionService struct {
	repo repository.TransactionRepository
}

func NewTransactionService(repo repository.TransactionRepository) *TransactionService {
	var transactionService TransactionService = TransactionService{
		repo: repo,
	}
	return &transactionService
}

func (r *TransactionService) ValidateTransaction(transactionFromFront domain.TransactionFromFront, from int64) bool {

}

func (r *TransactionService) CreateTransaction(transactionFromFront domain.TransactionFromFront, from int64) error {
	to := transactionFromFront.To
	amount := transactionFromFront.Amount
	timeOfCreation := time.Now()
	status := "SUCCESS"
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	if err := r.repo.CreateTransaction(ctx, from, to, amount, timeOfCreation, status); err != nil {
		return err
	}
	return nil
}

func (r *TransactionService) GetAmountOfSender() {

}
