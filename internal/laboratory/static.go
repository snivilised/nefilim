package lab

import (
	"io/fs"
)

type (
	Copy struct {
		Destination string
	}
	Create struct {
		Destination string
	}

	Pair struct {
		File      string
		Directory string
	}

	MakeDir struct {
		Single  string
		MakeAll string
	}

	Move struct {
		From        Pair
		Destination string
		To          Pair
	}

	Remove struct {
		File string
	}

	Rename struct {
		From Pair
		To   Pair
	}

	Write struct {
		Destination string
		Content     []byte
	}

	StaticFs struct {
		Copy     Copy
		Create   Create
		Existing Pair
		MakeDir  MakeDir
		Move     Move
		Remove   Remove
		Rename   Rename
		Scratch  string
		Write    Write
	}
	StaticOs struct{}
)

var (
	Perms = struct {
		File fs.FileMode
		Dir  fs.FileMode
	}{
		File: 0o666, //nolint:mnd // ok (pedantic)
		Dir:  0o777, //nolint:mnd // ok (pedantic)
	}

	Static = struct {
		Foo string
		FS  StaticFs
		OS  StaticOs
	}{
		Foo: "foo",
		FS: StaticFs{
			Copy: Copy{
				Destination: "scratch/paradise-lost.txt",
			},
			Create: Create{
				Destination: "scratch/pictures-of-you.CREATE.txt",
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
