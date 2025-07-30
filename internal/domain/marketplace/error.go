package marketplace

import "errors"

var (
	ErrOutOfStock          = errors.New("marketplace: out of stock")
	ErrInsufficientBalance = errors.New("marketplace: insufficient balance")
)
