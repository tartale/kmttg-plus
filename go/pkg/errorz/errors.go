package errorz

import (
	"errors"
	"fmt"
)

func ErrNotAuthenticated(message string) error {

	return fmt.Errorf("authentication failed: %s", message)
}

func ErrResponse(message string) error {

	return fmt.Errorf("error in tivo response: %s", message)
}

var ErrReconnected error = errors.New("connection was restarted")

func ErrBadRequest(message string) error {

	return fmt.Errorf("bad request: %s", message)
}
