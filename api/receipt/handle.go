package receipt

import (
	"errors"
	errorcode "fetch-takehome/pkg/errors"
	"fetch-takehome/pkg/receipt"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/render"
	"github.com/google/uuid"
)

type Handle struct {
	ReceiptService receipt.ReceiptModule
}

func ErrorResponse(statusCode int, message string) *Response {
	return &Response{
		Code: statusCode,
		body: map[string]interface{}{
			"code":    statusCode,
			"message": message,
		},
		contentType: "application/json",
	}

}

func (h Handle) GetIDPoints(w http.ResponseWriter, r *http.Request, id string) *Response {

	receipt_uuid, err := uuid.Parse(id)
	if err != nil {
		slog.Error("Error parsing receipt uuid", "err", err, "id", id)
		return ErrorResponse(http.StatusBadRequest, "unable to parse receipt id")
	}

	points, err := h.ReceiptService.GetReceiptPoint(r.Context(), receipt_uuid)
	if err != nil {
		slog.Error("Error finding receipt by uuid", "err", err, "id", receipt_uuid)
		if errors.Is(err, errorcode.ErrReceiptNotFound) {
			return ErrorResponse(http.StatusNotFound, "receipt not found")
		}
		return ErrorResponse(http.StatusInternalServerError, "unable to get receipt")
	}

	body := struct {
		Points *int64 `json:"points,omitempty"`
	}{
		Points: &points,
	}
	return GetIDPointsJSON200Response(body)
}

func (h Handle) PostProcess(w http.ResponseWriter, r *http.Request) *Response {

	data := PostProcessJSONBody{}
	err := render.DecodeJSON(r.Body, &data)
	if err != nil {
		return ErrorResponse(http.StatusBadRequest, "unable to decode body")
	}

	// Parse the date time into time.Time
	purchaseDateTime := strings.TrimSpace(data.PurchaseDate.String()) + " " + strings.TrimSpace(data.PurchaseTime)
	layout := "2006-01-02 15:04"
	timeStamp, err := time.Parse(layout, purchaseDateTime)
	if err != nil {
		slog.Error("Error parsing purchase date time", "err", err, "date", data.PurchaseDate.String(), "time", data.PurchaseTime)
		return ErrorResponse(http.StatusBadRequest, "unable to parse purchase date time")
	}

	// Parse the total value
	floatTotal, err := verifyFloatValue(data.Total)
	if err != nil {
		slog.Error("Error parsing total value", "err", err, "total", data.Total)
		return ErrorResponse(http.StatusBadRequest, "unable to parse total")
	}

	// Parse items into struct
	items := []receipt.Item{}
	for _, item := range data.Items {
		floatPrice, err := verifyFloatValue(item.Price)
		if err != nil {
			slog.Error("Error parsing price value", "err", err, "price", item.Price)
			return ErrorResponse(http.StatusBadRequest, "unable to parse item price")
		}
		temp := receipt.Item{
			Price:            floatPrice,
			ShortDescription: strings.TrimSpace(item.ShortDescription),
		}
		items = append(items, temp)
	}

	// create receipt
	receipt, err := h.ReceiptService.CreateReceipt(r.Context(), receipt.CreateReceiptParams{
		Items:        items,
		Total:        floatTotal,
		PurchaseTime: timeStamp,
		Retailer:     strings.TrimSpace(data.Retailer),
	})
	if err != nil {
		slog.Error("Error creating receipt", "err", err)
		return ErrorResponse(http.StatusInternalServerError, "fail to create receipt")
	}

	body := struct {
		ID string `json:"id"`
	}{
		ID: receipt.String(),
	}

	return PostProcessJSON200Response(body)
}

// ensure value is not less than 0
func verifyFloatValue(s string) (float64, error) {
	float, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, err
	}
	if float < 0 {
		return 0, errors.New("invalud value")
	}
	return float, nil
}
