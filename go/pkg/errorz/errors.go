package errorz

import (
	"errors"
)

var ErrAuthenticationFailed = errors.New("authentication failed")
var ErrReconnected error = errors.New("connection was restarted")
var ErrUnexpectedResponse = errors.New("unexpected response")
