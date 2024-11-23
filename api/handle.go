package receipt

import (
	"fetch-takehome/pkg/receipt"
	"net/http"

	"github.com/go-chi/render"
)

type Handle struct {
	ReceiptService receipt.ReceiptModule
}

func (h Handle) HelloPublic(w http.ResponseWriter, r *http.Request) {
	render.PlainText(w, r, "Hello")
}

func ErrorResponse(statusCode int, message string) *Response {
	return &Response{
		Code:        statusCode,
		body:        message,
		contentType: "application/json",
	}

}

func (h Handle) GetReceiptsIDPoints(w http.ResponseWriter, r *http.Request, id string) *Response {
	return ErrorResponse(http.StatusNotImplemented, "GetReceiptPoint not implemented")
}

func (h Handle) PostReceiptsProcess(w http.ResponseWriter, r *http.Request) *Response {
	return ErrorResponse(http.StatusNotImplemented, "PostReceipt not implemented")
}
