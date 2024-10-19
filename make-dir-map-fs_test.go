package nef_test

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing/fstest"

	nef "github.com/snivilised/nefilim"
)

type (
	makeDirMapFS struct {
		mapFS fstest.MapFS
	}
)

func (f *makeDirMapFS) FileExists(path string) bool {
	fi, err := f.mapFS.Stat(path)
	if err != nil {
		return false
	}

	if fi.IsDir() {
		return false
	}

	return true
}

func (f *makeDirMapFS) DirectoryExists(path string) bool {
	if strings.HasPrefix(path, string(filepath.Separator)) {
		path = path[1:]
	}

	fileInfo, err := f.mapFS.Stat(path)
	if err != nil {
		return false
	}

	if !fileInfo.IsDir() {
		return false
	}

	return true
}

func (f *makeDirMapFS) MakeDir(path string, perm os.FileMode) error {
	if exists := f.DirectoryExists(path); !exists {
		f.mapFS[path] = &fstest.MapFile{
			Mode: fs.ModeDir | perm,
		}
	}

	return nil
}

func (f *makeDirMapFS) MakeDirAll(path string, perm os.FileMode) error {
	var current string
	segments := filepath.SplitList(path)

	for _, part := range segments {
		if current == "" {
			current = part
		} else {
			current += string(filepath.Separator) + part
		}

		if exists := f.DirectoryExists(current); !exists {
			f.mapFS[current] = &fstest.MapFile{
				Mode: fs.ModeDir | perm,
			}
		}
	}

	return nil
}

func (f *makeDirMapFS) Ensure(as nef.PathAs,
) (at string, err error) {
	_ = as
	panic("NOT-IMPL: makeDirMapFS.Ensure")
}
