package nef

import (
	"io/fs"
)

// syntheticEntry is a synthetic DirEntry implementation that wraps a FileInfo.
type syntheticEntry struct {
	fi fs.FileInfo
}

// Name returns the name of the entry.
func (e *syntheticEntry) Name() string {
	return e.fi.Name()
}

// IsDir returns true if the entry is a directory.
func (e *syntheticEntry) IsDir() bool {
	return e.fi.IsDir()
}

// Type returns the type of the entry.
func (e *syntheticEntry) Type() fs.FileMode {
	return e.fi.Mode()
}

// Info returns the FileInfo for the entry.
func (e *syntheticEntry) Info() (fs.FileInfo, error) {
	return e.fi, nil
}

// FromFileInfo synthesises a DirEntry instance from a FileInfo
func FromFileInfo(fi fs.FileInfo) fs.DirEntry {
	return &syntheticEntry{
		fi: fi,
	}
}
