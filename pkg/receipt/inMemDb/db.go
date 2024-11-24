package db

import (
	"time"

	"github.com/google/uuid"
)

type Item struct {
	Price            float64
	ShortDescription string
}

type Receipt struct {
	ReceiptUuid  uuid.UUID
	Total        float64
	Retailer     string
	PurchaseTime time.Time
	Items        []Item
}

type InMemDB struct {
	Receipts map[uuid.UUID]Receipt
}

func New() *InMemDB {

	receipts := make(map[uuid.UUID]Receipt)

	db := InMemDB{
		Receipts: receipts,
	}
	return &db
}
