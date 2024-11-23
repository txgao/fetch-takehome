package receipt

import (
	"fetch-takehome/pkg/receipts"
	"net/http"
)

type Handle struct {
	ReceiptService receipts.ReceiptModule
}

func ErrorResponse(statusCode int, message string) *Response {
	return &Response{
		Code:        statusCode,
		body:        message,
		contentType: "application/json",
	}

}

func (h Handle) GetIDPoints(w http.ResponseWriter, r *http.Request, id string) *Response {
	return ErrorResponse(http.StatusNotImplemented, "GetReceiptPoint not implemented")
}

func (h Handle) PostProcess(w http.ResponseWriter, r *http.Request) *Response {
	return ErrorResponse(http.StatusNotImplemented, "PostReceipt not implemented")
}
