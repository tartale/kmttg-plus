package mindrpc

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Tivo Message Headers", func() {

	It("can be converted to a string in mRPC format", func() {
		headers := make(TivoMessageHeaders)

		headers.Set("Type", "bar")
		headers.Set("RpcId", "boo")

		headerString := headers.String()
		Expect(headerString).To(Equal("Type: bar\r\nRpcId: boo\r\n\r\n"))
		Expect(headerString).To(HaveLen(25))
	})

})
