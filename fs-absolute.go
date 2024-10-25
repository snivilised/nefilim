package nef

import (
	"io/fs"
	"os"
)

type absoluteFS struct{}

// NewUniversalABS creates an absolute universal file system
func NewUniversalABS() UniversalFS {
	return &absoluteFS{}
}

// NewTraverseABS creates an absolute traverse file system
func NewTraverseABS() TraverseFS {
	return &absoluteFS{}
}

// NewReaderABS creates an absolute reader file system
func NewReaderABS() ReaderFS {
	return &absoluteFS{}
}

// NewWriterABS creates an absolute writer file system
func NewWriterABS() WriterFS {
	return &absoluteFS{}
}

// FileExists does file exist at the path specified
func (f *absoluteFS) FileExists(name string) bool {
	info, err := f.Stat(name)
	if err != nil {
		return false
	}

	if info.IsDir() {
		return false
	}

	return true
}

// DirectoryExists does directory exist at the path specified
func (f *absoluteFS) DirectoryExists(name string) bool {
	info, err := f.Stat(name)
	if err != nil {
		return false
	}

	if !info.IsDir() {
		return false
	}

	return true
}

// Open
func (f *absoluteFS) Open(name string) (fs.File, error) {
	return os.Open(name)
}

// ReadFile reads the named file and returns the contents.
// A successful call returns err == nil, not err == EOF.
// Because ReadFile reads the whole file, it does not treat an EOF from Read
// as an error to be reported.
func (f *absoluteFS) ReadFile(name string) ([]byte, error) {
	return os.ReadFile(name)
}

// Stat returns a [FileInfo] describing the named file.
// If there is an error, it will be of type [*PathError].
func (f *absoluteFS) Stat(name string) (fs.FileInfo, error) {
	return os.Stat(name)
}

// ReadDir reads the named directory
// and returns a list of directory entries sorted by filename.
//
// If fs implements [ReadDirFS], ReadDir calls fs.ReadDir.
// Otherwise ReadDir calls fs.Open and uses ReadDir and Close
// on the returned file.
func (f *absoluteFS) ReadDir(name string) ([]fs.DirEntry, error) {
	return os.ReadDir(name)
}

// Mkdir creates a new directory with the specified name and permission
// bits (before umask).
// If there is an error, it will be of type *PathError.
func (f *absoluteFS) MakeDir(name string, perm os.FileMode) error {
	return os.Mkdir(name, perm)
}

// MkdirAll creates a directory named path,
// along with any necessary parents, and returns nil,
// or else returns an error.
// The permission bits perm (before umask) are used for all
// directories that MkdirAll creates.
// If path is already a directory, MkdirAll does nothing
// and returns nil.
func (f *absoluteFS) MakeDirAll(name string, perm os.FileMode) error {
	return os.MkdirAll(name, perm)
}

// Ensure is not currently implemented on absoluteFS
func (f *absoluteFS) Ensure(_ PathAs) (string, error) {
	panic("NOT-IMPL: absoluteFS.Ensure")
}

// Move is not currently implemented on absoluteFS
func (f *absoluteFS) Move(_, _ string) error {
	panic("NOT-IMPL: absoluteFS.Move")
}

// Change is not currently implemented on absoluteFS
func (f *absoluteFS) Change(_, _ string) error {
	panic("NOT-IMPL: absoluteFS.Change")
}

// Copy is not currently implemented on absoluteFS
func (f *absoluteFS) Copy(_, _ string) error {
	panic("NOT-IMPL: absoluteFS.Copy")
}

// CopyFS copies the file system fsys into the directory dir,
// creating dir if necessary.
//
// Files are created with mode 0o666 plus any execute permissions
// from the source, and directories are created with mode 0o777
// (before umask).
//
// CopyFS will not overwrite existing files, and returns an error
// if a file name in fsys already exists in the destination.
//
// Symbolic links in fsys are not supported. A *PathError with Err set
// to ErrInvalid is returned when copying from a symbolic link.
//
// Symbolic links in dir are followed.
//
// Copying stops at and returns the first error encountered.
func (f *absoluteFS) CopyFS(dir string, fsys fs.FS) error {
	return os.CopyFS(dir, fsys)
}

// Remove removes the named file or (empty) directory.
// If there is an error, it will be of type *PathError.
func (f *absoluteFS) Remove(name string) error {
	return os.Remove(name)
}

// RemoveAll removes path and any children it contains.
// It removes everything it can but returns the first error
// it encounters. If the path does not exist, RemoveAll
// returns nil (no error).
// If there is an error, it will be of type [*PathError].
func (f *absoluteFS) RemoveAll(path string) error {
	return os.RemoveAll(path)
}

// Rename renames (moves) 'from' to 'to'.
// If 'to' already exists and is not a directory, Rename replaces it.
// OS-specific restrictions may apply when 'from' and 'to' are in different directories.
// Even within the same directory, on non-Unix platforms Rename is not an atomic operation.
// If there is an error, it will be of type *LinkError.
func (f *absoluteFS) Rename(from, to string) error {
	return os.Rename(from, to)
}

// Create creates or truncates the named file. If the file already exists,
// it is truncated. If the file does not exist, it is created with mode 0o666
// (before umask). If successful, methods on the returned File can
// be used for I/O; the associated file descriptor has mode O_RDWR.
// If there is an error, it will be of type *PathError.
func (f *absoluteFS) Create(name string) (fs.File, error) {
	return os.Create(name)
}

// WriteFile writes data to the named file, creating it if necessary.
// If the file does not exist, WriteFile creates it with permissions perm (before umask);
// otherwise WriteFile truncates it before writing, without changing permissions.
// Since WriteFile requires multiple system calls to complete, a failure mid-operation
// can leave the file in a partially written state.
func (f *absoluteFS) WriteFile(name string, data []byte, perm os.FileMode) error {
	return os.WriteFile(name, data, perm)
}
