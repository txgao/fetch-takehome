package receipt

import (
	"context"
	"log/slog"
	"math"
	"time"
	"unicode"

	errorcode "fetch-takehome/pkg/errors"
	"fetch-takehome/pkg/receipt/db"

	inMemDb "fetch-takehome/pkg/receipt/inMemDb"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type (
	CreateReceiptParams struct {
		Total        float64
		PurchaseTime time.Time
		Retailer     string
		Items        []Item
	}

	Item struct {
		Price            float64
		ShortDescription string
	}

	ReceiptModule interface {
		CreateReceipt(ctx context.Context, params CreateReceiptParams) (uuid.UUID, error)
		GetReceiptPoint(ctx context.Context, receipt_id uuid.UUID) (int64, error)
	}
)

type ReceiptService struct {
	pgDb  *db.Queries
	inMem *inMemDb.InMemDB
}

const (
	PointPerAlphanumeric       = 1
	PointNoCents               = 50
	PointMultipleQuater        = 25
	PointOddDate               = 6
	PointTimeInRange           = 10
	PointEveryTwoItems         = 5
	PointItemPriceMultiplyRate = 0.2
	TwoPMHour                  = 14
	FourPMHour                 = 16
)

type Option func(*ReceiptService)

func WithDB(dbPool *pgxpool.Pool) Option {
	return func(rs *ReceiptService) {
		rs.pgDb = db.New(dbPool)
	}
}

func NewService(options ...Option) *ReceiptService {

	rs := &ReceiptService{
		inMem: inMemDb.New(),
	}
	for _, opt := range options {
		opt(rs)
	}

	return rs

}

func (receiptSvc *ReceiptService) createReceiptInMem(ctx context.Context, params CreateReceiptParams) (uuid.UUID, error) {

	uid := uuid.New()
	receipt := inMemDb.Receipt{
		ReceiptUuid:  uid,
		Total:        params.Total,
		Retailer:     params.Retailer,
		PurchaseTime: params.PurchaseTime,
		Items:        []inMemDb.Item{},
	}
	for _, item := range params.Items {
		receipt.Items = append(receipt.Items, inMemDb.Item{
			Price:            item.Price,
			ShortDescription: item.ShortDescription,
		})
	}

	receiptSvc.inMem.Receipts[uid] = receipt

	return receiptSvc.inMem.Receipts[uid].ReceiptUuid, nil
}

func (receiptSvc *ReceiptService) createReceiptInDb(ctx context.Context, params CreateReceiptParams) (uuid.UUID, error) {

	receipt_uuid, err := receiptSvc.pgDb.CreateReceipt(ctx, db.CreateReceiptParams{
		Total:        params.Total,
		PurchaseTime: params.PurchaseTime,
		Retailer:     params.Retailer,
	})
	if err != nil {
		slog.Error("fail to create receipt", "err", err)
		return uuid.Nil, err
	}
	for _, item := range params.Items {
		item_uuid, err := receiptSvc.pgDb.CreateItem(ctx, db.CreateItemParams{
			Price:            item.Price,
			ShortDescription: stringToPgText(item.ShortDescription),
		})
		if err != nil {
			slog.Error("fail to create item", "err", err)
			return uuid.Nil, err
		}

		_, err = receiptSvc.pgDb.CreateReceiptItem(ctx, db.CreateReceiptItemParams{
			ItemUuid:    item_uuid,
			ReceiptUuid: receipt_uuid,
		})
		if err != nil {
			slog.Error("fail to create receipt item", "err", err)
			return uuid.Nil, err
		}
	}

	return receipt_uuid, nil
}

func (receiptSvc *ReceiptService) CreateReceipt(ctx context.Context, params CreateReceiptParams) (uuid.UUID, error) {
	if receiptSvc.pgDb != nil {
		return receiptSvc.createReceiptInDb(ctx, params)
	} else {
		return receiptSvc.createReceiptInMem(ctx, params)
	}
}

func (receiptSvc *ReceiptService) getReceiptPointInMem(ctx context.Context, receipt_id uuid.UUID) (int64, error) {
	result := int64(0)

	receipt, ok := receiptSvc.inMem.Receipts[receipt_id]
	if !ok {
		return result, errorcode.ErrReceiptNotFound
	}

	result += calculateReceiptPoints(receipt)
	result += calculateItemsPoints(receipt.Items)

	return result, nil
}

func (receiptSvc *ReceiptService) getReceiptPointInDb(ctx context.Context, receipt_id uuid.UUID) (int64, error) {

	return int64(0), nil
}

func (receiptSvc *ReceiptService) GetReceiptPoint(ctx context.Context, receipt_id uuid.UUID) (int64, error) {

	if receiptSvc.pgDb != nil {
		return receiptSvc.getReceiptPointInDb(ctx, receipt_id)
	} else {
		return receiptSvc.getReceiptPointInMem(ctx, receipt_id)
	}
}

func countAlphanumeric(s string) int64 {
	count := int64(0)
	for _, char := range s {
		if unicode.IsLetter(char) || unicode.IsDigit(char) {
			count++
		}
	}
	return count
}

func calculateReceiptPoints(receipt inMemDb.Receipt) int64 {

	points := int64(0)

	// One point for every alphanumeric character in the retailer name.
	count := countAlphanumeric(receipt.Retailer)
	points += count * PointPerAlphanumeric
	slog.Info("alphanumeric character", "points", count*PointPerAlphanumeric)

	// 50 points if the total is a round dollar amount with no cents.
	_, frac := math.Modf(receipt.Total)
	if frac == 0 {
		points += PointNoCents
		slog.Info("total no cents", "points", PointNoCents)
	}

	// 25 points if the total is a multiple of 0.25.
	if receipt.Total != 0 && math.Mod(frac, 0.25) == 0 {
		points += PointMultipleQuater
		slog.Info("total multiple of 0.25", "points", PointMultipleQuater)
	}

	// 6 points if the day in the purchase date is odd.
	day := receipt.PurchaseTime.Day()
	if day%2 == 1 {
		points += PointOddDate
		slog.Info("odd date", "points", PointOddDate)
	}

	// 10 points if the time of purchase is after 2:00pm and before 4:00pm.
	hour := receipt.PurchaseTime.Hour()
	minute := receipt.PurchaseTime.Minute()

	if (hour == TwoPMHour && minute > 0) || (hour > TwoPMHour && hour < FourPMHour) {
		points += PointTimeInRange
		slog.Info("time in range", "points", PointTimeInRange)
	}

	return points
}

func calculateItemsPoints(items []inMemDb.Item) int64 {
	points := int64(0)

	//5 points for every two items on the receipt.
	count := int64(len(items) / 2)
	points += count * PointEveryTwoItems
	slog.Info("item count", "points", count*PointEveryTwoItems)

	// If length of the item description is a multiple of 3, multiply the price by 0.2 and round up
	for _, item := range items {
		slog.Info("	item:", "description", item.ShortDescription, "price", item.Price)

		if len(item.ShortDescription) != 0 && len(item.ShortDescription)%3 == 0 {
			cur_point := PointItemPriceMultiplyRate * item.Price
			points += int64(math.Ceil(cur_point))
			slog.Info("	multiple of 3", "points", math.Ceil(cur_point))
		}
	}

	return points

}

func stringToPgText(s string) pgtype.Text {
	return pgtype.Text{
		String: s,
		Valid:  true,
	}
}
