package nef_test

import (
	"io/fs"
	"os"

	. "github.com/onsi/ginkgo/v2" //nolint:revive // ok
	. "github.com/onsi/gomega"    //nolint:revive // ok
	nef "github.com/snivilised/nefilim"
	lab "github.com/snivilised/nefilim/internal/laboratory"
)

// NB: these tests should NEVER be run in parallel because they interact with
// local filesystem.
var _ = Describe("file systems", Ordered, func() {
	var root string

	BeforeAll(func() {
		root = Repo("test")
	})

	Context("fs: StatFS", func() {
		var fS fs.StatFS

		BeforeEach(func() {
			fS = nef.NewStatFS(nef.At{
				Root: root,
			})
		})

		Context("op: FileExists", func() {
			When("given: existing path", func() {
				It("🧪 should: return true", func() {
					_, err := fS.Stat(lab.Static.FS.Existing.File)
					Expect(err).To(Succeed())
				})
			})

			When("given: path does not exist", func() {
				It("🧪 should: return false", func() {
					_, err := fS.Stat(lab.Static.Foo)
					Expect(err).To(MatchError(os.ErrNotExist))
				})
			})
		})
	})

	Context("fs: ExistsInFS", func() {
		var fS nef.ExistsInFS

		BeforeEach(func() {
			fS = nef.NewExistsInFS(nef.At{
				Root: root,
			})
		})

		Context("op: FileExists", func() {
			When("given: existing path", func() {
				It("🧪 should: return true", func() {
					Expect(AsFile(lab.Static.FS.Existing.File)).To(ExistInFS(fS))
				})
			})

			When("given: path does not exist", func() {
				It("🧪 should: return false", func() {
					Expect(AsFile(lab.Static.Foo)).NotTo(ExistInFS(fS))
				})
			})
		})

		Context("op: DirectoryExists", func() {
			When("given: existing path", func() {
				It("🧪 should: return true", func() {
					Expect(AsDirectory(lab.Static.FS.Existing.Directory)).To(ExistInFS(fS))
				})
			})

			When("given: path does not exist", func() {
				It("🧪 should: return false", func() {
					Expect(AsDirectory(lab.Static.Foo)).NotTo(ExistInFS(fS))
				})
			})
		})
	})

	Context("fs: ReadFileFS", func() {
		var fS nef.ReadFileFS

		BeforeEach(func() {
			fS = nef.NewReadFileFS(nef.At{
				Root: root,
			})
		})

		Context("op: ReadFile", func() {
			When("given: existing path", func() {
				It("🧪 should: ", func() {
					_, err := fS.ReadFile(lab.Static.FS.Existing.File)
					Expect(err).To(Succeed())
				})
			})

			When("given: path does not exist", func() {
				It("🧪 should: ", func() {
					_, err := fS.ReadFile(lab.Static.Foo)
					Expect(err).NotTo(Succeed())
				})
			})
		})
	})

	Context("fs: RenameFS", func() {
		Context("op: Rename", func() {
			When("given: ", func() {
				It("🧪 should: ", func() {

				})
			})
		})
	})
})
