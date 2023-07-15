package mindrpc

import (
	"fmt"
)

func ErrUnauthorized(message string) error {

	return fmt.Errorf("unauthorized: %s", message)
}
