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

var _ = Describe("EnsurePathAt", Ordered, func() {
	var (
		mocks *nef.ResolveMocks
		fS    *makeDirMapFS
	)

	BeforeEach(func() {
		mocks = &nef.ResolveMocks{
			HomeFunc: func() (string, error) {
				return filepath.Join(string(filepath.Separator), "home", "prodigy"), nil
			},
			AbsFunc: func(_ string) (string, error) {
				return "", errors.New("not required for these tests")
			},
		}

		fS = &makeDirMapFS{
			mapFS: fstest.MapFS{
				filepath.Join("home", "prodigy"): &fstest.MapFile{
					Mode: os.ModeDir,
				},
			},
		}
	})

	DescribeTable("with mapFS",
		func(entry *ensureTE) {
			home, _ := mocks.HomeFunc()
			location := TrimRoot(filepath.Join(home, entry.relative))

			if entry.directory {
				location += string(filepath.Separator)
			}

			actual, err := nef.EnsurePathAt(location, "default-test.log", lab.Perms.File, fS)
			directory, _ := filepath.Split(actual)
			directory = filepath.Clean(directory)
			expected := TrimRoot(Path(home, entry.expected))

			Expect(err).Error().To(BeNil())
			Expect(actual).To(Equal(expected))
			Expect(AsDirectory(TrimRoot(directory))).To(ExistInFS(fS))
		},
		func(entry *ensureTE) string {
			return fmt.Sprintf("🧪 ===> given: '%v', should: '%v'", entry.given, entry.should)
		},

		Entry(nil, &ensureTE{
			given:    "path is file",
			should:   "create parent directory and return specified file path",
			relative: filepath.Join("logs", "test.log"),
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
