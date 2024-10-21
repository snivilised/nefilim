package nef

import (
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/snivilised/nefilim/internal/third/lo"
)

const (
	changeOpName = "Change"
)

type (
	changer interface {
		create() changer
		change(from, to string) error
	}

	changeFunc func(from, to string) error

	changers map[bitmask]changeFunc

	baseChanger struct {
		root    string
		fS      ChangerFS
		actions changers
	}
)

func (m *baseChanger) change(from, to string) error {
	mask := m.query(from, to)

	if action, exists := m.actions[mask]; exists {
		return action(from, to)
	}

	return NewInvalidBinaryFsOpError(moveOpName, from, to)
}

func (m *baseChanger) query(from, to string) bitmask {
	fromExists, fromIsDir := m.peek(from)
	toExists, toIsDir := m.peek(to)

	return bitmask{
		fromExists: fromExists,
		toExists:   toExists,
		fromIsDir:  fromIsDir,
		toIsDir:    toIsDir,
	}
}

func (m *baseChanger) peek(name string) (exists, isDir bool) { // => generic bitmap mgr
	if m.fS.DirectoryExists(name) {
		return true, true
	}

	if m.fS.FileExists(name) {
		return true, false
	}

	return false, false
}

func (m *baseChanger) changeItemWithName(from, to string) error {
	// 'to' includes the file name eg:
	// from/file.txt => to/file.txt
	//
	if strings.Contains(to, "/") {
		return NewInvalidPathError(to)
	}

	directory, _ := SplitParent(from)
	destination := join(directory, to)

	if from == destination {
		return nil
	}

	return os.Rename(
		filepath.Join(m.root, from),
		filepath.Join(m.root, destination),
	)
}

func (m *baseChanger) changeItemWithoutName(from, to string) error {
	// 'to' does not include the file name, so it has to be appended, eg:
	// from/file.txt => to/
	//
	return os.Rename(
		filepath.Join(m.root, from),
		filepath.Join(m.root, to, filepath.Base(from)),
	)
}

func (m *baseChanger) changeItemWithoutNameClash(from, to string) error {
	if strings.Contains(to, "/") {
		return NewInvalidPathError(to)
	}

	directory, _ := SplitParent(from)

	return os.Rename(
		filepath.Join(m.root, from),
		filepath.Join(m.root, directory, to),
	)
}

type lazyChanger struct { // create a generic lazyObject[T=op-interface, F=fS]
	once    sync.Once
	changer changer
}

func (l *lazyChanger) instance(root string, overwrite bool, fS ChangerFS) changer {
	l.once.Do(func() {
		l.changer = l.create(root, overwrite, fS)
	})

	return l.changer
}

func (l *lazyChanger) create(root string, overwrite bool, fS ChangerFS) changer {
	// create an interface for this function
	//
	return lo.TernaryF(overwrite,
		func() changer {
			return &overwriteChanger{
				baseChanger: baseChanger{
					root: root,
					fS:   fS,
				},
			}
		},
		func() changer {
			return &tentativeChanger{
				baseChanger: baseChanger{
					root: root,
					fS:   fS,
				},
			}
		},
	).create()
}

// changer.base

func (m *baseChanger) rejectOrNoOp(from, to string) error { // this is not a reject, rename this
	if strings.Contains(to, "/") {
		return NewInvalidPathError(to)
	}

	// both file names exists, but they may or may not be the same item. If
	// they are not in the same location then we reject the overwrite attempt
	// otherwise they are the same item and this should effectively be a no op.
	//

	directory, file := SplitParent(from)

	if filepath.Base(file) == to {
		return nil
	}

	return os.Rename(
		// should be delegated to a default of some kind as this is default behaviour...
		filepath.Join(m.root, from),
		filepath.Join(m.root, join(directory, to)),
	)
}

// changer.overwrite

type overwriteChanger struct {
	baseChanger
}

func (m *overwriteChanger) create() changer {
	m.actions = changers{
		{true, false, false, false}: m.changeItemWithName, // from exists as file, to does not exist
		{true, false, true, false}:  m.changeItemWithName, // from exists as dir, to does not exist
		// {true, true, false, true}:   m.moveItemWithoutName,                      // from exists as file,to exists as dir
		{true, true, true, true}:   m.changeItemWithoutNameClash, // from exists as dir, to exists as dir
		{true, true, false, false}: m.rejectOrNoOp,               // from and to refer to the same existing file
	}

	return m
}

func (m *overwriteChanger) changeItemWithName(from, to string) error {
	// 'to' includes the file name eg: !!! this can't be true, to can't contain /
	// from/file.txt => to/file.txt
	//
	if strings.Contains(to, "/") {
		return NewInvalidPathError(to)
	}

	directory, _ := SplitParent(from)
	destination := join(directory, to)

	if from == destination {
		return nil
	}

	return os.Rename(
		filepath.Join(m.root, from),
		filepath.Join(m.root, destination),
	)
}

// changer.tentative

type tentativeChanger struct {
	baseChanger
}

func (m *tentativeChanger) create() changer {
	m.actions = changers{
		{true, false, false, false}: m.changeItemWithName,      // from exists as file, to does not exist
		{true, false, true, false}:  m.changeDirectoryWithName, // from exists as dir, to does not exist
		// {true, true, false, true}:   m.moveItemWithoutName,      // from exists as file,to exists as dir
		{true, true, true, true}:   m.changeItemWithoutNameClash, // from exists as dir, to exists as dir
		{true, true, false, false}: m.rejectOrNoOp,               // from and to may refer to the same existing file
	}

	return m
}

func (m *tentativeChanger) changeDirectoryWithName(from, to string) error { // consolidate
	// 'to' includes the file name eg:
	// from/file.txt => to/file.txt
	//
	return m.changeItemWithoutNameClash(from, to) // TODO: check this is correct
}
