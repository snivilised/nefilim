package luna

import (
	"io/fs"
	"os"
	"strings"
	"testing/fstest"

	nef "github.com/snivilised/nefilim"
	lab "github.com/snivilised/nefilim/internal/laboratory"
	"github.com/snivilised/nefilim/internal/third/lo"
)

// MemFS is a memory fs based on fstest.MapFS intended to be used in
// unit tests. Clients can embed and override the methods defined here
// without having to provide a full implementation from scratch.
type MemFS struct {
	fstest.MapFS
}

var (
	_ nef.UniversalFS = (*MemFS)(nil)
)

func NewMemFS() *MemFS {
	return &MemFS{
		MapFS: fstest.MapFS{},
	}
}

func (f *MemFS) FileExists(name string) bool {
	if mapFile, found := f.MapFS[name]; found && !mapFile.Mode.IsDir() {
		return true
	}

	return false
}

func (f *MemFS) DirectoryExists(name string) bool {
	if mapFile, found := f.MapFS[name]; found && mapFile.Mode.IsDir() {
		return true
	}

	return false
}

func (f *MemFS) Create(name string) (*os.File, error) {
	if _, err := f.Stat(name); err == nil {
		return nil, fs.ErrExist
	}

	file := &fstest.MapFile{
		Mode: lab.Perms.File,
	}

	f.MapFS[name] = file
	// TODO: this needs a resolution using a file interface
	// rather than using os.File which is a struct not an
	// interface
	dummy := &os.File{}

	return dummy, nil
}

func (f *MemFS) MakeDir(name string, perm os.FileMode) error {
	if !fs.ValidPath(name) {
		return nef.NewInvalidPathError("MakeDir", name)
	}

	if _, found := f.MapFS[name]; !found {
		f.MapFS[name] = &fstest.MapFile{
			Mode: perm | os.ModeDir,
		}
	}

	return nil
}

func (f *MemFS) MakeDirAll(name string, perm os.FileMode) error {
	if !fs.ValidPath(name) {
		return nef.NewInvalidPathError("MakeDirAll", name)
	}

	segments := strings.Split(name, "/")

	_ = lo.Reduce(segments,
		func(acc []string, s string, _ int) []string {
			acc = append(acc, s)
			path := strings.Join(acc, "/")

			if _, found := f.MapFS[path]; !found {
				f.MapFS[path] = &fstest.MapFile{
					Mode: perm | os.ModeDir,
				}
			}

			return acc
		}, []string{},
	)

	return nil
}

// Ensure is not currently implemented on MemFS
func (f *MemFS) Ensure(_ nef.PathAs) (string, error) {
	return "", nil
}

// Move is not currently implemented on MemFS
func (f *MemFS) Move(_, _ string) error {
	return nil
}

// Change is not currently implemented on MemFS
func (f *MemFS) Change(_, _ string) error {
	return nil
}

// Copy is not currently implemented on MemFS
func (f *MemFS) Copy(_, _ string) error {
	return nil
}

func (f *MemFS) CopyFS(_ string, _ fs.FS) error {
	return nil
}

// Remove removes the named file or (empty) directory.
// If there is an error, it will be of type *PathError.
func (f *MemFS) Remove(name string) error {
	if _, found := f.MapFS[name]; found {
		delete(f.MapFS, name)
		return nil
	}

	return os.ErrNotExist
}

func (f *MemFS) RemoveAll(path string) error {
	keys := lo.Keys(f.MapFS)
	matched := lo.Filter(keys, func(item string, _ int) bool {
		return strings.HasPrefix(item, path)
	})

	if len(matched) == 0 {
		return os.ErrNotExist
	}

	for _, item := range matched {
		delete(f.MapFS, item)
	}

	return nil
}

func (f *MemFS) Rename(from, to string) error {
	if item, found := f.MapFS[from]; found {
		delete(f.MapFS, from)
		f.MapFS[to] = item

		return nil
	}

	return os.ErrNotExist
}

func (f *MemFS) WriteFile(name string, data []byte, perm os.FileMode) error {
	if _, err := f.Stat(name); err == nil {
		return fs.ErrExist
	}

	f.MapFS[name] = &fstest.MapFile{
		Data: data,
		Mode: perm,
	}

	return nil
}
