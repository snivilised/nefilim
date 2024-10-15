package collections_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2" //nolint:revive // gomega ok
	. "github.com/onsi/gomega"    //nolint:revive // gomega ok
)

func TestCollections(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Collections Suite")
}
