package nef_test

import (
	. "github.com/onsi/ginkgo/v2" //nolint:revive // ok
	. "github.com/onsi/gomega"    //nolint:revive // ok
	nef "github.com/snivilised/nefilim"
	lab "github.com/snivilised/nefilim/internal/laboratory"
	"github.com/snivilised/nefilim/test/luna"
)

var _ = Describe("op: write-file", Ordered, func() {
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

			Context("op: WriteFile", func() {
				When("given: file does not already exist", func() {
					It("🧪 should: write successfully", func() {
						Expect(require(root, lab.Static.FS.Scratch)).To(Succeed())
						name := lab.Static.FS.Write.Destination
						Expect(fS.WriteFile(
							name, lab.Static.FS.Write.Content, lab.Perms.File.Perm(),
						)).To(Succeed())
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

			Context("op: WriteFile", func() {
				When("given: file does not already exist", func() {
					It("🧪 should: write successfully", func() {
						file := lab.Static.FS.Write.Destination
						Expect(require(
							root, lab.Static.FS.Scratch, file,
						)).To(Succeed())
						Expect(fS.WriteFile(
							file, lab.Static.FS.Write.Content, lab.Perms.File.Perm(),
						)).To(Succeed())
					})
				})
			})
		})
	})
})
