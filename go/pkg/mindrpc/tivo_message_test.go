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

		headers.Set("Type", "bar")
		headers.Set("RpcId", "boo")

		headerString := headers.String()
		Expect(headerString).To(Equal("Type: bar\r\nRpcId: boo\r\n\r\n"))
		Expect(headerString).To(HaveLen(25))
	})

	It("can be populated as an authentication request payload", func() {
		tivoMessage := NewTivoMessage().WithAuthRequest(config.Values.MediaAccessKey)
		Expect(tivoMessage).ToNot(BeNil())
		Expect(tivoMessage.Payload.Credential.Key).To(Equal(config.Values.MediaAccessKey))
	})

	It("can be formatted as a MindRPC message", func() {
		tivoMessage := NewTivoMessage().WithAuthRequest(config.Values.MediaAccessKey)
		mindRpcMessage, err := tivoMessage.ToMindRpcMessage()
		fmt.Print(mindRpcMessage)
		Expect(err).ToNot(HaveOccurred())
	})
})
