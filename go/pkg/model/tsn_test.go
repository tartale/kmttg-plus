package model

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Tivo Helper", func() {

	It("can format a legal TSN as a ServerName", func() {
		testTivo := Tivo{
			Tsn: "012345678901234",
		}
		Expect(testTivo.ServerName()).To(Equal("012-3456-7890-1234"))
	})
})
