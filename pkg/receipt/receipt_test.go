package receipt

import (
	"context"
	"log/slog"
	"os"
	"testing"
	"time"

	inMemDb "fetch-takehome/pkg/receipt/inMemDb"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
)

var (
	dbPool *pgxpool.Pool
)

func setup() {
	var err error
	connStr := "postgres://content:pwd@localhost:5433/powercard_db"
	dbPool, err = pgxpool.New(context.Background(), connStr)
	if err != nil {
		slog.Error("Failed to connect to the database", "err", err)
	}
}

func teardown() {
	dbPool.Close()
}

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	teardown()
	os.Exit(code)
}

func TestCreateReceipt(t *testing.T) {

	rs := NewService(dbPool)

	receipt, err := rs.CreateReceipt(context.Background(), CreateReceiptParams{
		Total:        10.01,
		PurchaseTime: time.Now(),
		Retailer:     "Target",
		Items: []Item{
			Item{
				Price:            10.01,
				ShortDescription: "Toothbrush",
			},
		},
	})

	get_receipt, ok := rs.inMem.Receipts[receipt]

	assert.NoError(t, err, "Expected no error from CreateReceipt")
	assert.NotEmpty(t, receipt, "Expected receipt uuid to be non-empty")
	assert.True(t, ok, "Expected receipt uuid to be found")
	assert.Equal(t, 1, len(get_receipt.Items), "Expected receipt item to be found")
	assert.Equal(t, "Target", get_receipt.Retailer, "Expected retailer to be correct")

}

func TestCalculateItemsPoints(t *testing.T) {

	items_one := []inMemDb.Item{
		inMemDb.Item{
			Price:            1.0,
			ShortDescription: "abc",
		},
		inMemDb.Item{
			Price:            1.0,
			ShortDescription: "abc",
		},
		inMemDb.Item{
			Price:            1.0,
			ShortDescription: "abc",
		},
	}

	items_two := []inMemDb.Item{
		inMemDb.Item{
			Price:            1.0,
			ShortDescription: "abcd",
		},
	}

	points_one := calculateItemsPoints(items_one)
	points_two := calculateItemsPoints(items_two)
	assert.Equal(t, int64(8), points_one, "Expected first points to be 8")
	assert.Equal(t, int64(0), points_two, "Expected second points to be 0")

}

func TestCalculateReceiptPoints(t *testing.T) {

	timestring := "2022-01-02 13:01"
	layout := "2006-01-02 15:04"
	timeStamp, _ := time.Parse(layout, timestring)

	receipt := inMemDb.Receipt{
		Total:        10.24,
		Retailer:     "@@$$",
		PurchaseTime: timeStamp,
	}

	timestring = "2022-01-01 14:01"
	timeStamp, _ = time.Parse(layout, timestring)
	receipt_two := inMemDb.Receipt{
		Total:        10.01,
		Retailer:     "@@$$",
		PurchaseTime: timeStamp,
	}

	timestring = "2022-01-02 13:01"
	timeStamp, _ = time.Parse(layout, timestring)
	receipt_three := inMemDb.Receipt{
		Total:        10.5,
		Retailer:     "aa",
		PurchaseTime: timeStamp,
	}

	points_one := calculateReceiptPoints(receipt)
	points_two := calculateReceiptPoints(receipt_two)
	points_three := calculateReceiptPoints(receipt_three)

	assert.Equal(t, int64(0), points_one, "Expected first points to be 8")
	assert.Equal(t, int64(16), points_two, "Expected second points to be 0")
	assert.Equal(t, int64(27), points_three, "Expected second points to be 0")

}
