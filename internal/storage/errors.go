package storage

import "errors"

var ErrNotFound = errors.New("not found")
var ErrAliasConflict = errors.New("alias conflict")
var ErrInvalidAlias = errors.New("invalid alias")
