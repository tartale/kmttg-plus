package mindrpc

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Tivo Message", func() {

	It("headers can be converted to a string in wire format", func() {
		headers := make(TivoMessageHeaders)

		headers.Set("Foo", "bar")
		headers.Set("BarBaz", "boo")

		headerString := headers.String()
		Expect(headerString).To(Equal("BarBaz: boo\r\nFoo: bar\r\n\r\n"))
		Expect(headerString).To(HaveLen(25))
	})

	FIt("can be converted to Mind RPC format", func() {
		headers := make(TivoMessageHeaders)

		headers.Set("Foo", "bar")
		headers.Set("BarBaz", "boo")

		headerString := headers.String()
		Expect(headerString).To(Equal("BarBaz: boo\r\nFoo: bar\r\n\r\n"))
		Expect(headerString).To(HaveLen(25))
	})
})
