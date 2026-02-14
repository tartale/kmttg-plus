package errorz

import (
	"errors"
)

var (
	ErrAuthenticationFailed       = errors.New("authentication failed")
	ErrReconnected          error = errors.New("connection was restarted")
	ErrUnexpectedResponse         = errors.New("unexpected response")
)
