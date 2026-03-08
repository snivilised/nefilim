package luna

import (
	"io"
	"io/fs"
	"os"
	"strings"
	"testing/fstest"
	"time"

	nef "github.com/snivilised/nefilim"
	lab "github.com/snivilised/nefilim/internal/laboratory"
	"github.com/snivilised/nefilim/internal/third/lo"
)

// MemFS is a memory fs based on fstest.MapFS intended to be used in
// unit tests. Clients can embed and override the methods defined here
// without having to provide a full implementation from scratch.
type MemFS struct {
	fstest.MapFS
	calc nef.PathCalc
}

var (
	_ nef.UniversalFS = (*MemFS)(nil)
)

// NewMemFS returns a new in-memory file system implementing nef.UniversalFS for tests.
func NewMemFS() *MemFS {
	return &MemFS{
		MapFS: fstest.MapFS{},
		calc:  &nef.RelativeCalc{},
	}
}

// Calc returns the path calculator used by this file system.
func (f *MemFS) Calc() nef.PathCalc {
	return f.calc
}

// IsRelative reports whether paths are relative to a root (always true for MemFS).
func (f *MemFS) IsRelative() bool {
	return true
}

// FileExists reports whether a regular file exists at name.
func (f *MemFS) FileExists(name string) bool {
	if mapFile, found := f.MapFS[name]; found && !mapFile.Mode.IsDir() {
		return true
	}

	return false
}

// DirectoryExists reports whether a directory exists at name.
func (f *MemFS) DirectoryExists(name string) bool {
	if mapFile, found := f.MapFS[name]; found && mapFile.Mode.IsDir() {
		return true
	}

	return false
}

// Create creates or truncates the named file; returns fs.ErrExist if it already exists.
func (f *MemFS) Create(name string) (fs.File, error) {
	if _, err := f.Stat(name); err == nil {
		return nil, fs.ErrExist
	}

	adapter := &FileAdapter{name: name}
	f.MapFS[name] = &fstest.MapFile{Data: adapter.data}

	return adapter, nil
}

// MakeDir creates a single directory at name with the given permissions.
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

// MakeDirAll creates the directory and any parents as needed.
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

// CopyFS copies the given fs.FS into the directory; not implemented on MemFS (returns nil).
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

// RemoveAll removes path and any children; returns os.ErrNotExist if path does not exist.
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

// Rename renames the item at from to to; returns os.ErrNotExist if from does not exist.
func (f *MemFS) Rename(from, to string) error {
	if item, found := f.MapFS[from]; found {
		delete(f.MapFS, from)
		f.MapFS[to] = item

		return nil
	}

	return os.ErrNotExist
}

// WriteFile writes data to the named file, creating it if necessary; returns fs.ErrExist if it already exists.
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

// FileInfoAdapter adapts in-memory file metadata to os.FileInfo for MemFS.
type FileInfoAdapter struct {
	name string
	size int64
	dir  bool
}

// Name returns the base name of the file or directory.
func (fi *FileInfoAdapter) Name() string {
	return fi.name
}

// Size returns the length in bytes for files; meaning is implementation-dependent for directories.
func (fi *FileInfoAdapter) Size() int64 {
	return fi.size
}

// Mode returns the file mode bits (uses lab.Perms.File for MemFS).
func (fi *FileInfoAdapter) Mode() os.FileMode {
	return lab.Perms.File
}

// ModTime returns the modification time (zero value for MemFS).
func (fi *FileInfoAdapter) ModTime() time.Time {
	var t time.Time
	return t
}

// IsDir returns true for directories, false for files.
func (fi *FileInfoAdapter) IsDir() bool {
	return fi.dir
}

// Sys returns underlying data source (nil for MemFS).
func (fi *FileInfoAdapter) Sys() interface{} {
	return nil
}

// FileAdapter is an in-memory fs.File used by MemFS for read/write.
type FileAdapter struct {
	name string
	data []byte
	pos  int64
}

// Read reads up to len(p) bytes from the file into p.
func (f *FileAdapter) Read(p []byte) (n int, err error) {
	if f.pos >= int64(len(f.data)) {
		return 0, io.EOF
	}

	n = copy(p, f.data[f.pos:])
	f.pos += int64(n)

	return n, nil
}

// Write appends p to the file content and advances the position.
func (f *FileAdapter) Write(p []byte) (n int, err error) {
	f.data = append(f.data[:f.pos], p...)
	n = len(p)
	f.pos += int64(n)

	return n, nil
}

// Close closes the file (no-op for MemFS).
func (f *FileAdapter) Close() error {
	return nil
}

// Stat returns FileInfo for the open file.
func (f *FileAdapter) Stat() (os.FileInfo, error) {
	return &FileInfoAdapter{
		name: f.name,
		size: int64(len(f.data)),
	}, nil
}
