package mindrpc

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/tartale/kmttg-plus/go/pkg/config"
)

var _ = Describe("Tivo Client", func() {

	It("can create a TLS config from the certificates", func() {
		tlsConfig, err := tlsConfigFromCertificates()
		Expect(err).ToNot(HaveOccurred())
		Expect(tlsConfig).NotTo(BeNil())
	})

	It("can create a new tivo RPC client", func() {
		if testTivo == nil {
			Skip("skipping test; to enable, populate the KMTTG_TEST_TIVO env variable with information for an existing Tivo")
		}

		client, err := NewTivoClient(testTivo.Address)
		Expect(err).ToNot(HaveOccurred())
		Expect(client).NotTo(BeNil())

		err = client.Close()
		Expect(err).ToNot(HaveOccurred())
	})

	PIt("can authenticate a client session", func() {
		if testTivo == nil {
			Skip("skipping test; to enable, populate the KMTTG_TEST_TIVO env variable with information for an existing Tivo")
		}

		client, err := NewTivoClient(testTivo.Address)
		Expect(err).ToNot(HaveOccurred())
		Expect(client).NotTo(BeNil())
		defer client.Close()

		authMessage := NewTivoMessage().WithAuthPayload(config.Values.MediaAccessKey)
		err = client.SendRequest(*authMessage)
		Expect(err).ToNot(HaveOccurred())

		client.ReceiveResponse(context.Background())
	})
})
