package nef

import (
	"path/filepath"
)

// SplitParent returns the directory and base name of path, using
// filepath.Dir and filepath.Base.
func SplitParent(path string) (d, f string) {
	d = filepath.Dir(path)
	f = filepath.Base(path)

	return d, f
}
