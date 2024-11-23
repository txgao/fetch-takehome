package receipt

import (
	"context"

	"fetch-takehome/pkg/receipt/db"

	"github.com/google/uuid"
)

type (
	CreateReceiptParams struct {
		ContentType string
		Filename    string
		Category    string
		TenantUuid  uuid.UUID
		UserUuid    uuid.UUID
		CompanyUuid uuid.UUID
	}

	ReceiptModule interface {
		CreateReceipt(ctx context.Context, params CreateReceiptParams) (uuid.UUID, error)
	}
)

type ReceiptService struct {
	pgDb *db.Queries
}

func (receiptSvc *ReceiptService) CreateReceipt(ctx context.Context, params CreateReceiptParams) (uuid.UUID, error) {
	return uuid.Nil, nil
}
