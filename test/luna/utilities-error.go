package luna

import (
	"errors"
	"fmt"
	"os"

	"github.com/onsi/gomega"
)

// IsLinkError asserts that err is or wraps an os.LinkError, failing the test with reason if not.
func IsLinkError(err error, reason string) {
	var linkErr *os.LinkError
	gomega.Expect(errors.As(err, &linkErr)).To(gomega.BeTrue(),
		fmt.Sprintf("not LinkError, %q", reason),
	)
}
