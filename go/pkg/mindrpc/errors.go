package mindrpc

import (
	"fmt"
)

func ErrNotAuthenticated(message string) error {

	return fmt.Errorf("unauthorized: %s", message)
}
