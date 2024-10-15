package nef

import (
	"os"
	"path/filepath"
	"strings"
)

// EnsurePathAt ensures that the specified path exists (including any non
// existing intermediate directories). Given a path and a default filename,
// the specified path is created in the following manner:
// - If the path denotes a file (path does not end is a directory separator), then
// the parent folder is created if it doesn't exist on the file-system provided.
// - If the path denotes a directory, then that directory is created.
//
// The returned string represents the file, so if the path specified was a
// directory path, then the defaultFilename provided is joined to the path
// and returned, otherwise the original path is returned un-modified.
// Note: filepath.Join does not preserve a trailing separator, therefore to make sure
// a path is interpreted as a directory and not a file, then the separator has
// to be appended manually onto the end of the path.
// If vfs is not provided, then the path is ensured directly on the native file
// system.

func EnsurePathAt(path, defaultFilename string, perm os.FileMode,
	fS ...MakeDirFS,
) (at string, err error) {
	var (
		directory, file string
	)

	if strings.HasSuffix(path, string(os.PathSeparator)) {
		directory = path
		file = defaultFilename
	} else {
		directory, file = filepath.Split(path)
	}

	if len(fS) > 0 {
		if !fS[0].DirectoryExists(directory) {
			err = fS[0].MakeDirAll(directory, perm)
		}
	} else {
		err = os.MkdirAll(directory, perm)
	}

	return filepath.Clean(filepath.Join(directory, file)), err
}
