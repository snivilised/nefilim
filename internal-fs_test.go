package nef

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("internal-fs", func() {
	Context("Calc", Label("coverage"), func() {
		Context("existsInFS disambiguator", func() {
			It("🧪 should: return the calc", func() {
				entities := compose("root").mutate(true)
				Expect(entities.writer.existsInFS.Calc()).NotTo(BeNil())
			})
		})

		Context("makeDirAllFS disambiguator", func() {
			It("🧪 should: return the calc", func() {
				entities := compose("root").mutate(true)
				Expect(entities.writer.makeDirAllFS.Calc()).NotTo(BeNil())
			})
		})

		Context("readerFS disambiguator", func() {
			It("🧪 should: return the calc", func() {
				entities := compose("root").mutate(true)
				Expect(entities.reader.Calc()).NotTo(BeNil())
			})
		})

		Context("writerFS disambiguator", func() {
			It("🧪 should: return the calc", func() {
				entities := compose("root").mutate(true)
				Expect(entities.writer.Calc()).NotTo(BeNil())
			})
		})

		Context("openFS disambiguator", func() {
			It("🧪 should: return the calc", func() {
				entities := compose("root").mutate(true)
				Expect(entities.open.Calc()).NotTo(BeNil())
			})
		})
	})
})
