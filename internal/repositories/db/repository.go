package db

import (
	"context"
	"docmate/internal/model"
	"docmate/internal/services/transaction"
	"log/slog"

	"gorm.io/gorm"
)

type Repository struct {
	client *gorm.DB
}

func NewRepository(client *gorm.DB) *Repository {
	return &Repository{
		client: client,
	}
}

func (r *Repository) CreateTransaction(ctx context.Context) (*model.TXClient, error) {
	tx := r.client.Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}
	slog.Info("transaction started...")

	return transaction.NewTXClient(ctx, tx), nil
}

func (r *Repository) dbClient(txc *model.TXClient) *gorm.DB {
	if txc == nil {
		return r.client
	}

	return txc.Get().(*gorm.DB)
}
