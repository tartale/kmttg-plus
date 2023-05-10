package mindrpc

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/tartale/kmttg-plus/go/pkg/model"
)

var testTivo *model.Tivo

func TestMindRPC(t *testing.T) {
	initializeTestTivo()
	RegisterFailHandler(Fail)
	RunSpecs(t, "Mind RPC Test Suite")
}

func initializeTestTivo() {
	testTivoJSON := os.Getenv("KMTTG_TEST_TIVO")
	if testTivoJSON == "" {
		return
	}
	testTivo = &model.Tivo{}
	err := json.Unmarshal([]byte(testTivoJSON), testTivo)
	if err != nil {
		panic(fmt.Errorf("%s: %s", "error trying to initialize test tivo", err))
	}
}
