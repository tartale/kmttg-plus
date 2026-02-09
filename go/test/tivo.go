package test

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/tartale/kmttg-plus/go/pkg/model"
)

var Tivo *model.Tivo

func init() {
	testTivoJSON := os.Getenv("KMTTG_TEST_TIVO")
	if testTivoJSON == "" {
		return
	}
	Tivo = &model.Tivo{}
	err := json.Unmarshal([]byte(testTivoJSON), Tivo)
	if err != nil {
		panic(fmt.Errorf("%s: %s", "error trying to initialize test tivo", err))
	}
}
