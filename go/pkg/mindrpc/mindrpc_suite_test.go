package mindrpc

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestMindRPC(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Mind RPC Test Suite")
}
