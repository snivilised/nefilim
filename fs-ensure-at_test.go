package nef_test

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"testing/fstest"

	. "github.com/onsi/ginkgo/v2" //nolint:revive // ok
	. "github.com/onsi/gomega"    //nolint:revive // ok

	nef "github.com/snivilised/nefilim"
	lab "github.com/snivilised/nefilim/internal/laboratory"
)

var _ = Describe("Ensure", Ordered, func() {
	var (
		mocks *nef.ResolveMocks
		mapFS *makeDirMapFS
		root  string
		fS    nef.MakeDirFS
	)

	BeforeAll(func() {
		root = Repo("test")
	})

	BeforeEach(func() {
		mocks = &nef.ResolveMocks{
			HomeFunc: func() (string, error) {
				return filepath.Join(string(filepath.Separator), "home", "prodigy"), nil
			},
			AbsFunc: func(_ string) (string, error) { // no-op
				return "", errors.New("not required for these tests")
			},
		}

		mapFS = &makeDirMapFS{
			mapFS: fstest.MapFS{
				filepath.Join("home", "prodigy"): &fstest.MapFile{
					Mode: os.ModeDir,
				},
			},
		}
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

	DescribeTable("with mapFS",
		func(entry *ensureTE) {
			home, _ := mocks.HomeFunc()
			location := TrimRoot(filepath.Join(home, entry.relative))

			if entry.directory {
				location += string(filepath.Separator)
			}

			actual, err := nef.EnsurePathAt(location, lab.Static.FS.Ensure.Default.File, lab.Perms.File, mapFS)
			directory, _ := filepath.Split(actual)
			directory = filepath.Clean(directory)
			expected := TrimRoot(Path(home, entry.expected))

			Expect(err).Error().To(BeNil())
			Expect(actual).To(Equal(expected))
			Expect(AsDirectory(TrimRoot(directory))).To(ExistInFS(mapFS))
		},
		func(entry *ensureTE) string {
			return fmt.Sprintf("🧪 ===> given: '%v', should: '%v'", entry.given, entry.should)
		},

		XEntry(nil, &ensureTE{
			given:    "path is file",
			should:   "create parent directory and return specified file path",
			relative: filepath.Join("logs", "test.log"), // home/logs/test.log
			expected: "logs/test.log",
		}),

		XEntry(nil, &ensureTE{
			given:     "path is directory",
			should:    "create parent directory and return default file path",
			relative:  "logs/",
			directory: true,
			expected:  "logs/default-test.log",
		}),
	)
})
