package nef_test

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	. "github.com/onsi/ginkgo/v2" //nolint:revive // ok
	. "github.com/onsi/gomega"    //nolint:revive // ok
	nef "github.com/snivilised/nefilim"
	lab "github.com/snivilised/nefilim/internal/laboratory"
)

func TestNefilim(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Nefilim Suite")
}

type (
	ensureTE struct {
		given     string
		should    string
		relative  string
		expected  string
		directory bool
	}

	RPEntry struct {
		given  string
		should string
		path   string
		expect string
	}

	funcFS[T any] func(entry fsTE[T], fS T)

	fsTE[T any] struct {
		given     string
		should    string
		note      string
		op        string
		overwrite bool
		directory bool
		require   string
		target    string
		from      string
		to        string
		arrange   funcFS[T]
		action    funcFS[T]
	}
)

var (
	fakeHome      = filepath.Join(string(filepath.Separator), "home", "rabbitweed")
	fakeAbsCwd    = filepath.Join(string(filepath.Separator), "home", "rabbitweed", "music", "xpander")
	fakeAbsParent = filepath.Join(string(filepath.Separator), "home", "rabbitweed", "music")
)

func scratch(root string) { // should we re-create scratch too, so the tests don't have to?
	scratchPath := filepath.Join(root, lab.Static.FS.Scratch)

	if _, err := os.Stat(scratchPath); err == nil {
		Expect(os.RemoveAll(scratchPath)).To(Succeed(),
			fmt.Sprintf("failed to delete existing directory %q", scratchPath),
		)
	}
}

// require ensures that a path exists. If files are also provided,
// it will create these files too. The files are relative to the root
// and should be prefixed by parent; that is to say, when a test needs
// scratch/foo.txt, parent = 'scratch' and file = 'scratch/foo.txt';
// ie te file still needs to be relative to root, not parent.
func require(root, parent string, files ...string) error {
	if err := os.MkdirAll(filepath.Join(root, parent), lab.Perms.Dir.Perm()); err != nil {
		return fmt.Errorf("failed to create directory: %q (%w)", parent, err)
	}

	for _, name := range files {
		handle, err := os.Create(filepath.Join(root, name))
		if err != nil {
			return fmt.Errorf("failed to create file: %q (%w)", name, err)
		}

		handle.Close()
	}

	return nil
}

func fakeHomeResolver() (string, error) {
	return fakeHome, nil
}

func fakeAbsResolver(path string) (string, error) {
	if strings.HasPrefix(path, "..") {
		return filepath.Join(fakeAbsParent, path[2:]), nil
	}

	if strings.HasPrefix(path, ".") {
		return filepath.Join(fakeAbsCwd, path[1:]), nil
	}

	return path, nil
}

func TrimRoot(root string) string {
	// omit leading '/', because test-fs stupidly doesn't like it,
	// so we have to jump through hoops
	if strings.HasPrefix(root, string(filepath.Separator)) {
		return root[1:]
	}

	pattern := `^[a-zA-Z]:[\\/]*`
	re := regexp.MustCompile(pattern)

	return re.ReplaceAllString(root, "")
}

// Path creates a path from the parent combined with the relative path. The relative
// path is a file system path so should only contain forward slashes, not the standard
// file path separator as denoted by filepath.Separator, typically used when interacting
// with the local file system. Do not use trailing "/".
func Path(parent, relative string) string {
	if relative == "" {
		return parent
	}

	return parent + "/" + relative
}

// Repo gets the path of the repo with relative joined on
func Repo(relative string) string {
	cmd := exec.Command("git", "rev-parse", "--show-toplevel")
	output, _ := cmd.Output()
	repo := strings.TrimSpace(string(output))

	return Path(repo, relative)
}

func Join(segments ...string) string {
	return strings.Join(segments, "/")
}

func Normalise(p string) string {
	return strings.ReplaceAll(p, "/", string(filepath.Separator))
}

func Because(name, because string) string {
	return fmt.Sprintf("❌ for item named: '%v', because: '%v'", name, because)
}

func Reason(name string) string {
	return fmt.Sprintf("❌ for item named: '%v'", name)
}

func Log() string {
	return Repo("Test/test.log")
}

func IsLinkError(err error, reason string) {
	var linkErr *os.LinkError
	Expect(errors.As(err, &linkErr)).To(BeTrue(), fmt.Sprintf("not LinkError, %q", reason))
}

func IsInvalidPathError(err error, reason string) {
	Expect(nef.IsInvalidPathError(err)).To(BeTrue(),
		fmt.Sprintf("not NewInvalidPathError, %q", reason),
	)
}

func IsSameDirectoryMoveRejectionError(err error, reason string) {
	Expect(nef.IsRejectSameDirMoveError(err)).To(BeTrue(),
		fmt.Sprintf("not SameDirectoryMoveRejectionError, %q", reason),
	)
}

func IsDifferentDirectoryChangeRejectionError(err error, reason string) {
	Expect(nef.IsRejectDifferentDirChangeError(err)).To(BeTrue(),
		fmt.Sprintf("not DifferentDirectoryChangeRejectionError, %q", reason),
	)
}
