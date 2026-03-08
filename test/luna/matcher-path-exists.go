package luna

import (
	"fmt"

	"github.com/onsi/gomega/types"
	nef "github.com/snivilised/nefilim"
)

// PathExistsMatcher is a Gomega matcher that checks whether a path exists as a file or directory in an nef.ExistsInFS.
type PathExistsMatcher struct {
	fS interface{}
}

// AsDirectory is a path to be checked as a directory when used with ExistInFS.
type AsDirectory string

// AsFile is a path to be checked as a file when used with ExistInFS.
type AsFile string

// ExistInFS returns a Gomega matcher that asserts a path exists in the given file system (use AsDirectory or AsFile for actual).
func ExistInFS(fs interface{}) types.GomegaMatcher {
	return &PathExistsMatcher{
		fS: fs,
	}
}

// Match runs the matcher; actual must be AsDirectory or AsFile.
func (m *PathExistsMatcher) Match(actual interface{}) (bool, error) {
	FS, fileSystemOK := m.fS.(nef.ExistsInFS)
	if !fileSystemOK {
		return false, fmt.Errorf("❌ matcher expected an nef.ExistsInFS instance (%T)", FS)
	}

	if actualPath, dirOK := actual.(AsDirectory); dirOK {
		return FS.DirectoryExists(string(actualPath)), nil
	}

	if actualPath, fileOK := actual.(AsFile); fileOK {
		return FS.FileExists(string(actualPath)), nil
	}

	return false, fmt.Errorf("❌ matcher expected an AsDirectory or AsFile instance (%T)", actual)
}

// FailureMessage returns the message shown when the matcher fails (expected path to exist).
func (m *PathExistsMatcher) FailureMessage(actual interface{}) string {
	return fmt.Sprintf("🔥 Expected\n\t%v\npath to exist", actual)
}

// NegatedFailureMessage returns the message shown when the negated matcher fails (expected path not to exist).
func (m *PathExistsMatcher) NegatedFailureMessage(actual interface{}) string {
	return fmt.Sprintf("🔥 Expected\n\t%v\npath NOT to exist\n", actual)
}
