package nef_test

import (
	"errors"
	"fmt"
	"path/filepath"

	. "github.com/onsi/ginkgo/v2" //nolint:revive // ok
	. "github.com/onsi/gomega"    //nolint:revive // ok

	nef "github.com/snivilised/nefilim"
	lab "github.com/snivilised/nefilim/internal/laboratory"
	"github.com/snivilised/nefilim/test/luna"
)

var _ = Describe("Ensure", Ordered, func() {
	var (
		root string
		fS   nef.MakeDirFS
	)

	BeforeAll(func() {
		root = luna.Repo("test")
	})

	BeforeEach(func() {
		scratch(root)

		fS = nef.NewMakeDirFS(nef.Rel{
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
				_, file := fS.Calc().Split(lab.Static.FS.Ensure.Log.File)
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
				_, file := fS.Calc().Split(lab.Static.FS.Ensure.Default.File)
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
			arrange: func(entry fsTE[nef.MakeDirFS], fS nef.MakeDirFS) {
				parent := fS.Calc().Join(entry.require, entry.from)
				Expect(require(root, parent)).To(Succeed())
			},
			action: func(entry fsTE[nef.MakeDirFS], fS nef.MakeDirFS) {
				_, file := fS.Calc().Split(lab.Static.FS.Ensure.Default.File)
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
				Expect(luna.AsDirectory(ensureAt)).To(luna.ExistInFS(fS))
				_, file = fS.Calc().Split(entry.target)
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
				parent := fS.Calc().Join(entry.require, entry.from)
				Expect(require(root, parent)).To(Succeed())
			},
			action: func(entry fsTE[nef.MakeDirFS], fS nef.MakeDirFS) {
				_, file := fS.Calc().Split(lab.Static.FS.Ensure.Default.File)
				result, err := fS.Ensure(
					nef.PathAs{
						Name:    entry.target,
						Default: file,
						Perm:    lab.Perms.Dir,
					},
				)
				Expect(err).To(Succeed())
				ensureAt := lab.Static.FS.Ensure.Default.Directory
				Expect(luna.AsDirectory(ensureAt)).To(luna.ExistInFS(fS))
				Expect(result).To(Equal(file))
			},
		}),
	)
})

var _ = Describe("Ensure", Ordered, func() {
	const (
		home = "home/prodigy"
	)

	var (
		mocks *nef.ResolveMocks
		fS    nef.UniversalFS
		calc  nef.PathCalc
		root  string
	)

	BeforeAll(func() {
		root = luna.Repo("test")
	})

	BeforeEach(func() {
		fS = nef.NewUniversalABS()
		calc = fS.Calc()

		scratch(root)

		mocks = &nef.ResolveMocks{
			HomeFunc: func() (string, error) {
				return calc.Join(root, "scratch", home), nil
			},
			AbsFunc: func(_ string) (string, error) {
				return "", errors.New("not required for these tests")
			},
		}
	})

	DescribeTable("with absolute fs",
		func(entry *ensureTE) {
			home, _ := mocks.HomeFunc()
			location := calc.Join(home, entry.relative)

			if entry.directory {
				location += string(filepath.Separator)
			}

			actual, err := fS.Ensure(nef.PathAs{
				Name:    location,
				Default: "default-test.log",
				Perm:    lab.Perms.Dir,
			})

			directory, _ := calc.Split(actual)
			directory = calc.Clean(directory)
			expected := luna.Combine(home, entry.expected)

			Expect(err).Error().To(BeNil())
			Expect(actual).To(Equal(expected))
			Expect(luna.AsDirectory(directory)).To(luna.ExistInFS(fS))
		},
		func(entry *ensureTE) string {
			return fmt.Sprintf("🧪 ===> given: '%v', should: '%v'", entry.given, entry.should)
		},

		Entry(nil, &ensureTE{
			given:    "path is file",
			should:   "create parent directory and return specified file path",
			relative: filepath.Join("logs", "test.log"), // (can't use calc here, not set yet)
			expected: "logs/test.log",
		}),

		Entry(nil, &ensureTE{
			given:     "path is directory",
			should:    "create parent directory and return default file path",
			relative:  "logs/",
			directory: true,
			expected:  "logs/default-test.log",
		}),
	)
})
