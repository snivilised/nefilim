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

// 📦 pkg: luna

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

func NewMemFS() *MemFS {
	return &MemFS{
		MapFS: fstest.MapFS{},
		calc:  &nef.RelativeCalc{},
	}
}

func (f *MemFS) Calc() nef.PathCalc {
	return f.calc
}

func (f *MemFS) IsRelative() bool {
	return true
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

func (f *MemFS) Create(name string) (fs.File, error) {
	if _, err := f.Stat(name); err == nil {
		return nil, fs.ErrExist
	}

	adapter := &FileAdapter{name: name}
	f.MapFS[name] = &fstest.MapFile{Data: adapter.data}

	return adapter, nil
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

// File

type FileInfoAdapter struct {
	name string
	size int64
	dir  bool
}

func (fi *FileInfoAdapter) Name() string {
	return fi.name
}

func (fi *FileInfoAdapter) Size() int64 {
	return fi.size
}

func (fi *FileInfoAdapter) Mode() os.FileMode {
	return lab.Perms.File
}

func (fi *FileInfoAdapter) ModTime() time.Time {
	var t time.Time
	return t
}

func (fi *FileInfoAdapter) IsDir() bool {
	return fi.dir
}

func (fi *FileInfoAdapter) Sys() interface{} {
	return nil
}

type FileAdapter struct {
	name string
	data []byte
	pos  int64
}

func (f *FileAdapter) Read(p []byte) (n int, err error) {
	if f.pos >= int64(len(f.data)) {
		return 0, io.EOF
	}

	n = copy(p, f.data[f.pos:])
	f.pos += int64(n)

	return n, nil
}

func (f *FileAdapter) Write(p []byte) (n int, err error) {
	f.data = append(f.data[:f.pos], p...)
	n = len(p)
	f.pos += int64(n)

	return n, nil
}

func (f *FileAdapter) Close() error {
	return nil
}

func (f *FileAdapter) Stat() (os.FileInfo, error) {
	return &FileInfoAdapter{
		name: f.name,
		size: int64(len(f.data)),
	}, nil
}
