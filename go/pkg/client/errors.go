package client

import (
	"fmt"
)

func ErrNotAuthenticated(message string) error {

	return fmt.Errorf("authentication failed: %s", message)
}
