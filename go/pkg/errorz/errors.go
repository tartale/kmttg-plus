package errorz

import (
	"errors"
)

var ErrAuthenticationFailed = errors.New("authentication failed")
var ErrReconnected error = errors.New("connection was restarted")
var ErrResponse = errors.New("response error")
var ErrInvalidArgument = errors.New("invalid argument")
var ErrBadRequest = errors.New("bad request")
var ErrNotFound = errors.New("not found")
