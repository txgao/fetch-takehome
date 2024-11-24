package receipt

import (
	"context"
	"errors"
	"math"
	"time"
	"unicode"

	"fetch-takehome/pkg/receipt/db"

	inMemDb "fetch-takehome/pkg/receipt/inMemDb"

	"github.com/google/uuid"
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
	PointMultipleQuater        = 30
	PointOddDate               = 6
	PointTimeInRange           = 10
	PointEveryTwoItems         = 5
	PointItemPriceMultiplyRate = 0.2
	TwoPMHour                  = 14
	FourPMHour                 = 16
)

func NewService(dbPool *pgxpool.Pool) *ReceiptService {
	DB := db.New(dbPool)

	rs := &ReceiptService{
		pgDb:  DB,
		inMem: inMemDb.New(),
	}

	return rs

}

func (receiptSvc *ReceiptService) CreateReceipt(ctx context.Context, params CreateReceiptParams) (uuid.UUID, error) {

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

	// create_param := db.CreateReceiptParams{
	// 	Total:        params.Total,
	// 	PurchaseTime: params.PurchaseTime,
	// 	Retailer:     params.Retailer,
	// }
	// receipt, err := receiptSvc.pgDb.CreateReceipt(ctx, create_param)
	// if err != nil {
	// 	return uuid.Nil, err
	// }
	// return receipt, nil

}

func (receiptSvc *ReceiptService) GetReceiptPoint(ctx context.Context, receipt_id uuid.UUID) (int64, error) {

	result := int64(0)

	receipt, ok := receiptSvc.inMem.Receipts[receipt_id]
	if !ok {
		return result, errors.New("Receipt not found")
	}

	result += calculateReceiptPoints(receipt)
	result += calculateItemsPoints(receipt.Items)

	return result, nil

	// result := int64(0)
	// recipt, err := receiptSvc.pgDb.GetReceiptById(ctx, receipt_id)
	// if err != nil {
	// 	if errors.Is(err, pgx.ErrNoRows) {
	// 		return result, err
	// 	}
	// 	return result, err
	// }
	// items, err := receiptSvc.pgDb.GetItemsByReceipt(ctx, receipt_id)
	// if err != nil {
	// 	return 0, err
	// }

	// result += calculateReceiptPoints(recipt)
	// result += calculateItemsPoints(items)

	// return result, nil

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

	// 50 points if the total is a round dollar amount with no cents.
	_, frac := math.Modf(receipt.Total)
	if frac == 0 {
		points += PointNoCents
	}

	// 25 points if the total is a multiple of 0.25.
	if math.Mod(frac, 0.25) == 0 {
		points += PointMultipleQuater
	}

	// 6 points if the day in the purchase date is odd.
	day := receipt.PurchaseTime.Day()
	if day%2 == 1 {
		points += PointOddDate
	}

	// 10 points if the time of purchase is after 2:00pm and before 4:00pm.
	hour := receipt.PurchaseTime.Hour()
	if hour > TwoPMHour && hour < FourPMHour {
		points += PointTimeInRange
	}

	return 0
}

func calculateItemsPoints(items []inMemDb.Item) int64 {
	points := int64(0)

	//5 points for every two items on the receipt.
	count := int64(len(items) / 2)
	points += count * PointEveryTwoItems

	// If length of the item description is a multiple of 3, multiply the price by 0.2 and round up
	for _, item := range items {

		if len(item.ShortDescription)%3 == 0 {
			cur_point := PointItemPriceMultiplyRate * item.Price
			points += int64(math.Ceil(cur_point))
		}
	}

	return points

}
