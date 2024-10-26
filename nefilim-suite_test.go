package nef_test

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	. "github.com/onsi/ginkgo/v2" //nolint:revive // ok
	. "github.com/onsi/gomega"    //nolint:revive // ok
	nef "github.com/snivilised/nefilim"
	lab "github.com/snivilised/nefilim/internal/laboratory"
	"github.com/snivilised/nefilim/test/luna"
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

func requires(fS nef.WriterFS, root, parent string, files ...string) error {
	if err := fS.MakeDirAll(filepath.Join(root, parent), lab.Perms.Dir.Perm()); err != nil {
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

func errorHomeResolver() (string, error) {
	return "", errors.New("failed to resolve home")
}

func errorAbsResolver(_ string) (string, error) {
	return "", errors.New("failed to resolve abs")
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
	return luna.Repo("Test/test.log")
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
