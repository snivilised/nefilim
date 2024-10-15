package nef

import (
	"path/filepath"
)

type tentativeMover struct {
	baseMover
}

func (m *tentativeMover) create() mover {
	m.actions = movers{
		{true, false, false, false}: m.moveItemWithName,         // from exists as file, to does not exist
		{true, false, true, false}:  m.moveDirectoryWithName,    // from exists as dir, to does not exist
		{true, true, false, true}:   m.moveItemWithoutName,      // from exists as file,to exists as dir
		{true, true, true, true}:    m.moveItemWithoutNameClash, // from exists as dir, to exists as dir
		{true, true, false, false}:  m.rejectOverwriteOrNoOp,    // from and to may refer to the same existing file
	}

	return m
}

func (m *tentativeMover) moveDirectoryWithName(from, to string) error {
	// 'to' includes the file name eg:
	// from/file.txt => to/file.txt
	//
	if filepath.Dir(from) == filepath.Dir(to) {
		return NewRejectSameDirMoveError(moveOpName, from, to)
	}

	return m.moveItemWithName(from, to)
}

func (m *tentativeMover) moveItemWithoutName(from, to string) error {
	// 'to' does not include the file name, so it has to be appended, eg:
	// from/file.txt => to/
	//
	if _, err := m.fS.Stat(filepath.Join(to, filepath.Base(from))); err == nil {
		return NewInvalidBinaryFsOpError("Move", from, to)
	}

	return m.baseMover.moveItemWithoutName(from, to)
}

func (m *tentativeMover) rejectOverwriteOrNoOp(from, to string) error {
	// both file names exists, but they may or may not be the same item. If
	// they are not in the same location then we reject the overwrite attempt
	// otherwise they are the same item and this should effectively be a no op.
	//
	if filepath.Dir(from) != filepath.Dir(to) {
		return NewInvalidBinaryFsOpError(moveOpName, from, to)
	}

	return nil
}
