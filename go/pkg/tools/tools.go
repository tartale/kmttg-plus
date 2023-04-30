//go:build tools

package tools

import (
	_ "github.com/99designs/gqlgen"
	_ "mvdan.cc/gofumpt"
)

//go:generate go install github.com/99designs/gqlgen
