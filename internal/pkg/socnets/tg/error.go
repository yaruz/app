package tg

import (
	"github.com/pkg/errors"
)

// ErrSessionPasswordNeeded is error for case when session password needed
var ErrSessionPasswordNeeded error = errors.New("Session password needed.")
