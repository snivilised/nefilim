package nef_test

import (
	. "github.com/onsi/ginkgo/v2" //nolint:revive // ok
	. "github.com/onsi/gomega"    //nolint:revive // ok
	nef "snivilised.com/nefilim"
)

var _ = Describe("Dummy", func() {
	Context("given: foo", func() {
		It("should: bar", func() {
			Expect(nef.Greet("nefilim")).To(Equal("greetings: 'nefilim'"))
		})
	})
})
