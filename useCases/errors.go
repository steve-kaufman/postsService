package useCases

import "errors"

var ErrInternal = errors.New("internal error")
var ErrNotFound = errors.New("post not found")
var ErrCantChangeLikes = errors.New("likes cant be changed")
