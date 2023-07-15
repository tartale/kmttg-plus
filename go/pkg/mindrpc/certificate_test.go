package mindrpc

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/tartale/kmttg-plus/go/pkg/config"
	"github.com/tartale/kmttg-plus/go/test"
)

var _ = Describe("Tivo Certificates", func() {

	It("can be read from the configured certificate file", func() {
		if test.Tivo == nil {
			Skip("skipping test; to enable, populate the KMTTG_TEST_TIVO env variable with information for an existing Tivo")
		}

		certificatePath, err := config.CertificatePath()
		Expect(err).ToNot(HaveOccurred())

		certs, certPool, err := GetCertificates(certificatePath)
		Expect(err).ToNot(HaveOccurred())
		Expect(certs).ToNot(BeNil())
		Expect(certs.Certificate).To(HaveLen(3))
		Expect(certs.PrivateKey).ToNot(BeNil())
		Expect(certPool).ToNot(BeNil())
	})
})
