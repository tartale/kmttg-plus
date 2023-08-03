package message

import (
	"github.com/tartale/go/pkg/jsontime"
	"github.com/tartale/kmttg-plus/go/pkg/logz"
)

func init() {

	jsontime.Logger = logz.NopLogger
}
