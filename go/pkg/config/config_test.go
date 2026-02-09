package config

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Config", func() {

	It("can return the path to the certificate file", func() {
		certificatePath, err := CertificatePath()
		Expect(err).ToNot(HaveOccurred())
		Expect(certificatePath).To(HaveSuffix(certificateFilename))
	})
})
