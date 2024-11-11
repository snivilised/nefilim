package nef_test

import (
	"fmt"

	. "github.com/onsi/ginkgo/v2" //nolint:revive // ok
	. "github.com/onsi/gomega"    //nolint:revive // ok
	nef "github.com/snivilised/nefilim"
	lab "github.com/snivilised/nefilim/internal/laboratory"
	"github.com/snivilised/nefilim/test/luna"
)

var _ = Describe("op: make-dir/all", Ordered, func() {
	var root string

	BeforeAll(func() {
		root = luna.Repo("test")
	})

	Context("tentative", func() {
		Context("fs: MakeDirFS", func() {
			var (
				fS nef.MakeDirFS
			)

			BeforeEach(func() {
				fS = nef.NewMakeDirFS(nef.Rel{
					Root:      root,
					Overwrite: false,
				})
				Expect(fS.IsRelative()).To(BeTrue())
				scratch(root)
			})

			Context("op: MakeDir", func() {
				When("given: path does not exist", func() {
					It("🧪 should: complete ok", func() {
						path := lab.Static.FS.Scratch
						Expect(fS.MakeDir(path, lab.Perms.Dir.Perm())).To(
							Succeed(), fmt.Sprintf("failed to MakeDir %q", path),
						)

						Expect(luna.AsDirectory(path)).To(luna.ExistInFS(fS))
					})
				})

				When("given: path already exists", func() {
					It("🧪 should: complete ok", func() {
						path := lab.Static.FS.Existing.Directory
						Expect(fS.MakeDir(path, lab.Perms.Dir.Perm())).To(
							Succeed(), fmt.Sprintf("failed to MakeDir %q", path),
						)
					})
				})
			})

			Context("op: MakeDirAll", func() {
				When("given: path does not exist", func() {
					It("🧪 should: complete ok", func() {
						path := lab.Static.FS.MakeDir.MakeAll
						Expect(fS.MakeDirAll(path, lab.Perms.Dir.Perm())).To(
							Succeed(), fmt.Sprintf("failed to MakeDir %q", path),
						)

						Expect(luna.AsDirectory(path)).To(luna.ExistInFS(fS))
					})
				})

				When("given: path already exists", func() {
					It("🧪 should: complete ok", func() {
						path := lab.Static.FS.Existing.Directory
						Expect(fS.MakeDir(path, lab.Perms.Dir.Perm())).To(
							Succeed(), fmt.Sprintf("failed to MakeDir %q", path),
						)

						Expect(luna.AsDirectory(path)).To(luna.ExistInFS(fS))
					})
				})
			})
		})
	})
})
