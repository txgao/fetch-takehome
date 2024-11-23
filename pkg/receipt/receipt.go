package receipt

import (
	"context"
	"time"

	"fetch-takehome/pkg/receipt/db"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type (
	CreateReceiptParams struct {
		Total        float64
		PurchaseDate time.Time
		PurchaseTime time.Time
		Retailer     string
	}

	ReceiptModule interface {
		CreateReceipt(ctx context.Context, params CreateReceiptParams) (uuid.UUID, error)
	}
)

type ReceiptService struct {
	pgDb *db.Queries
}

func NewService(dbPool *pgxpool.Pool) *ReceiptService {
	DB := db.New(dbPool)

	rs := &ReceiptService{
		pgDb: DB,
	}

	return rs

}

func (receiptSvc *ReceiptService) CreateReceipt(ctx context.Context, params CreateReceiptParams) (uuid.UUID, error) {
	return uuid.Nil, nil
}
