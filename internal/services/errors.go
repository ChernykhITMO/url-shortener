package services

import "errors"

var ErrInvalidURL = errors.New("invalid URL")
var ErrNotFound = errors.New("not found")
var ErrInvalidAlias = errors.New("invalid alias")
var ErrAliasTaken = errors.New("alias taken")
