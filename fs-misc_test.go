package nef_test

import (
	"io/fs"
	"os"

	. "github.com/onsi/ginkgo/v2" //nolint:revive // ok
	. "github.com/onsi/gomega"    //nolint:revive // ok
	nef "github.com/snivilised/nefilim"
	lab "github.com/snivilised/nefilim/internal/laboratory"
	"github.com/snivilised/nefilim/test/luna"
)

// NB: these tests should NEVER be run in parallel because they interact with
// local filesystem.
var _ = Describe("file systems", Ordered, func() {
	var root string

	BeforeAll(func() {
		root = luna.Repo("test")
	})

	Context("fs: StatFS", func() {
		var fS fs.StatFS

		BeforeEach(func() {
			fS = nef.NewStatFS(nef.Rel{
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
			fS = nef.NewExistsInFS(nef.Rel{
				Root: root,
			})
		})

		Context("op: FileExists", func() {
			When("given: existing path", func() {
				It("🧪 should: return true", func() {
					Expect(luna.AsFile(lab.Static.FS.Existing.File)).To(luna.ExistInFS(fS))
				})
			})

			When("given: path does not exist", func() {
				It("🧪 should: return false", func() {
					Expect(luna.AsFile(lab.Static.Foo)).NotTo(luna.ExistInFS(fS))
				})
			})
		})

		Context("op: DirectoryExists", func() {
			When("given: existing path", func() {
				It("🧪 should: return true", func() {
					Expect(luna.AsDirectory(lab.Static.FS.Existing.Directory)).To(luna.ExistInFS(fS))
				})
			})

			When("given: path does not exist", func() {
				It("🧪 should: return false", func() {
					Expect(luna.AsDirectory(lab.Static.Foo)).NotTo(luna.ExistInFS(fS))
				})
			})
		})
	})

	Context("fs: ReadFileFS", func() {
		var fS nef.ReadFileFS

		BeforeEach(func() {
			fS = nef.NewReadFileFS(nef.Rel{
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

	Context("MakeDirFS", func() {
		Context("New", func() {
			It("🧪 should: create", func() {
				Expect(nef.NewMakeDirFS(nef.Rel{
					Root: root,
				})).NotTo(BeNil())
			})
		})
	})

	Context("ReaderFS", func() {
		Context("New", func() {
			It("🧪 should: create", func() {
				Expect(nef.NewReaderFS(nef.Rel{
					Root: root,
				})).NotTo(BeNil())
			})
		})
	})

	Context("WriterFS", func() {
		Context("New", func() {
			It("🧪 should: create", func() {
				Expect(nef.NewWriterFS(nef.Rel{
					Root: root,
				})).NotTo(BeNil())
			})
		})
	})

	Context("ReadDirFS", func() {
		Context("New", func() {
			It("🧪 should: create", func() {
				Expect(nef.NewReadDirFS(nef.Rel{
					Root: root,
				})).NotTo(BeNil())
			})
		})
	})

	Context("TraverseFS", func() {
		Context("New", func() {
			It("🧪 should: create", func() {
				Expect(nef.NewTraverseFS(nef.Rel{
					Root: root,
				})).NotTo(BeNil())
			})
		})
	})
})
