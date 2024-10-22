package nef

import (
	"strings"
)

const (
	delimiter = "/"
)

// utils required for relative file systems, that typically replace those
// defined in filepath that are based upon filepath.Separator.

// Parent (equivalent to filepath.Dir) returns all but the last element of path,
// typically the path's directory.
func Parent(path string) string {
	if !strings.Contains(path, delimiter) {
		return ""
	}

	index := strings.LastIndex(path, delimiter)
	return path[:index]
}

func Join(segments ...string) string {
	// required so we can avoid the use of the file utilities defined
	// in filepath, which are not appropriate for a relative file systems
	// because they are based upon translating delimiters into platform
	// specific separators.
	return strings.Join(segments, "/")
}
