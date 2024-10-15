package nef_test

import (
	. "github.com/onsi/ginkgo/v2" //nolint:revive // ok

	nef "github.com/snivilised/nefilim"
)

var _ = Describe("op: copy/all", Ordered, func() {
	var (
		root   string
		fS     nef.UniversalFS
		single string
	)
	BeforeAll(func() {
		root = Repo("test")
	})

	BeforeEach(func() {
		fS = nef.NewUniversalFS(nef.At{
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
