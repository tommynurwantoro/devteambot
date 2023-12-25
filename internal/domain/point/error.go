package point

import "errors"

var (
	// ErrDataNotFound .
	ErrDataNotFound = errors.New("point: data not found")
	ErrLimitReached = errors.New("point: limit reached")
)
