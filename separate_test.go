package nef_test

import (
	"io/fs"
	"path/filepath"

	. "github.com/onsi/ginkgo/v2" //nolint:revive // ok
	. "github.com/onsi/gomega"    //nolint:revive // ok

	nef "github.com/snivilised/nefilim"
	"github.com/snivilised/nefilim/internal/third/lo"
	"github.com/snivilised/nefilim/test/luna"
)

var _ = Describe("Separate", Ordered, func() {
	var (
		fS             nef.UniversalFS
		calc           nef.PathCalc
		root, separate string

		expectations = struct {
			foo,
			bar,
			baz,
			x,
			y string
		}{
			foo: "foo.txt",
			bar: "bar.txt",
			baz: "baz.txt",
			x:   "x",
			y:   "y",
		}
	)

	BeforeAll(func() {
		root = luna.Repo("test")
		separate = filepath.Join("scratch", "separate")
	})

	BeforeEach(func() {
		scratch(root)
		fS = nef.NewUniversalABS()
		calc = fS.Calc()
	})

	When("directory contains mixed entries", func() {
		It("🧪 should: separate files from directories", func() {

			Expect(requires(fS, root, separate,
				calc.Join(separate, expectations.foo),
				calc.Join(separate, expectations.bar),
				calc.Join(separate, expectations.baz),
			)).To(Succeed())
			Expect(requires(fS, root, filepath.Join(separate, "x"))).To(Succeed())
			Expect(requires(fS, root, filepath.Join(separate, "y"))).To(Succeed())

			full := filepath.Join(root, separate)
			entries, err := fS.ReadDir(full)
			Expect(err).To(Succeed())
			files, directories := nef.Separate(entries)

			fileNames := lo.Map(files, func(entry fs.DirEntry, _ int) string {
				return entry.Name()
			})
			Expect(fileNames).To(ContainElements(
				expectations.foo, expectations.bar, expectations.baz,
			))

			dirNames := lo.Map(directories, func(entry fs.DirEntry, _ int) string {
				return entry.Name()
			})
			Expect(dirNames).To(ContainElements(
				expectations.x, expectations.y,
			))
		})
	})

	When("directory contains only file entries", func() {
		It("🧪 should: return files", func() {
			Expect(requires(fS, root, separate,
				calc.Join(separate, expectations.foo),
				calc.Join(separate, expectations.bar),
				calc.Join(separate, expectations.baz),
			)).To(Succeed())

			full := filepath.Join(root, separate)
			entries, err := fS.ReadDir(full)
			Expect(err).To(Succeed())
			files, directories := nef.Separate(entries)

			fileNames := lo.Map(files, func(entry fs.DirEntry, _ int) string {
				return entry.Name()
			})
			Expect(fileNames).To(ContainElements(
				expectations.foo, expectations.bar, expectations.baz,
			))
			Expect(directories).NotTo(BeNil())
			Expect(directories).To(BeEmpty())
		})
	})

	When("directory contains only directory entries", func() {
		It("🧪 should: return files", func() {
			Expect(requires(fS, root, separate)).To(Succeed())
			Expect(requires(fS, root, filepath.Join(separate, "x"))).To(Succeed())
			Expect(requires(fS, root, filepath.Join(separate, "y"))).To(Succeed())

			full := filepath.Join(root, separate)
			entries, err := fS.ReadDir(full)
			Expect(err).To(Succeed())
			files, directories := nef.Separate(entries)

			Expect(files).NotTo(BeNil())
			Expect(files).To(BeEmpty())

			dirNames := lo.Map(directories, func(entry fs.DirEntry, _ int) string {
				return entry.Name()
			})
			Expect(dirNames).To(ContainElements(
				expectations.x, expectations.y,
			))
		})
	})
})
