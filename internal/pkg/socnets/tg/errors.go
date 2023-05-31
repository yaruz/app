package tg

import (
	"errors"
)

// ErrNotFound is error for case when entity not found
var ErrNotFound error = errors.New("Not found.")

// ErrBadRequest is error for case when bad request
var ErrBadParams error = errors.New("Bad params.")

// ErrBadRequest is error for case when bad request
var ErrBadRequest error = errors.New("Bad request.")

// ErrInternal is error for case when smth went wrong
var ErrInternal error = errors.New("Internal error.")

var ErrSignUpRequired error = errors.New("Sign up required.")

var ErrPasswordAuthNeeded error = errors.New("Password auth needed.")

var ErrUserEmpty error = errors.New("User is empty.")
