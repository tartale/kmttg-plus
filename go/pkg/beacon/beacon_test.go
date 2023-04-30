package beacon

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Tivo Beacon", func() {

	It("can listen for Tivos on the network", func() {
		// TODO: determine a way to mock the mDNS traffic
		Expect(true)
	})
})
