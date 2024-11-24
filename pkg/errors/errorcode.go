package errorcode

import "errors"

// define common errors
var (
	ErrReceiptNotFound = errors.New("receipt not found")
)
