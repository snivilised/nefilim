package luna_test

import (
	"os"

	. "github.com/onsi/ginkgo/v2" //nolint:revive // ok
	. "github.com/onsi/gomega"    //nolint:revive // ok

	nef "github.com/snivilised/nefilim"
	lab "github.com/snivilised/nefilim/internal/laboratory"
	"github.com/snivilised/nefilim/luna"
)

var (
	data = []byte("some content")
)

var _ = Describe("MemFs", func() {
	var fS nef.UniversalFS

	BeforeEach(func() {
		fS = luna.NewMemFS()
	})

	Context("Remove", func() {
		When("given: item exists", func() {
			It("🧪 should: remove item", func() {
				path := "widgets/foo.txt"
				Expect(fS.WriteFile(path, data, lab.Perms.File)).To(Succeed())
				Expect(fS.Remove(path)).To(Succeed())
				Expect(fS.FileExists(path)).To(BeFalse())
			})
		})

		When("given: item does NOT exist", func() {
			It("🧪 should: return error", func() {
				path := "missing/foo.txt"
				Expect(fS.Remove(path)).To(MatchError(os.ErrNotExist),
					"expected error os.ErrNotExist",
				)
			})
		})
	})

	Context("RemoveAll", func() {
		When("given: some items exist", func() {
			It("🧪 should: remove only matching item", func() {
				foo := "a/b/foo.txt"
				Expect(fS.WriteFile(foo, data, lab.Perms.File)).To(Succeed())
				bar := "a/b/bar.txt"
				Expect(fS.WriteFile(bar, data, lab.Perms.File)).To(Succeed())
				baz := "bin/baz.txt"
				Expect(fS.WriteFile(baz, data, lab.Perms.File)).To(Succeed())
				path := "a/b"

				Expect(fS.RemoveAll(path)).To(Succeed(), "failed to remove all in path")
				Expect(fS.FileExists(foo)).To(BeFalse(), "foo should not exist")
				Expect(fS.FileExists(bar)).To(BeFalse(), "bar should not exist")
				Expect(fS.FileExists(baz)).To(BeTrue(), "baz should still exist")
			})
		})
	})
})
