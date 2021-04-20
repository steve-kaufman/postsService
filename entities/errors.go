package entities

import "errors"

var ErrNeedsTitle = errors.New("title is required")
var ErrTooLong = errors.New("content must be less than 500 characters")
