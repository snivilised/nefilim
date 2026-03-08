package lab

import (
	"io/fs"
)

type (
	// Change holds paths for a change operation (from, destination, to).
	Change struct {
		// From is the path of the item to change
		From        Pair
		// Destination is the path to change the item to
		Destination string
		// To is the path to change the item to
		To          Pair
	}

	// Copy holds the destination path for a copy operation.
	Copy struct {
		// Destination is the path to copy the item to
		Destination string
	}
	// Create holds the destination path for a create operation.
	Create struct {
		// Destination is the path to create the item at
		Destination string
	}

	// Ensure holds paths for an ensure operation (home, default, log).
	Ensure struct {
		// Home is the home path to ensure
		Home    string
		// Default is the default path to ensure
		Default Pair
		// Log is the log path to ensure
		Log     Pair
		// Log is the log path to ensure
	}

	// Pair holds a file path and a directory path, used in test fixtures.
	Pair struct {
		// File is the file path
		File      string
		// Directory is the directory path
		Directory string
	}

	// MakeDir holds paths for single directory and recursive make-dir operations.
	MakeDir struct {
		// Single is the single directory to make
		Single  string
		// MakeAll is the all directories to make
		MakeAll string
	}

	// Move holds paths for a move operation (from, destination, to).
	Move struct {
		// From is the path of the item to move
		From        Pair
		// Destination is the path to move the item to
		Destination string
		// To is the path to move the item to
		To          Pair
	}

	// Remove holds the path of the file or directory to remove.
	Remove struct {
		// File is the file to remove
		File string
	}

	// Rename holds from and to paths for a rename operation.
	Rename struct {
		// From is the path of the item to rename
		From Pair
		// To is the path to rename the item to
		To   Pair
	}

	// Write holds destination path and content for a write operation.
	Write struct {
		// Destination is the path to write the item to
		Destination string
		// Content is the content to write to the item
		Content     []byte
	}

	// StaticFs holds fixed paths for file system tests (change,
	// copy, create, ensure, etc.).
	StaticFs struct {
		// Change is the change operation
		Change   Change
		// Copy is the copy operation
		Copy     Copy
		// Create is the create operation
		Create   Create
		// Ensure makes sure that a path exists and
		// creates it if it does not exist.
		Ensure   Ensure
		// Existing is the existing path
		Existing Pair
		// MakeDir is the make directory operation
		MakeDir  MakeDir
		// Move is the move operation
			Move     Move
		// Remove is the remove operation
		Remove   Remove
		// Rename is the rename operation
		Rename   Rename
		// Scratch is the scratch path
		Scratch  string
		Write    Write
	}
	// StaticOs is a placeholder for OS-level test fixture data.
	StaticOs struct{}
)

var (
	// Perms holds default file and directory permissions used in tests.
	Perms = struct {
		// File is the default file permission
		File fs.FileMode
		// Dir is the default directory permission
		Dir  fs.FileMode
	}{
		File: 0o666, //nolint:mnd // ok (pedantic)
		Dir:  0o777, //nolint:mnd // ok (pedantic)
	}

	// Static holds the default test fixture (FS paths, OS placeholder, and Foo).
	Static = struct {
		// Foo is the foo path
		Foo string
		// FS is the file system test fixture
		FS  StaticFs
		OS  StaticOs
	}{
		Foo: "foo",
		FS: StaticFs{
			Change: Change{
				From: Pair{
					File:      "scratch/mad-as-hell.CHANGE-FROM.txt",
					Directory: "scratch/no-geography-CHANGE-FROM",
				},
				Destination: "scratch/no-geography-CHANGE-TO",
				To: Pair{
					File:      "mad-as-hell.CHANGE-TO.txt",
					Directory: "scratch/no-geography-CHANGE-TO",
				},
			},
			Copy: Copy{
				Destination: "scratch/paradise-lost.txt",
			},
			Create: Create{
				Destination: "scratch/pictures-of-you.CREATE.txt",
			},
			Ensure: Ensure{
				Home: "home/marina",
				Default: Pair{
					File:      "scratch/home/marina/logs/default-test.log",
					Directory: "scratch/home/marina/logs",
				},
				Log: Pair{
					File:      "scratch/home/marina/logs/test.log",
					Directory: "scratch/home/marina/logs",
				},
			},
			Existing: Pair{
				File:      "data/fS/paradise-lost.txt",
				Directory: "data/fS",
			},
			MakeDir: MakeDir{
				Single:  "leftfield",
				MakeAll: "scratch/leftfield/tourism",
			},
			Move: Move{
				From: Pair{
					File:      "scratch/the-same-deep-water-as-you.MOVE-FROM.txt",
					Directory: "scratch/closedown-MOVE-FROM",
				},
				Destination: "scratch/disintegration",
				To: Pair{
					File:      "scratch/disintegration/the-same-deep-water-as-you.MOVE-FROM.txt",
					Directory: "scratch/disintegration/closedown-MOVE-FROM",
				},
			},
			Remove: Remove{
				File: "scratch/paradise-regained.REMOVE.txt",
			},
			Rename: Rename{
				From: Pair{
					File:      "scratch/love-under-will.RENAME-FROM.txt",
					Directory: "scratch/earth-inferno-FROM",
				},
				To: Pair{
					File:      "scratch/love-under-will.RENAME-TO.txt",
					Directory: "scratch/earth-inferno-TO",
				},
			},
			Scratch: "scratch",
			Write: Write{
				Destination: "scratch/disintegration.WRITE.txt",
				Content:     []byte("disintegration"),
			},
		},
		OS: StaticOs{},
	}
)
