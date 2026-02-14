package certificate_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/tartale/kmttg-plus/go/pkg/certificate"
)

var _ = Describe("Tivo Certificates", func() {
	It("can be read from the configured certificate file", func() {
		certs, certPool, err := certificate.Read()
		Expect(err).ToNot(HaveOccurred())
		Expect(certs).ToNot(BeNil())
		Expect(certs.Certificate).To(HaveLen(3))
		Expect(certs.PrivateKey).ToNot(BeNil())
		Expect(certPool).ToNot(BeNil())
	})
})
