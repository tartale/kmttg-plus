package errorz

import (
	"fmt"
)

func ErrNotAuthenticated(message string) error {

	return fmt.Errorf("authentication failed: %s", message)
}

func ErrResponse(message string) error {

	return fmt.Errorf("error in tivo response: %s", message)
}
