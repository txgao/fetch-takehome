package db

import (
	"context"

	"github.com/google/uuid"
)

const getItemsByReceipt = `-- name: GetItemsByReceipt :many
SELECT i.item_uuid, i.price, i.short_description
FROM receipt_items ri, item i
WHERE ri.item_uuid = i.item_uuid
AND ri.receipt_uuid = $1
`

func (q *Queries) GetItemsByReceipt(ctx context.Context, receiptUuid uuid.UUID) ([]Item, error) {
	rows, err := q.db.Query(ctx, getItemsByReceipt, receiptUuid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Item
	for rows.Next() {
		var i Item
		if err := rows.Scan(&i.ItemUuid, &i.Price, &i.ShortDescription); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getReceiptById = `-- name: GetReceiptById :one
SELECT receipt_uuid, total, purchase_date, purchase_time, retailer
FROM receipt
WHERE receipt_uuid = $1
`

func (q *Queries) GetReceiptById(ctx context.Context, receiptUuid uuid.UUID) (Receipt, error) {
	row := q.db.QueryRow(ctx, getReceiptById, receiptUuid)
	var i Receipt
	err := row.Scan(
		&i.ReceiptUuid,
		&i.Total,
		&i.PurchaseDate,
		&i.PurchaseTime,
		&i.Retailer,
	)
	return i, err
}
