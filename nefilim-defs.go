package nef

import (
	"io/fs"
	"os"
)

// 📦 pkg: nef - contains local file system abstractions for navigation.
// Since there are no standard write-able file system interfaces,
// we need to define proprietary ones here in this package.

type (
	// Rel represents generic info required to create a relative file system.
	// Relative just means that a file system is created with a root path and
	// the operations on the file system are invoked with paths that must be
	// relative to the root.
	Rel struct {
		Root      string
		Overwrite bool
	}

	// FileSystems contains the logical file systems required
	// for navigation.
	FileSystems struct {
		// T is the file system that contains just the functionality required
		// for traversal. It can also represent other file systems including afero,
		// providing the appropriate adapters are in place.
		T TraverseFS
	}

	// ExistsInFS contains methods that check the existence of file system items.
	ExistsInFS interface {
		// FileExists does file exist at the path specified
		FileExists(name string) bool

		// DirectoryExists does directory exist at the path specified
		DirectoryExists(name string) bool
	}

	// ReadFileFS file system non streaming reader
	ReadFileFS interface {
		fs.FS
		// Read reads file at path, from file system specified
		ReadFile(name string) ([]byte, error)
	}

	// ReaderFS
	ReaderFS interface {
		fs.StatFS
		fs.ReadDirFS
		ExistsInFS
		ReadFileFS
	}

	// PathAs used with Ensure to define how to ensure that a path exists
	// at the location specified
	PathAs struct {
		Name    string
		Default string
		Perm    os.FileMode
		AsFile  bool
	}

	// MakeDirFS is a file system with a MkDirAll method.
	MakeDirFS interface {
		ExistsInFS
		MakeDir(name string, perm os.FileMode) error
		MakeDirAll(name string, perm os.FileMode) error
		// Ensure makes sure that a path exists (PathAs.Name). If the path exists
		// as a file then no directories need to be created and this file name
		// (PathAs.Name) is returned. If the path exists as a directory, then again
		// no directories are created, but the default (PathAs.Default) is returned.
		//
		// If the path does not exist, then 1 of 2 things can happen. If PathAs.AsFile
		// is set to true, then the parent of the path is created, and file portion
		// of the path is returned. When PathAs.AsFile is not set, ie the path
		// provided is to be interpreted as a directory, then this directory is
		// created and the default is returned.
		Ensure(as PathAs) (string, error)
	}

	// MoveFS
	MoveFS interface {
		Move(from, to string) error
	}

	// MoverFS
	MoverFS interface {
		MoveFS
		ExistsInFS
		fs.StatFS
	}

	ChangeFS interface {
		Change(from, to string) error
	}

	ChangerFS interface {
		ChangeFS
		ExistsInFS
		fs.StatFS
	}

	// CopyFS
	CopyFS interface {
		Copy(from, to string) error
		// CopyFS copies the file system fsys into the directory dir,
		// creating dir if necessary.
		CopyFS(dir string, fsys fs.FS) error
	}

	// RemoveFS
	RemoveFS interface {
		Remove(name string) error
		RemoveAll(path string) error
	}

	// RenameFS
	RenameFS interface {
		Rename(from, to string) error
	}

	// WriteFileFS file system non streaming writer
	WriteFileFS interface {
		// Create creates or truncates the named file.
		Create(name string) (*os.File, error)
		// Write writes file at path, to file system specified
		WriteFile(name string, data []byte, perm os.FileMode) error
	}

	// WriterFS
	WriterFS interface {
		ChangeFS
		CopyFS
		ExistsInFS
		MakeDirFS
		MoveFS
		RemoveFS
		RenameFS
		WriteFileFS
	}

	// TraverseFS non streaming file system with reader and some
	// writer capabilities
	TraverseFS interface {
		MakeDirFS
		ReaderFS
		WriteFileFS
	}

	// UniversalFS the file system that can do it all
	UniversalFS interface {
		ReaderFS
		WriterFS
	}
)
