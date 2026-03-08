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
		// Root is the root path of the file system
		Root      string
		// Overwrite is true if the file system should overwrite existing files
		Overwrite bool
	}

	// FSUtility provides the path calculator and relative-root flag used by the file system.
	FSUtility interface {
		// Calc is the path calculator used by the FS
		Calc() PathCalc

		// IsRelative determines if the methods invoked on the file
		// system should use paths that are relative to a root specified
		// when created.
		IsRelative() bool
	}

	// ExistsInFS contains methods that check the existence of file system items.
	ExistsInFS interface {
		FSUtility
		// FileExists does file exist at the path specified
		FileExists(name string) bool

		// DirectoryExists does directory exist at the path specified
		DirectoryExists(name string) bool
	}

	// ReadFileFS is a file system that supports reading a file's contents in one call (non-streaming).
	ReadFileFS interface {
		fs.FS
		// Read reads file at path, from file system specified
		ReadFile(name string) ([]byte, error)
	}

	// ReaderFS is a file system that can stat, read directories, check existence, and read file contents.
	ReaderFS interface {
		fs.StatFS
		fs.ReadDirFS
		ExistsInFS
		ReadFileFS
	}

	// PathAs used with Ensure to define how to ensure that a path exists
	// at the location specified
	PathAs struct {
		// Name is the path to ensure
		Name    string
		// Default is the default value to return if the path does not exist
		Default string
		// Perm is the permission to set on the path
		Perm    os.FileMode
		// AsFile is true if the path is to be interpreted as a file, false if
		// it is to be interpreted as a directory
		AsFile  bool
	}

	// MakeDirFS is a file system with a MkDirAll method.
	MakeDirFS interface {
		ExistsInFS
		// MakeDir creates a new directory with the specified name and permission
		// bits (before umask).
		// If there is an error, it will be of type *PathError.
		MakeDir(name string, perm os.FileMode) error
		// MakeDirAll creates a directory named path,
		// along with any necessary parents, and returns nil,
		// or else returns an error.
		// The permission bits perm (before umask) are used for all
		// directories that MkdirAll creates.
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

	// MoveFS is a file system that supports moving an item from one path to another.
	MoveFS interface {
		// Move moves an item from one path to another
		Move(from, to string) error
	}

	// MoverFS extends MoveFS with existence checks and stat; used for move operations.
	MoverFS interface {
		MoveFS
		ExistsInFS
		fs.StatFS
	}

	// ChangeFS is a file system that supports changing an item (e.g. overwrite in place)
	// from one path to another.
	ChangeFS interface {
		// Change changes an item from one path to another
		Change(from, to string) error
	}

	// ChangerFS extends ChangeFS with existence checks and stat; used for change operations.
	ChangerFS interface {
		ChangeFS
		ExistsInFS
		fs.StatFS
	}

	// CopyFS is a file system that supports copying within the FS and copying another
	// fs.FS into a directory.
	CopyFS interface {
		// Copy copies an item from one path to another
		Copy(from, to string) error
		// CopyFS copies the file system fsys into the directory dir,
		// creating dir if necessary.
		CopyFS(dir string, fsys fs.FS) error
	}

	// RemoveFS is a file system that supports removing a single item or a directory tree.
	RemoveFS interface {
		// Remove removes the named file or (empty) directory.
		// If there is an error, it will be of type *PathError.
		Remove(name string) error
		// RemoveAll removes path and any children it contains.
		// It removes everything it can but returns the first error
		// it encounters. If the path does not exist, RemoveAll
		// returns nil (no error).
		// If there is an error, it will be of type *PathError.
		RemoveAll(path string) error
	}

	// RenameFS is a file system that supports renaming an item from one path to another.
	RenameFS interface {
		Rename(from, to string) error
	}

	// WriteFileFS is a file system that supports creating and writing files in one
	// call (non-streaming).
	WriteFileFS interface {
		FSUtility
		// Create creates or truncates the named file.
		Create(name string) (fs.File, error)
		// Write writes file at path, to file system specified
		WriteFile(name string, data []byte, perm os.FileMode) error
	}

	// WriterFS is a file system that supports change, copy, make dir, move, remove,
	// rename, and write.
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

	// UniversalFS is a file system that provides both read and write
	// capabilities (ReaderFS and WriterFS).
	UniversalFS interface {
		ReaderFS
		WriterFS
	}
)
