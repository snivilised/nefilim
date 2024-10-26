package nef_test

import (
	. "github.com/onsi/ginkgo/v2" //nolint:revive // ok

	nef "github.com/snivilised/nefilim"
	"github.com/snivilised/nefilim/test/luna"
)

var _ = Describe("op: copy/all", Ordered, func() {
	var (
		root   string
		fS     nef.UniversalFS
		single string
	)
	BeforeAll(func() {
		root = luna.Repo("test")
	})

	BeforeEach(func() {
		fS = nef.NewUniversalFS(nef.Rel{
			Root:      root,
			Overwrite: false,
		})
		scratch(root)
	})

	Context("op: Copy", func() {
		When("given: ", func() {
			It("🧪 should: ", func() {
				_ = fS
				_ = single
			})
		})
	})

	Context("op: CopyAll", func() {
		When("given: ", func() {
			It("🧪 should: ", func() {

			})
		})
	})
})
