package nef_test

import (
	"os"

	. "github.com/onsi/ginkgo/v2" //nolint:revive // ok
	. "github.com/onsi/gomega"    //nolint:revive // ok
	nef "github.com/snivilised/nefilim"
	lab "github.com/snivilised/nefilim/internal/laboratory"
	"github.com/snivilised/nefilim/test/luna"
)

var _ = Describe("op: create", Ordered, func() {
	var root string

	BeforeAll(func() {
		root = luna.Repo("test")
	})

	Context("fs: WriteFileFS", func() {
		BeforeEach(func() {
			scratch(root)
		})

		Context("overwrite", func() {
			var fS nef.WriteFileFS

			BeforeEach(func() {
				fS = nef.NewWriteFileFS(nef.Rel{
					Root:      root,
					Overwrite: true,
				})
				Expect(fS.IsRelative()).To(BeTrue())
			})

			Context("op: Create", func() {
				When("given: file does not already exist", func() {
					It("🧪 should: create successfully", func() {
						Expect(require(root, lab.Static.FS.Scratch)).To(Succeed())
						name := lab.Static.FS.Create.Destination
						file, err := fS.Create(name)
						Expect(err).To(Succeed())
						defer file.Close()

						Expect(luna.AsFile(name)).To(luna.ExistInFS(fS))
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

						Expect(luna.AsFile(name)).To(luna.ExistInFS(fS))
					})
				})
			})
		})

		Context("tentative", func() {
			var fS nef.WriteFileFS

			BeforeEach(func() {
				fS = nef.NewWriteFileFS(nef.Rel{
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
