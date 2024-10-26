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

const (
	home = "home/prodigy"
)

var _ = Describe("EnsurePathAt", Ordered, func() {
	var (
		mocks *nef.ResolveMocks
		fS    nef.UniversalFS
		root  string
	)

	BeforeAll(func() {
		root = luna.Repo("test")
	})

	BeforeEach(func() {
		scratch(root)

		mocks = &nef.ResolveMocks{
			HomeFunc: func() (string, error) {
				return filepath.Join(root, "scratch", home), nil
			},
			AbsFunc: func(_ string) (string, error) {
				return "", errors.New("not required for these tests")
			},
		}

		fS = nef.NewUniversalABS()
	})

	DescribeTable("with absolute fs",
		func(entry *ensureTE) {
			home, _ := mocks.HomeFunc()
			location := filepath.Join(home, entry.relative)

			if entry.directory {
				location += string(filepath.Separator)
			}

			actual, err := nef.EnsurePathAt(location, "default-test.log", lab.Perms.Dir, fS)
			directory, _ := filepath.Split(actual)
			directory = filepath.Clean(directory)
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
