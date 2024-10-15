package nef_test

import (
	"os"

	. "github.com/onsi/ginkgo/v2" //nolint:revive // ok
	. "github.com/onsi/gomega"    //nolint:revive // ok
	nef "github.com/snivilised/nefilim"
	lab "github.com/snivilised/nefilim/internal/laboratory"
)

var _ = Describe("op: create", Ordered, func() {
	var root string

	BeforeAll(func() {
		root = Repo("test")
	})

	Context("fs: WriteFileFS", func() {
		BeforeEach(func() {
			scratch(root)
		})

		Context("overwrite", func() {
			var fS nef.WriteFileFS

			BeforeEach(func() {
				fS = nef.NewWriteFileFS(nef.At{
					Root:      root,
					Overwrite: true,
				})
			})

			Context("op: Create", func() {
				When("given: file does not already exist", func() {
					It("🧪 should: create successfully", func() {
						Expect(require(root, lab.Static.FS.Scratch)).To(Succeed())
						name := lab.Static.FS.Create.Destination
						file, err := fS.Create(name)
						Expect(err).To(Succeed())
						defer file.Close()

						Expect(AsFile(name)).To(ExistInFS(fS))
					})
				})

				When("given: file exists", func() {
					It("🧪 should: create successfully", func() {
						Expect(require(
							root, lab.Static.FS.Scratch, lab.Static.FS.Create.Destination,
						)).To(Succeed())
						name := lab.Static.FS.Create.Destination
						file, err := fS.Create(name)
						Expect(err).To(Succeed())
						defer file.Close()

						Expect(AsFile(name)).To(ExistInFS(fS))
					})
				})
			})
		})

		Context("tentative", func() {
			var fS nef.WriteFileFS

			BeforeEach(func() {
				fS = nef.NewWriteFileFS(nef.At{
					Root:      root,
					Overwrite: false,
				})
			})

			Context("op: Create", func() {
				When("given: file exists", func() {
					It("🧪 should: fail", func() {
						file := lab.Static.FS.Create.Destination
						Expect(require(
							root, lab.Static.FS.Scratch, file,
						)).To(Succeed())
						_, err := fS.Create(file)
						Expect(err).To(MatchError(os.ErrExist))
					})
				})
			})
		})
	})
})
