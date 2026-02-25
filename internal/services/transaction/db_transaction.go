package transaction

import (
	"context"
	"docmate/internal/model"
	"log/slog"

	"gorm.io/gorm"
)

type DBTransactionService struct {
	repo model.TXRepo
}

func NewTXClient(ctx context.Context, client *gorm.DB) *model.TXClient {
	return &model.TXClient{
		Ctx:    ctx,
		Client: client,
	}
}

func NewDBTransaction(repo model.TXRepo) *DBTransactionService {
	return &DBTransactionService{
		repo: repo,
	}
}

func (svc *DBTransactionService) CreateTransaction(ctx context.Context) (*model.TXClient, error) {
	txc, err := svc.repo.CreateTransaction(ctx)
	if err != nil {
		slog.Error("transaction creation failed", "error", err.Error())

		return nil, err
	}

	return txc, nil
}
