package beacon

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestBeacon(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Beacon Test Suite")
}
