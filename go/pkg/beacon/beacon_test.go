package beacon

import (
	"context"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/tartale/kmttg-plus/go/pkg/model"
	"github.com/tartale/kmttg-plus/go/pkg/tivo"
	"github.com/tartale/kmttg-plus/go/test"
)

var _ = Describe("Tivo Beacon", func() {

	It("can listen for Tivos on the network", func() {
		if test.Tivo == nil {
			Skip("skipping test; to enable, populate the KMTTG_TEST_TIVO env variable with information for an existing Tivo")
		}

		ctx := context.Background()
		ctx, cancelFunc := context.WithCancel(ctx)
		defer cancelFunc()

		go Listen(ctx)

		Eventually(func() []*model.Tivo { return tivo.List() }).
			WithTimeout(10 * time.Second).
			WithPolling(1 * time.Second).
			Should(ContainElement(HaveField("Name", test.Tivo.Name)))
	})
})
