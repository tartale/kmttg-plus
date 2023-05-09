package mindrpc

import (
	"context"
	"os"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var (
	testTivoAddress = os.Getenv("KMTTG_TEST_TIVO_ADDRESS")
	testTivoMak     = os.Getenv("KMTTG_TEST_TIVO_MAK")
)

var _ = PDescribe("Tivo Client", func() {

	It("can create a TLS config from the certificates", func() {
		tlsConfig, err := tlsConfigFromCertificates()
		Expect(err).ToNot(HaveOccurred())
		Expect(tlsConfig).NotTo(BeNil())
	})

	It("can create a new tivo RPC client", func() {
		if testTivoAddress == "" {
			Skip("skipping test; to enable, provide a valid address for an existing Tivo")
		}

		client, err := NewTivoClient(testTivoAddress)
		Expect(err).ToNot(HaveOccurred())
		Expect(client).NotTo(BeNil())

		err = client.Close()
		Expect(err).ToNot(HaveOccurred())
	})

	It("can authenticate a client session", func() {
		if testTivoAddress == "" || testTivoMak == "" {
			Skip("skipping test; to enable, provide a valid address and media access key for an existing Tivo")
		}

		client, err := NewTivoClient(testTivoAddress)
		Expect(err).ToNot(HaveOccurred())
		Expect(client).NotTo(BeNil())
		defer client.Close()

		authMessage := NewTivoMessage().WithAuthPayload(testTivoMak)
		err = client.SendRequest(*authMessage)
		Expect(err).ToNot(HaveOccurred())

		client.ReceiveResponse(context.Background())
	})
})
