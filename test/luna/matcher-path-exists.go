package luna

import (
	"fmt"

	"github.com/onsi/gomega/types"
	nef "github.com/snivilised/nefilim"
)

type PathExistsMatcher struct {
	fS interface{}
}

type AsDirectory string
type AsFile string

func ExistInFS(fs interface{}) types.GomegaMatcher {
	return &PathExistsMatcher{
		fS: fs,
	}
}

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

func (m *PathExistsMatcher) FailureMessage(actual interface{}) string {
	return fmt.Sprintf("🔥 Expected\n\t%v\npath to exist", actual)
}

func (m *PathExistsMatcher) NegatedFailureMessage(actual interface{}) string {
	return fmt.Sprintf("🔥 Expected\n\t%v\npath NOT to exist\n", actual)
}
