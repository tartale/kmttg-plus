package mindrpc

import (
	"fmt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/tartale/kmttg-plus/go/pkg/config"
)

var _ = Describe("Tivo Message", func() {

	It("headers can be converted to a string in wire format", func() {
		headers := make(TivoMessageHeaders)

		headers.Set("Foo", "bar")
		headers.Set("BarBaz", "boo")

		headerString := headers.String()
		Expect(headerString).To(Equal("BarBaz: boo\r\nFoo: bar\r\n"))
		Expect(headerString).To(HaveLen(23))
	})

	It("can be populated as an authentication payload", func() {
		tivoMessage := NewTivoMessage().WithAuthPayload(config.Values.MediaAccessKey)
		Expect(tivoMessage).ToNot(BeNil())
		Expect(tivoMessage.Payload.Credential.Key).To(Equal(config.Values.MediaAccessKey))
	})

	It("can be formatted as a MindRPC message", func() {
		tivoMessage := NewTivoMessage().WithAuthPayload(config.Values.MediaAccessKey)
		mindRpcMessage, err := tivoMessage.ToMindRpcMessage()
		fmt.Println(mindRpcMessage)
		Expect(err).ToNot(HaveOccurred())
	})
})
