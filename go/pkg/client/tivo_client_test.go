package client_test

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/tartale/kmttg-plus/go/pkg/client"
	"github.com/tartale/kmttg-plus/go/test"
)

var _ = Describe("Tivo Client", func() {

	It("can create a TLS config from the certificates", func() {
		if test.Tivo == nil {
			Skip("skipping test; to enable, populate the KMTTG_TEST_TIVO env variable with information for an existing Tivo")
		}

		tlsConfig, err := client.NewTLSConfig(test.Tivo)
		Expect(err).ToNot(HaveOccurred())
		Expect(tlsConfig).NotTo(BeNil())
		Expect(tlsConfig.ServerName).To(Equal(test.Tivo.ServerName()))
	})

	It("can create a new tivo RPC client", func() {
		if test.Tivo == nil {
			Skip("skipping test; to enable, populate the KMTTG_TEST_TIVO env variable with information for an existing Tivo")
		}

		tivoClient, err := client.NewTivoClient(test.Tivo)
		Expect(err).ToNot(HaveOccurred())
		Expect(tivoClient).NotTo(BeNil())

		err = tivoClient.Close()
		Expect(err).ToNot(HaveOccurred())
	})

	It("can authenticate a client session", func() {
		if test.Tivo == nil {
			Skip("skipping test; to enable, populate the KMTTG_TEST_TIVO env variable with information for an existing Tivo")
		}

		tivoClient, err := client.NewTivoClient(test.Tivo)
		Expect(err).ToNot(HaveOccurred())
		Expect(tivoClient).NotTo(BeNil())
		defer tivoClient.Close()

		err = tivoClient.Authenticate(context.Background())
		Expect(err).ToNot(HaveOccurred())
	})

	It("can get a list of recordings", func() {
		if test.Tivo == nil {
			Skip("skipping test; to enable, populate the KMTTG_TEST_TIVO env variable with information for an existing Tivo")
		}

		tivoClient, err := client.NewTivoClient(test.Tivo)
		Expect(err).ToNot(HaveOccurred())
		Expect(tivoClient).NotTo(BeNil())
		defer tivoClient.Close()

		err = tivoClient.Authenticate(context.Background())
		Expect(err).ToNot(HaveOccurred())

		recordings, err := tivoClient.GetRecordingList(context.Background())
		Expect(err).ToNot(HaveOccurred())
		Expect(len(recordings)).To(BeNumerically(">", 0))
	})
})
