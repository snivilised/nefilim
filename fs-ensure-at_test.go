package nef_test

import (
	"fmt"
	"path/filepath"

	. "github.com/onsi/ginkgo/v2" //nolint:revive // ok
	. "github.com/onsi/gomega"    //nolint:revive // ok

	nef "github.com/snivilised/nefilim"
	lab "github.com/snivilised/nefilim/internal/laboratory"
)

var _ = Describe("Ensure", Ordered, func() {
	var (
		root string
		fS   nef.MakeDirFS
	)

	BeforeAll(func() {
		root = Repo("test")
	})

	BeforeEach(func() {
		scratch(root)

		fS = nef.NewMakeDirFS(nef.At{
			Root: root,
		})
	})

	DescribeTable("local-fs",
		func(entry fsTE[nef.MakeDirFS]) {
			if entry.arrange != nil {
				entry.arrange(entry, fS)
			}
			entry.action(entry, fS)
		},
		func(entry fsTE[nef.MakeDirFS]) string {
			return fmt.Sprintf("🧪 ===> given: '%v', %v should: '%v'",
				entry.given, entry.op, entry.should,
			)
		},

		Entry(nil, fsTE[nef.MakeDirFS]{
			given:   "path exists as file",
			should:  "return filename of path",
			op:      "Ensure",
			require: lab.Static.FS.Scratch,
			target:  lab.Static.FS.Ensure.Log.File,
			arrange: func(entry fsTE[nef.MakeDirFS], _ nef.MakeDirFS) {
				directory := lab.Static.FS.Ensure.Default.Directory
				Expect(require(root, directory, entry.target)).To(Succeed())
			},
			action: func(entry fsTE[nef.MakeDirFS], fS nef.MakeDirFS) {
				result, err := fS.Ensure(
					nef.PathAs{
						Name:    entry.target,
						Default: lab.Static.FS.Ensure.Default.File,
						Perm:    lab.Perms.Dir,
					},
				)
				Expect(err).To(Succeed())
				_, file := filepath.Split(lab.Static.FS.Ensure.Log.File)
				Expect(result).To(Equal(file))
			},
		}),

		Entry(nil, fsTE[nef.MakeDirFS]{
			given:   "path exists as directory",
			should:  "return default",
			op:      "Ensure",
			require: lab.Static.FS.Scratch,
			target:  lab.Static.FS.Ensure.Log.Directory,
			arrange: func(entry fsTE[nef.MakeDirFS], _ nef.MakeDirFS) {
				Expect(require(root, entry.target)).To(Succeed())
			},
			action: func(entry fsTE[nef.MakeDirFS], fS nef.MakeDirFS) {
				_, file := filepath.Split(lab.Static.FS.Ensure.Default.File)
				result, err := fS.Ensure(
					nef.PathAs{
						Name:    entry.target,
						Default: file,
						Perm:    lab.Perms.Dir,
					},
				)
				Expect(err).To(Succeed())
				Expect(result).To(Equal(file))
			},
		}),

		Entry(nil, fsTE[nef.MakeDirFS]{
			given:   "file does not exist",
			should:  "create parent directory and return file",
			op:      "Ensure",
			require: lab.Static.FS.Scratch,
			target:  lab.Static.FS.Ensure.Log.File,
			from:    lab.Static.FS.Ensure.Home,
			arrange: func(entry fsTE[nef.MakeDirFS], _ nef.MakeDirFS) {
				parent := Join(entry.require, entry.from)
				Expect(require(root, parent)).To(Succeed())
			},
			action: func(entry fsTE[nef.MakeDirFS], fS nef.MakeDirFS) {
				_, file := filepath.Split(lab.Static.FS.Ensure.Default.File)
				result, err := fS.Ensure(
					nef.PathAs{
						Name:    entry.target,
						Default: file,
						Perm:    lab.Perms.Dir,
						AsFile:  true,
					},
				)
				Expect(err).To(Succeed())
				ensureAt := lab.Static.FS.Ensure.Default.Directory
				Expect(AsDirectory(ensureAt)).To(ExistInFS(fS))
				_, file = filepath.Split(entry.target)
				Expect(result).To(Equal(file))
			},
		}),

		Entry(nil, fsTE[nef.MakeDirFS]{
			given:   "directory does not exist",
			should:  "create directory and return default",
			op:      "Ensure",
			require: lab.Static.FS.Scratch,
			target:  lab.Static.FS.Ensure.Log.Directory,
			from:    lab.Static.FS.Ensure.Home,
			arrange: func(entry fsTE[nef.MakeDirFS], _ nef.MakeDirFS) {
				parent := Join(entry.require, entry.from)
				Expect(require(root, parent)).To(Succeed())
			},
			action: func(entry fsTE[nef.MakeDirFS], fS nef.MakeDirFS) {
				_, file := filepath.Split(lab.Static.FS.Ensure.Default.File)
				result, err := fS.Ensure(
					nef.PathAs{
						Name:    entry.target,
						Default: file,
						Perm:    lab.Perms.Dir,
					},
				)
				Expect(err).To(Succeed())
				ensureAt := lab.Static.FS.Ensure.Default.Directory
				Expect(AsDirectory(ensureAt)).To(ExistInFS(fS))
				Expect(result).To(Equal(file))
			},
		}),
	)
})
