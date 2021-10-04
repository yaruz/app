package yaruserror

import (
	"github.com/pkg/errors"
)

// ErrNotSet is error for case when some thing is not set
var ErrNotSet error = errors.New("Not set.")

// ErrNotFound is error for case when entity not found
var ErrNotFound error = errors.New("Not found.")

// ErrAlreadyExists is error for case when some value is already exists
var ErrAlreadyExists error = errors.New("Already exists.")

// ErrEmptyParams is error for case when params are empty
var ErrEmptyParams error = errors.New("Empty params.")

// ErrBadRequest is error for case when bad request
var ErrBadParams error = errors.New("Bad params.")

// ErrBadRequest is error for case when bad request
var ErrBadRequest error = errors.New("Bad request.")

// ErrInternal is error for case when smth went wrong
var ErrInternal error = errors.New("Internal error.")

var ErrTokenHasExpired error = errors.New("Token has expired.")
