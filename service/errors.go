package service

import "errors"

var ErrInternal = errors.New("internal error")
var ErrNotFound = errors.New("post not found")
var ErrNeedsTitle = errors.New("title is required")
var ErrTooLong = errors.New("content must be less than 500 characters")
var ErrCantChangeLikes = errors.New("likes cant be changed")
