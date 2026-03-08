package luna_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestLuna(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Luna Suite")
}
